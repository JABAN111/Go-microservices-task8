package client

import (
	"context"

	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core/errors"
	wordspb "yadro.com/course/proto/words"
)

type WordsClient struct {
	log    *slog.Logger
	client wordspb.WordsClient
	conn   *grpc.ClientConn
}

func NewWordsClient(address string, log *slog.Logger) (*WordsClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to WordsService", "error", err)
		return nil, err
	}

	return &WordsClient{
		log:    log,
		client: wordspb.NewWordsClient(conn),
		conn:   conn,
	}, nil
}

func (c *WordsClient) Close() error {
	if err := c.conn.Close(); err != nil {
		c.log.Error("ERROR while closing connection:", "error", err)
		return err
	}
	c.log.Debug("Words client are closed")
	return nil
}

func (c *WordsClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *WordsClient) Norm(ctx context.Context, phrase string) ([]string, error) {
	result, err := c.client.Norm(ctx, &wordspb.WordsRequest{Phrase: phrase})
	if err != nil {

		if status.Code(err) == codes.ResourceExhausted {
			return nil, errors.ErrResourceExhausted
		}
		return nil, err
	}
	return result.Words, nil
}
