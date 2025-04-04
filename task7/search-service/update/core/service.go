package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"
)

type Service struct {
	log        *slog.Logger
	db         DB
	xkcd       XKCD
	words      Words
	numWorkers int
	isUpdating atomic.Bool
	mu         sync.Mutex
}

func NewService(
	log *slog.Logger, db DB, xkcd XKCD, words Words, numWorkers int,
) (*Service, error) {
	if numWorkers < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", numWorkers)
	}
	return &Service{
		log:        log,
		db:         db,
		xkcd:       xkcd,
		words:      words,
		numWorkers: numWorkers,
	}, nil
}

func (s *Service) Update(ctx context.Context) error {
	start := time.Now()
	if ok := s.mu.TryLock(); !ok {
		return ErrAlreadyExists
	}
	defer s.mu.Unlock()
	defer s.isUpdating.Store(false)
	s.isUpdating.Store(true)

	existingIDs, err := s.db.IDs(ctx)
	hashSetExisting := s.sliceToMap(existingIDs)

	if err != nil {
		return fmt.Errorf("failed to get existing comics: %w", err)
	}

	latestID, err := s.xkcd.LastID(ctx)
	if err != nil || latestID < 0 {
		return fmt.Errorf("failed to get last XKCD ID: %w", err)
	}

	idToFetch := s.getNewId(ctx, latestID, hashSetExisting)
	fetcherInfo := s.fetchMissedComics(ctx, idToFetch)
	comicsCh := s.processInfo(ctx, fetcherInfo)
	if err = s.db.AddAllComics(ctx, comicsCh); err != nil {
		s.log.Error("Fail to upload to db")
		return err
	}

	s.log.Debug("Elapsed", "time", time.Since(start))
	return nil
}

func (s *Service) Stats(ctx context.Context) (ServiceStats, error) {
	stats, err := s.db.DbStats(ctx)
	if err != nil {
		s.log.Error("Error while getting stats", "error", err)
		return ServiceStats{}, err
	}
	last, err := s.xkcd.LastID(ctx)
	if err != nil {
		s.log.Error("Fail to get last XKCD ID", "error", err)
		return ServiceStats{}, err
	}
	return ServiceStats{
		DBStats:     stats,
		ComicsTotal: last,
	}, nil
}

func (s *Service) Status(_ context.Context) ServiceStatus {
	if s.isUpdating.Load() {
		return StatusRunning
	}
	return StatusIdle
}

func (s *Service) Drop(ctx context.Context) error {
	err := s.db.Drop(ctx)
	if err != nil {
		s.log.Error("Failed to drop db")
		return err
	}

	return nil
}

func (s *Service) sliceToMap(ids []int) map[int]bool {
	idMap := make(map[int]bool, len(ids))
	for _, id := range ids {
		idMap[id] = true
	}
	return idMap
}

func (s *Service) getNewId(ctx context.Context, last int, exist map[int]bool) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 1; i <= last; i++ {
			if exist[i] {
				continue
			}
			select {
			case <-ctx.Done():
			case ch <- i:
			}
		}
	}()

	return ch
}

func (s *Service) fetchMissedComics(ctx context.Context, ids <-chan int) <-chan XKCDInfo {
	ch := make(chan XKCDInfo)
	sema := make(chan struct{}, s.numWorkers)
	var wg sync.WaitGroup

	for id := range ids {
		wg.Add(1)
		go func() {
			defer s.log.Debug("Finish to fetch", "id", id)
			defer wg.Done()
			sema <- struct{}{}
			defer func() { <-sema }()
			if id == 404 {
				ch <- XKCDInfo{ID: 404, Description: "Not found", Title: "Not found"}
				return
			}
			comic, err := s.xkcd.Get(ctx, id)
			if err != nil {
				s.log.Error("Fail to fetch comics", "error", err)
				return
			}
			select {
			case ch <- comic:
			case <-ctx.Done():
			}
		}()
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

func (s *Service) processInfo(ctx context.Context, infoCh <-chan XKCDInfo) <-chan Comics {
	comicsCh := make(chan Comics)

	var wg sync.WaitGroup
	wg.Add(s.numWorkers)

	for range s.numWorkers {
		go func() {
			defer wg.Done()
			for info := range infoCh {
				words, err := s.words.Norm(ctx,
					fmt.Sprintf("%s %s %s %s %s", info.Title, info.Description, info.Alt, info.News, info.SafeTitle))
				if err != nil {
					s.log.Warn("Fail to normalize comics data", "error", err)
					continue
				}
				comicsCh <- Comics{
					ID:     info.ID,
					URL:    info.URL,
					ImgUrl: info.ImgUrl,
					Words:  words,
				}
			}
		}()
	}
	go func() {
		wg.Wait()
		close(comicsCh)
	}()

	return comicsCh
}
