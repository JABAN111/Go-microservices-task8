package normalizer

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	wordspb "yadro.com/course/proto/words"
	"yadro.com/course/words/internal/domain/service"
	"yadro.com/course/words/pkg/logger"
	"yadro.com/course/words/pkg/shutdown"
)

const (
	maxRequestSize = 20 * 1024 // 4 KiB
	language       = "english"
)

var log = logger.GetInstance()
var shtdwnCtx, _ = shutdown.GetShutdownContext()

var normalizeService = service.NewNormalizeService()

type Normalizer interface {
	Norm(ctx context.Context, in *wordspb.WordsRequest) (*wordspb.WordsReply, error)
	Ping(ctx context.Context)
}

type normalize struct {
	wordspb.UnimplementedWordsServer
}

func NewNormalizer() wordspb.WordsServer {
	return &normalize{}
}

func (n *normalize) Ping(_ context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

func (n *normalize) Norm(ctx context.Context, in *wordspb.WordsRequest) (*wordspb.WordsReply, error) {

	reqCtx, reqCancel := context.WithCancel(ctx)
	defer reqCancel()

	if len(in.Phrase) > maxRequestSize {
		log.Warn("Got request bigger than 4 KiB")
		return nil, status.Error(codes.ResourceExhausted, "request size exceeds 4 KiB")
	}

	go func() {
		select {
		case <-shtdwnCtx.Done():
			reqCancel()
		case <-ctx.Done():
		}
	}()
	words, err := normalizeService.Normalize(reqCtx, in.Phrase, language)
	if err != nil {
		log.Warn("Error has occured while normalizing words, reason", "err", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &wordspb.WordsReply{
		Words: words,
	}, nil
}
