package core

import (
	"context"
)

type SearchReply struct {
	Comics []Comics
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
	Status(ctx context.Context) (UpdateStatus, error)
	Update(ctx context.Context) error
	Stats(ctx context.Context) (*StatsReply, error)
	Drop(ctx context.Context) error
}

type SearchServicePort interface {
	GrpcClient
	Search(ctx context.Context, query string, limit int64) (*SearchReply, error)
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
