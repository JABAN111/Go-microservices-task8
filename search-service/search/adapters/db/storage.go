package db

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/lib/pq"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"yadro.com/course/search/core"
)

type DB struct {
	log        *slog.Logger
	conn       *sqlx.DB
	numWorkers int
}

func New(log *slog.Logger, address string, numWorkers int, maxOpenConn int, connLifeTime time.Duration) (*DB, error) {
	db, err := sqlx.Connect("pgx", address)
	if err != nil {
		log.Error("connection problem", "address", address, "error", err)
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxOpenConn)
	db.SetConnMaxLifetime(connLifeTime)

	log.Debug("connected to db", "address", address)
	return &DB{
		log:        log,
		conn:       db,
		numWorkers: numWorkers,
	}, nil
}

type comicsResponse struct {
	Id     int            `db:"id"`
	Url    string         `db:"url"`
	ImgUrl string         `db:"img_url"`
	Words  pq.StringArray `db:"words"`
}

func (cr *comicsResponse) toComics() core.Comics {
	return core.Comics{
		ID:     cr.Id,
		URL:    cr.Url,
		ImgUrl: cr.ImgUrl,
		Words:  cr.Words,
	}
}

func (db *DB) Ping() error {
	err := db.conn.Ping()
	db.log.Info("pinging db", "error", err)
	return err
}

func (db *DB) GetAll(ctx context.Context) ([]core.Comics, error) {
	var comicsResponse []comicsResponse

	err := db.conn.SelectContext(ctx, &comicsResponse, "SELECT id, img_url, words FROM comics")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	result := make([]core.Comics, len(comicsResponse))
	for i, resp := range comicsResponse {
		result[i] = resp.toComics()
	}

	return result, nil
}
