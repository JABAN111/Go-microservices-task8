package grpc

import (
	"context"
	"log/slog"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	searchpb "yadro.com/course/proto/search"

	"yadro.com/course/search/core"
)

type Server struct {
	searchpb.UnimplementedSearchServer
	log      *slog.Logger
	searcher core.Searcher
}

func NewServer(log *slog.Logger, searcher core.Searcher) *Server {
	return &Server{
		log:      log,
		searcher: searcher,
	}
}

func (s *Server) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	s.log.Debug("Pinged")
	return nil, nil
}

func (s *Server) Search(_ context.Context, in *searchpb.SearchRequest) (*searchpb.RecommendedComics, error) {
	recs, err := s.searcher.Search(in.Phrase, in.Limit)
	if err != nil {
		s.log.Error("Failed to search")
		return nil, status.Error(codes.Internal, "Internal error has been occured")
	}

	ansArr := make([]*searchpb.Comics, len(recs))
	for r, rec := range recs {
		ansArr[r] = &searchpb.Comics{
			Id:     rec.Comic.ID,
			ImgUrl: rec.Comic.ImgUrl,
		}
	}

	log.Debug("Sending comics: ", "comics", ansArr)
	return &searchpb.RecommendedComics{
		Comics: ansArr,
	}, nil
}
