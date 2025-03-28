package client

import (
	"context"
	"errors"
	"log/slog"

	"yadro.com/course/api/core"

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

func (c *UpdateClient) Status(ctx context.Context) (core.UpdateStatus, error) {
	reply, err := c.client.Status(ctx, nil)
	if err != nil {
		return core.StatusUpdateUnknown, err
	}
	switch reply.Status {
	case updatepb.Status_STATUS_IDLE:
		return core.StatusUpdateIdle, nil
	case updatepb.Status_STATUS_RUNNING:
		return core.StatusUpdateRunning, nil
	}
	return core.StatusUpdateUnknown, errors.New("unknown status")
}

func (c *UpdateClient) Update(ctx context.Context) error {
	_, err := c.client.Update(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *UpdateClient) Stats(ctx context.Context) (*core.StatsReply, error) {
	stats, err := c.client.Stats(ctx, &emptypb.Empty{})
	if err != nil {
		return &core.StatsReply{}, err
	}
	return &core.StatsReply{
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
