package core

import (
	"context"
)

type GrpcClient interface {
	Close() error
	Pinger
}
type Pinger interface {
	Ping(context.Context) error
}

// Worder :D TODO: how to adapt this?
type Worder interface {
	GrpcClient
	Norm(ctx context.Context, phrase string) ([]string, error)
}

type Updater interface {
	GrpcClient
	Status(ctx context.Context) (UpdateStatus, error)
	Update(ctx context.Context) error
	Stats(ctx context.Context) (*Stats, error)
	Drop(ctx context.Context) error
}

type Searcher interface {
	GrpcClient
	Search(ctx context.Context, query string, limit int64) ([]Comics, error)
	ISearch(ctx context.Context, query string, limit int64) ([]Comics, error)
}
