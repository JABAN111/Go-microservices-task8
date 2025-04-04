package grpc

import (
	"context"
	"log/slog"

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
		s.log.Error("Search failed", "error", err)
		return nil, status.Error(codes.Internal, "An internal error has occurred")
	}

	s.log.Debug("Returning search results", "results", recs)
	return s.convertToRecommendedComics(recs), nil
}

func (s *Server) ISearch(ctx context.Context, in *searchpb.SearchRequest) (*searchpb.RecommendedComics, error) {
	recs, err := s.searcher.ISearch(ctx, in.Phrase, in.Limit)
	if err != nil {
		s.log.Error("ISearch failed", "error", err)
		return nil, status.Error(codes.Internal, "An internal error has occurred")
	}

	s.log.Debug("Returning ISearch results", "results", recs)
	return s.convertToRecommendedComics(recs), nil
}

func (s *Server) convertToRecommendedComics(recs []core.ComicMatch) *searchpb.RecommendedComics {
	ansArr := make([]*searchpb.Comics, len(recs))
	for i, rec := range recs {
		ansArr[i] = &searchpb.Comics{
			Id:     int64(rec.Comic.ID),
			ImgUrl: rec.Comic.ImgUrl,
		}
	}
	return &searchpb.RecommendedComics{Comics: ansArr}
}
