package ports

import (
	"context"
)

type StatsReply struct {
	WordsTotal    int64
	WordsUnique   int64
	ComicsTotal   int64
	ComicsFetched int64
}

type Pinger interface {
	Ping(context.Context) error
}

type WordsServicePort interface {
	GrpcClient
	Norm(ctx context.Context, phrase string) ([]string, error)
}

type UpdateServicePort interface {
	GrpcClient
	Ping(ctx context.Context) error
	Status(ctx context.Context) (string, error)
	Update(ctx context.Context) error
	Stats(ctx context.Context) (*StatsReply, error)
	Drop(ctx context.Context) error
}

type GrpcClient interface {
	Close() error
	Ping(ctx context.Context) error
}

type GrpcManager interface {
	Register(name string, client GrpcClient)
	GetClient(name string) (GrpcClient, error)
	CloseAll(ctx context.Context)
	PingAll(ctx context.Context) map[string]string
}
