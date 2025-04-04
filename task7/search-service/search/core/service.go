package core

import (
	"context"
	"log/slog"
	"maps"
	"sort"
	"sync"
	"time"
)

// дефакто comicsHashSet map[Comics]bool
// map[Comics]bool недопустимо, а хранить указателями кощунство из С
type comicsHashSet map[int]Comics

type Service struct {
	log        *slog.Logger
	db         DB
	words      Words
	numWorkers int
	mu         sync.Mutex
	indComics  map[string]comicsHashSet
	indexTTL   time.Duration
}

func NewService(ctx context.Context, log *slog.Logger, db DB, indexTTL time.Duration, words Words, numWorkers int) *Service {

	s := &Service{
		log:        log,
		db:         db,
		indexTTL:   indexTTL,
		words:      words,
		numWorkers: numWorkers,
		mu:         sync.Mutex{},
		indComics:  make(map[string]comicsHashSet),
	}
	go s.ticker(ctx)

	return s
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
	return s.getCommicsMatch(keyWordsSet, comics)
}

func (s *Service) ISearch(ctx context.Context, phrase string, limit int64) ([]ComicMatch, error) {
	keyWords, err := s.words.Norm(ctx, phrase)
	if err != nil {
		return nil, err
	}

	res := make(comicsHashSet)

	keyWordsSet := make(map[string]bool)
	for _, word := range keyWords {
		mergeHashSets(res, s.indComics[word])
		keyWordsSet[word] = true
	}

	var comics []Comics
	for _, comic := range res {
		comics = append(comics, comic)
	}

	recsSorted, err := s.getCommicsMatch(keyWordsSet, comics)
	if err != nil {
		return nil, err
	}

	len := min(int(limit), len(recsSorted))
	finalRecs := recsSorted[:len]

	return finalRecs, nil

}

func mergeHashSets(a, b comicsHashSet) comicsHashSet {
	maps.Copy(a, b)
	return a
}

func (s *Service) getCommicsMatch(keyWordsSet map[string]bool, comics []Comics) ([]ComicMatch, error) {
	sema := make(chan struct{}, s.numWorkers)
	resultChan := make(chan ComicMatch)

	wg := sync.WaitGroup{}
	wg.Add(len(comics))

	for _, comic := range comics {
		go func() {
			sema <- struct{}{}
			defer func() {
				wg.Done()
				<-sema
			}()
			resultChan <- ComicMatch{
				Comic: comic,
				Count: matchInComics(keyWordsSet, comic.Words)}
		}()
	}

	var resultSlice []ComicMatch
	go func() {
		wg.Wait()
		close(resultChan)
	}()
	for res := range resultChan {
		resultSlice = append(resultSlice, res)
	}
	sort.Slice(resultSlice, func(i, j int) bool {
		return resultSlice[i].Count > resultSlice[j].Count
	})
	return resultSlice, nil
}

func (s *Service) ticker(ctx context.Context) {
	s.log.Info("Ticker started")
	ticker := time.NewTicker(s.indexTTL)
	s.reindex(ctx)
	for {
		<-ticker.C
		s.log.Debug("Reindex started")
		s.reindex(ctx)
	}
}

func (s *Service) reindex(ctx context.Context) {
	if !s.mu.TryLock() {
		s.log.Warn("Failed to lock for reindex")
		return
	}
	defer s.mu.Unlock()

	comics, err := s.db.GetAll(ctx)
	if err != nil {
		s.log.Error("Failed to get comics for reindex")
		return
	}

	type wordComic struct {
		word  string
		comic Comics
	}
	resChan := make(chan wordComic)

	sema := make(chan struct{}, s.numWorkers)
	var wg sync.WaitGroup
	wg.Add(len(comics))
	for _, comic := range comics {
		go func() {
			sema <- struct{}{}
			for _, word := range comic.Words {
				resChan <- wordComic{word: word, comic: comic}
			}
			wg.Done()
			<-sema
		}()
	}

	go func() {
		wg.Wait()
		close(resChan)
	}()

	s.log.Info("Reindex starting", "size of indexes words before", len(s.indComics))

	for wi := range resChan {

		if s.indComics[wi.word] == nil {
			s.indComics[wi.word] = make(comicsHashSet)
		}
		s.indComics[wi.word][wi.comic.ID] = wi.comic
	}

	s.log.Info("Reindex finished", "size of indexes words", len(s.indComics))
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
