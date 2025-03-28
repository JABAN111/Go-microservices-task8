package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"yadro.com/course/pkg/util"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/search/core"
)

type DB struct {
	log        *slog.Logger
	conn       *sqlx.DB
	numWorkers int
}

func New(log *slog.Logger, address string, numWorkers int) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}
	log.Debug("connected to db", "address", address)
	return &DB{
		log:        log,
		conn:       db,
		numWorkers: numWorkers,
	}, nil
}

func (db *DB) Ping() error {
	err := db.conn.Ping()
	db.log.Info("pinging db", "error", err)
	return err
}

func (db *DB) GetAll(ctx context.Context) ([]core.Comics, error) {
	result := make([]core.Comics, 0)

	rows, err := db.conn.DB.QueryContext(ctx, "SELECT id, img_url, words FROM comics")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer util.SafeClose(rows)

	for rows.Next() {
		var comic core.Comics
		err := rows.Scan(&comic.ID, &comic.ImgUrl, pq.Array(&comic.Words))
		if err != nil {
			db.log.Warn("Error while scanning rows", "error", err)
		}

		result = append(result, comic)
	}
	return result, nil
}
