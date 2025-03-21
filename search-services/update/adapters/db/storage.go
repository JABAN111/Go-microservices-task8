package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"yadro.com/course/update/core"
)

type DB struct {
	log        *slog.Logger
	conn       *sqlx.DB
	batchSize  int
	numWorkers int
}

func New(log *slog.Logger, address string, batchSize, numWorkers int) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}
	log.Debug("connected to db", "address", address)
	return &DB{
		log:        log,
		conn:       db,
		batchSize:  batchSize,
		numWorkers: numWorkers,
	}, nil
}

func (db *DB) Ping() error {
	err := db.conn.Ping()
	db.log.Info("pinging db", "error", err)
	return err
}

func (db *DB) Add(ctx context.Context, comics core.Comics) error {
	res, err := db.conn.DB.ExecContext(ctx, "INSERT INTO comics (id, url, words) values ($1, $2, $3)", comics.ID, comics.URL, comics.Words)
	if err != nil {
		log.Error("error adding comic", "error", err)
	}
	db.log.Debug("Add finished", "res", res)
	return nil
}

func (db *DB) AddAllComics(ctx context.Context, comicsList <-chan core.Comics) error {
	var comics []map[string]interface{}
	var mu sync.Mutex
	sem := make(chan struct{}, db.numWorkers)
	var wg sync.WaitGroup

	for comic := range comicsList {
		sem <- struct{}{}
		wg.Add(1)
		go func(comic core.Comics) {
			defer func() {
				<-sem
				wg.Done()
			}()

			data := map[string]interface{}{
				"id":    comic.ID,
				"url":   comic.URL,
				"words": comic.Words,
			}

			mu.Lock()
			comics = append(comics, data)
			mu.Unlock()
		}(comic)
	}

	wg.Wait()

	if len(comics) == 0 {
		return nil
	}

	query := "INSERT INTO comics (id, url, words) VALUES (:id, :url, :words)"
	_, err := db.conn.NamedExecContext(ctx, query, comics)
	if err != nil {
		db.log.Error("Failed to batch insert comics", "error", err)
		return err
	}

	return nil
}

func (db *DB) DbStats(ctx context.Context) (core.DBStats, error) {
	var stats core.DBStats

	rows, err := db.conn.DB.QueryContext(ctx, "SELECT words_total, words_unique, comics_fetched FROM db_stats order by id desc limit 1;")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return core.DBStats{
				WordsTotal:    0,
				WordsUnique:   0,
				ComicsFetched: 0,
			}, nil
		}
		return stats, err
	}
	defer core.SafeClose(rows)

	if rows.Next() {
		err := rows.Scan(&stats.WordsTotal, &stats.WordsUnique, &stats.ComicsFetched)
		if err != nil {
			return stats, err
		}
	} else {
		return core.DBStats{
			WordsTotal:    0,
			WordsUnique:   0,
			ComicsFetched: 0,
		}, nil
	}

	return stats, nil
}

func (db *DB) UpdateStats(ctx context.Context, cntUniqueWords int, comicsInTotal int) error {
	return db.updateDbStatsAndServiceStats(ctx, cntUniqueWords, comicsInTotal)
}

func (db *DB) updateDbStatsAndServiceStats(ctx context.Context, cntUniqueWords int, comicsInTotal int) (err error) {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer commitOrRollback(tx, &err, db.log)

	dbStatsStmt, err := tx.PrepareContext(ctx, `
		INSERT INTO db_stats (words_total, words_unique, comics_fetched)
		SELECT 
    		(select sum(count) from words_stats) AS words_total,
   			COALESCE(
   				(SELECT words_unique FROM db_stats ORDER BY id desc LIMIT 1), 0) + $1 AS words_unique,
    		(SELECT COUNT(*) FROM comics) AS comics_fetched
		RETURNING id;
	`)
	if err != nil {
		return err
	}
	defer core.SafeClose(dbStatsStmt)

	var dbStatsId int
	err = dbStatsStmt.QueryRowContext(ctx, cntUniqueWords).Scan(&dbStatsId)
	if err != nil {
		return err
	}

	serviceStatsStmt, err := tx.PrepareContext(ctx, `INSERT INTO service_stats values ($1, $2)`)
	if err != nil {
		return err
	}
	defer core.SafeClose(serviceStatsStmt)

	_, err = serviceStatsStmt.ExecContext(ctx, dbStatsId, comicsInTotal)
	if err != nil {
		return err
	}

	db.log.Info("Inserted new db_stats record", "id", dbStatsId)
	return nil
}

func (db *DB) IDs(ctx context.Context) ([]int, error) {
	res, err := db.conn.DB.QueryContext(ctx, "SELECT id FROM comics order by id")
	if err != nil {
		return nil, err
	}
	defer core.SafeClose(res)

	var ids []int
	for res.Next() {
		var id int
		err := res.Scan(&id)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (db *DB) Drop(ctx context.Context) error {
	return db.drop(ctx)
}

func (db *DB) drop(ctx context.Context) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		db.log.Error("Failed to begin transaction", "error", err)
		return err
	}
	defer commitOrRollback(tx, &err, db.log)

	queries := []string{
		"DELETE FROM comics CASCADE",
		"DELETE FROM service_stats CASCADE",
		"DELETE FROM db_stats CASCADE",
	}

	for _, query := range queries {
		if _, err = tx.ExecContext(ctx, query); err != nil {
			db.log.Error("Failed to execute query", "query", query, "error", err)
			return err
		}
	}

	return nil
}

func (db *DB) AddWordStats(ctx context.Context, wordsArr []string) error {
	const batchSize = 100
	for i := 0; i < len(wordsArr); i += batchSize {
		endOfBatch := i + batchSize
		if endOfBatch > len(wordsArr) {
			endOfBatch = len(wordsArr)
		}

		_, err := db.addWordStats(ctx, wordsArr[i:endOfBatch])
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) addWordStats(ctx context.Context, wordsArr []string) ([]string, error) {
	//TODO: своеобразный костылек, надо бы залатать
	for i, word := range wordsArr {
		wordsArr[i] = strings.ReplaceAll(word, "'", "''")
	}

	joinedWords := "'" + strings.Join(wordsArr, "', '") + "'"

	wordsArray := fmt.Sprintf("ARRAY[%s]", joinedWords)

	query := fmt.Sprintf("SELECT add_word_stats(%s)", wordsArray)

	_, err := db.conn.ExecContext(ctx, query)
	if err != nil {
		db.log.Error("Failed to insert batch using function", "error", err)
		return nil, err
	}

	return wordsArr, nil
}

func commitOrRollback(tx *sql.Tx, err *error, log *slog.Logger) {
	if *err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			log.Error("Failed to rollback transaction", "rollback_error", rollbackErr)
		} else {
			log.Debug("Transaction rolled back due to error")
		}
	} else {
		if commitErr := tx.Commit(); commitErr != nil {
			log.Error("Failed to commit transaction", "commit_error", commitErr)
		} else {
			log.Debug("Transaction committed")
		}
	}
}
