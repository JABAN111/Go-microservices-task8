package core

import (
	"context"
)

type Updater interface {
	Update(context.Context) error
	Stats(context.Context) (ServiceStats, error)
	Status(context.Context) ServiceStatus
	Drop(context.Context) error
}

type DB interface {
	AddAllComics(context.Context, <-chan Comics) error
	Add(context.Context, Comics) error
	DbStats(context.Context) (DBStats, error)
	Drop(context.Context) error
	IDs(context.Context) ([]int, error)
	AddWordStats(ctx context.Context, wordsList []string) error
	UpdateStats(ctx context.Context, cntUniqueWords int, comicsInTotal int) error
}

type XKCD interface {
	Get(context.Context, int) (XKCDInfo, error)
	LastID(context.Context) (int, error)
}

type Words interface {
	Norm(ctx context.Context, phrase string) ([]string, error)
	GetWords(ctx context.Context, phrase string) ([]string, error)
}
