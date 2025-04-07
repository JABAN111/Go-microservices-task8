package rest

import (
	"context"
	"time"
)

// custom implementation of token bucket based on semaphore

type Bucket struct {
	tokensCh chan struct{}
}

func NewRateLimiter(ctx context.Context, maxRequest int) *Bucket {
	buck := &Bucket{
		tokensCh: make(chan struct{}, maxRequest),
	}
	for range maxRequest {
		buck.tokensCh <- struct{}{}
	}

	refillInterval := time.Millisecond * 1200 / time.Duration(maxRequest)
	go buck.process(ctx, time.Duration(refillInterval))
	return buck
}

func (b *Bucket) process(ctx context.Context, refillInterval time.Duration) {
	tick := time.NewTicker(refillInterval)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			select {
			case b.tokensCh <- struct{}{}:
			default:
			}
		}
	}
}

func (b *Bucket) Wait(ctx context.Context) error {
	select {
	case <-b.tokensCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
