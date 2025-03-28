package client

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"yadro.com/course/api/core"
	searchpb "yadro.com/course/proto/search"
)

type SearchClient struct {
	log    *slog.Logger
	client searchpb.SearchClient
	conn   *grpc.ClientConn
}

func NewSearchClient(address string, log *slog.Logger) (*SearchClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to WordsService", "error", err)
		return nil, err
	}

	return &SearchClient{
		log:    log,
		client: searchpb.NewSearchClient(conn),
		conn:   conn,
	}, nil
}

func (c *SearchClient) Close() error {
	if err := c.conn.Close(); err != nil {
		c.log.Error("ERROR while closing connection:", "error", err)
		return err
	}
	c.log.Debug("Words client are closed")
	return nil
}

func (c *SearchClient) Ping(ctx context.Context) error {
	_, err := c.client.Ping(ctx, &emptypb.Empty{})
	return err
}

func (c *SearchClient) Search(ctx context.Context, query string, limit int64) (*core.SearchReply, error) {
	recs, err := c.client.Search(ctx, &searchpb.SearchRequest{
		Phrase: query,
		Limit:  limit,
	})
	if err != nil {
		return &core.SearchReply{}, err
	}

	resultComics := make([]core.Comics, len(recs.Comics))
	for i, rec := range recs.Comics {
		resultComics[i] = core.Comics{ImgUrl: rec.ImgUrl, ID: rec.Id}
	}

	return &core.SearchReply{
		Comics: resultComics,
	}, nil

}
