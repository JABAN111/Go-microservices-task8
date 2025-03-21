package core

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	log "github.com/sirupsen/logrus"
)

type Service struct {
	log        *slog.Logger
	db         DB
	xkcd       XKCD
	words      Words
	numWorkers int
	mu         sync.Mutex
	isUpdating bool
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
	if s.isUpdating {
		return status.Errorf(codes.FailedPrecondition, "update already in progress")
	}

	s.mu.Lock()
	s.isUpdating = true
	defer func() {
		s.isUpdating = false
		s.mu.Unlock()
	}()
	s.log.Debug("Starting update process")

	existingIDs, err := s.db.IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to get existing comics: %w", err)
	}

	hashSetExisting := s.sliceToMap(existingIDs)

	latestID, err := s.xkcd.LastID(ctx)
	if err != nil || latestID < 0 {
		return fmt.Errorf("failed to get last XKCD ID: %w", err)
	}

	sema := make(chan struct{}, s.numWorkers)
	localMu := sync.Mutex{}
	comicsChan := make(chan Comics)
	wordsArr := make([]string, 0)
	cntUniqueWords := 0

	var wg sync.WaitGroup

	for id := 1; id <= latestID; id++ {
		if hashSetExisting[id] {
			continue
		}

		wg.Add(1)

		go func(id int) {
			log.Debug("Goroutine started")
			defer wg.Done()
			if err := ctx.Err(); err != nil {
				s.log.Warn("Update process canceled")
				return
			}
			sema <- struct{}{}

			defer func() { <-sema }()

			xkcdInfo, err := s.xkcd.Get(ctx, id)
			if err != nil {
				s.log.Warn("Failed to fetch XKCD data", "id", id, "error", err)
				return
			}

			normWords, err := s.words.Norm(ctx, xkcdInfo.Description)

			if err != nil {
				s.log.Warn("Failed to normalize description", "id", id, "error", err)
				return
			}
			words, err := s.words.GetWords(ctx, xkcdInfo.Description)

			if err != nil {
				log.Error("Failed to get words", "id", id, "error", err)
				return
			}
			localMu.Lock()
			wordsArr = append(wordsArr, words...)
			cntUniqueWords += len(normWords)
			localMu.Unlock()
			comicsChan <- Comics{
				ID:    id,
				URL:   xkcdInfo.URL,
				Words: normWords,
			}

			s.log.Debug("Update comics", "id", id)
		}(id)
	}

	go func() {
		wg.Wait()
		close(comicsChan)
	}()

	if err := s.db.AddAllComics(ctx, comicsChan); err != nil {
		log.Error("Failed to add all", "error", err)
		return err
	}
	if err := s.db.AddWordStats(ctx, wordsArr); err != nil {
		log.Error("Failed to add wordstats", "error", err)
		return err
	}

	if err := s.db.UpdateStats(ctx, cntUniqueWords, latestID); err != nil {
		log.Error("Failed to update comics", "error", err)
		return err
	}

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
		ComicsTotal: last - 1, // шуточный комикс 404 обкостылен
	}, nil
}

func (s *Service) Status(_ context.Context) ServiceStatus {
	if s.isUpdating {
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
	var mu sync.Mutex
	var wg sync.WaitGroup
	sema := make(chan struct{}, s.numWorkers)

	for _, id := range ids {
		wg.Add(1)
		sema <- struct{}{}

		go func(id int) {
			defer wg.Done()
			defer func() { <-sema }()

			mu.Lock()
			idMap[id] = true
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	return idMap
}
