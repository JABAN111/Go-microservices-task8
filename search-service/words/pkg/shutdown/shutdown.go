package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"yadro.com/course/words/pkg/logger"
)

var log = logger.GetInstance()
var (
	once         sync.Once
	shtdwnCtx    context.Context
	shtdwncancel context.CancelFunc
)

func GetShutdownContext() (context.Context, context.CancelFunc) {
	once.Do(func() {
		shtdwnCtx, shtdwncancel = context.WithCancel(context.Background())
	})
	return shtdwnCtx, shtdwncancel
}

func RunWithShutdown(s *grpc.Server, startServer func() error) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		<-ctx.Done()
		log.Info("Shutdown initiated, stopping gRPC server...")
		s.GracefulStop()
		time.Sleep(1 * time.Second) //NOTE: giving time to finish all operation of normalizing
		log.Info("gRPC server stopped")
		os.Exit(1)
		return nil
	})

	group.Go(startServer)

	return group.Wait()
}
