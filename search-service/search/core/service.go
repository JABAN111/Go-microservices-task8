package core

import (
	"context"
	"log/slog"
	"sort"
	"sync"
)

type Service struct {
	log        *slog.Logger
	db         DB
	words      Words
	numWorkers int
	mu         sync.Mutex
}

func NewService(log *slog.Logger, db DB, words Words, numWorkers int) *Service {
	return &Service{
		log:        log,
		db:         db,
		words:      words,
		numWorkers: numWorkers,
		mu:         sync.Mutex{},
	}
}

func (s *Service) Search(phrase string, limit int64) ([]ComicMatch, error) {
	s.log.Debug("Got request on search: ", "phrase", phrase, "limit", limit)

	recs, err := s.search(phrase)

	if err != nil {
		return nil, err
	}
	len := min(int(limit), len(recs))
	finalRecs := recs[:len]

	return finalRecs, nil
}

func (s *Service) search(phrase string) ([]ComicMatch, error) {
	sema := make(chan struct{}, s.numWorkers)

	resultChan := make(chan ComicMatch)

	keyWords, err := s.words.Norm(context.Background(), phrase)
	if err != nil {
		s.log.Error("Failed to normalize the phrase")
		return nil, err
	}
	keyWordsSet := make(map[string]bool)
	for _, word := range keyWords {
		keyWordsSet[word] = true
	}

	comics, err := s.db.GetAll(context.Background())
	if err != nil {
		s.log.Error("Error has been occured while getting all comics", "error", err)
		return nil, err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(comics))

	for _, comic := range comics {
		go func() {
			sema <- struct{}{}
			defer func() {
				<-sema
				wg.Done()
			}()

			resultChan <- ComicMatch{
				Comic: comic,
				Count: matchInComics(keyWordsSet, comic.Words)}
		}()
	}

	resultSlice := make([]ComicMatch, 0)
	go func() {
		for res := range resultChan {
			resultSlice = append(resultSlice, res)
		}
		close(resultChan)
	}()

	wg.Wait()
	sort.Slice(resultSlice, func(i, j int) bool {
		return resultSlice[i].Count > resultSlice[j].Count
	})
	return resultSlice, nil
}

func matchInComics(keyWordsSet map[string]bool, words []string) int {
	counter := 0
	for _, word := range words {
		if keyWordsSet[word] {
			counter++
		}
	}
	return counter
}
