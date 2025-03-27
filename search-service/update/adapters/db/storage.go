package db

import (
	"context"
	"log/slog"
	"sync"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"yadro.com/course/pkg/util"
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
	res, err := db.conn.ExecContext(ctx, "INSERT INTO comics (id, url, words) values ($1, $2, $3)", comics.ID, comics.URL, comics.Words)
	if err != nil {
		log.Error("error adding comic", "error", err)
	}
	db.log.Debug("Add finished", "res", res)
	return nil
}

func (db *DB) AddAllComics(ctx context.Context, comicsChan <-chan core.Comics) error {
	var comics []map[string]interface{}
	var mu sync.Mutex
	sem := make(chan struct{}, db.numWorkers)
	var wg sync.WaitGroup

	for comic := range comicsChan {
		sem <- struct{}{}
		wg.Add(1)
		go func(comic core.Comics) {
			defer func() {
				<-sem
				wg.Done()
			}()

			data := map[string]interface{}{
				"id":     comic.ID,
				"url":    comic.URL,
				"imgUrl": comic.ImgUrl,
				"words":  comic.Words,
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

	query := "INSERT INTO comics (id, url, img_url, words) VALUES (:id, :url, :imgUrl, :words)"
	res, err := db.conn.NamedExecContext(ctx, query, comics)
	if err != nil {
		db.log.Error("Failed to batch insert comics", "error", err)
		return err
	}
	aff, err := res.RowsAffected()
	if err != nil {
		return err
	}
	db.log.Debug("Inserting finished", "inserted", aff)
	return nil
}

func (db *DB) DbStats(ctx context.Context) (core.DBStats, error) {
	var stats core.DBStats
	err := db.conn.GetContext(
		ctx, &stats.ComicsFetched,
		"SELECT COUNT(*) FROM comics")
	if err != nil {
		return core.DBStats{}, err
	}
	err = db.conn.GetContext(
		ctx, &stats.WordsTotal,
		"SELECT coalesce(SUM(array_length(words,1)), 0) FROM comics",
	)
	if err != nil {
		return core.DBStats{}, err
	}
	err = db.conn.GetContext(
		ctx, &stats.WordsUnique,
		"SELECT count(*) FROM (SELECT distinct(unnest(words)) FROM comics)",
	)
	if err != nil {
		return core.DBStats{}, err
	}

	return stats, nil
}

func (db *DB) IDs(ctx context.Context) ([]int, error) {
	var ids []int

	err := db.conn.SelectContext(ctx, &ids, "SELECT id FROM comics ORDER BY id")
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (db *DB) Drop(ctx context.Context) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		db.log.Error("Failed to begin transaction", "error", err)
		return err
	}
	defer util.CommitOrRollback(tx, &err, db.log)

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
