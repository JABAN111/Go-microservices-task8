package core

import "context"

type Searcher interface {
	Search(phrase string, limit int64) ([]ComicMatch, error)
}

type DB interface {
	GetAll(ctx context.Context) ([]Comics, error)
}

type XKCD interface {
	Get(id int) Comics
}

type Words interface {
	Norm(ctx context.Context, phrase string) ([]string, error)
}
