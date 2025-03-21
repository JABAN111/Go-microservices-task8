package db_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	conn "yadro.com/course/update/adapters/db"
	"yadro.com/course/update/core"
)

func TestXxx(t *testing.T) {
	err := db.Ping()
	assert.NoError(t, err)
}

var db *conn.DB
var localConn *sqlx.DB

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=jabapwd",
			"POSTGRES_USER=jaba",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://jaba:jabapwd@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	if err = resource.Expire(120); err != nil {
		panic(err)
	}

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {

		db, err = conn.New(slog.Default(), databaseUrl, 50, 50)
		if err != nil {
			return err
		}
		localConn, err = sqlx.Connect("pgx", databaseUrl)
		if err != nil {
			log.Error("connection problem", "address", databaseUrl, "error", err)
			return err
		}

		log.Info("Database connected")
		if err = db.Migrate(); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func TestAddComicsAndCheck(t *testing.T) {
	ctx := context.Background()
	id := 1
	expected := core.Comics{
		ID:    id,
		URL:   "https://xkcd.com/1",
		Words: []string{"xkcd", "test"},
	}

	err := db.Add(ctx, expected)
	assert.NoError(t, err)

	var actual core.Comics
	err = localConn.DB.QueryRowContext(ctx, "SELECT id, url, words FROM comics WHERE id = $1", id).Scan(&actual.ID, &actual.URL, pq.Array(&actual.Words))
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual, "Added and retrieved comics should match")
}
func Test_IdsTest(t *testing.T) {
	ctx := context.Background()
	err := db.Drop(ctx)
	assert.NoError(t, err)

	id := 1
	expected := core.Comics{
		ID:    id,
		URL:   "https://xkcd.com/1",
		Words: []string{"xkcd", "test"},
	}

	err = db.Add(ctx, expected)
	assert.NoError(t, err)

	ids, err := db.IDs(ctx)
	assert.NoError(t, err)

	assert.Contains(t, ids, id, "Added comic should be in the list of IDs")
}
func TestDrop(t *testing.T) {
	localConn.MustExec("Truncate comics")

	ctx := context.Background()
	id := 1
	comic := core.Comics{
		ID:    id,
		URL:   "https://xkcd.com/1",
		Words: []string{"xkcd", "test"},
	}

	err := db.Add(ctx, comic)
	assert.NoError(t, err)

	var count int
	err = localConn.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM comics").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "Comic should exist before drop")

	err = db.Drop(ctx)
	assert.NoError(t, err)

	err = localConn.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM comics").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 0, count, "Table should be empty after drop")
}

func TestStats(t *testing.T) {
	ctx := context.Background()

	_, err := localConn.DB.ExecContext(ctx, "DELETE FROM db_stats")
	assert.NoError(t, err)

	expected := core.DBStats{
		WordsTotal:    100,
		WordsUnique:   50,
		ComicsFetched: 10,
	}
	_, err = localConn.DB.ExecContext(ctx, "INSERT INTO db_stats (words_total, words_unique, comics_fetched) VALUES ($1, $2, $3)",
		expected.WordsTotal, expected.WordsUnique, expected.ComicsFetched)
	assert.NoError(t, err)

	actual, err := db.DbStats(ctx)
	assert.NoError(t, err)

	assert.Equal(t, expected, actual, "Stats should return correct values")
}

func TestAddAllComics(t *testing.T) {
	localConn.MustExec("Truncate comics")

	ctx := context.Background()

	comicsList := []core.Comics{
		{ID: 1, URL: "https://xkcd.com/1", Words: []string{"xkcd", "test"}},
		{ID: 2, URL: "https://xkcd.com/2", Words: []string{"xkcd", "example"}},
	}

	comicsChannel := make(chan core.Comics, len(comicsList))

	go func() {
		for _, comic := range comicsList {
			comicsChannel <- comic
		}
		close(comicsChannel)
	}()

	err := db.AddAllComics(ctx, comicsChannel)
	assert.NoError(t, err)

	for _, comic := range comicsList {
		var actual core.Comics
		err = localConn.DB.QueryRowContext(ctx, "SELECT id, url, words FROM comics WHERE id = $1", comic.ID).Scan(&actual.ID, &actual.URL, pq.Array(&actual.Words))
		assert.NoError(t, err)
		assert.Equal(t, comic, actual, "Added and retrieved comics should match")
	}
}
