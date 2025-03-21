package client

import (
	"context"
	"log/slog"

	"yadro.com/course/api/core/ports"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	updatepb "yadro.com/course/proto/update"
)

type UpdateClient struct {
	log    *slog.Logger
	client updatepb.UpdateClient
	conn   *grpc.ClientConn
}

// NewUpdateClient TODO стоит ли обобщить эту логику
func NewUpdateClient(address string, log *slog.Logger) (*UpdateClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to WordsService", "error", err)
		return nil, err
	}

	return &UpdateClient{
		log:    log,
		client: updatepb.NewUpdateClient(conn),
		conn:   conn,
	}, nil
}

func (c *UpdateClient) Close() error {
	if err := c.conn.Close(); err != nil {
		c.log.Error("ERROR while closing connection:", "error", err)
		return err
	}
	c.log.Debug("Words client are closed")
	return nil
}

func (c *UpdateClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *UpdateClient) Status(ctx context.Context) (string, error) {
	status, err := c.client.Status(ctx, &emptypb.Empty{})
	if err != nil {
		return "", err
	}
	return status.String(), nil
}

func (c *UpdateClient) Update(ctx context.Context) error {
	_, err := c.client.Update(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *UpdateClient) Stats(ctx context.Context) (*ports.StatsReply, error) {
	stats, err := c.client.Stats(ctx, &emptypb.Empty{})
	if err != nil {
		return &ports.StatsReply{}, err
	}
	return &ports.StatsReply{
		WordsTotal:    stats.WordsTotal,
		WordsUnique:   stats.WordsUnique,
		ComicsTotal:   stats.ComicsTotal,
		ComicsFetched: stats.ComicsFetched,
	}, nil
}

func (c *UpdateClient) Drop(ctx context.Context) error {
	_, err := c.client.Drop(ctx, &emptypb.Empty{})
	if err != nil {
		return err
	}
	return nil
}
