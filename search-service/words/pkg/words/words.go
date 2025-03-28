package words

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
	"github.com/kljensen/snowball/english"
	"github.com/kljensen/snowball/french"
	"github.com/kljensen/snowball/hungarian"
	"github.com/kljensen/snowball/russian"
	"github.com/kljensen/snowball/spanish"
	"github.com/kljensen/snowball/swedish"

	"github.com/kljensen/snowball/norwegian"
	"yadro.com/course/words/pkg/logger"
	hashset "yadro.com/course/words/pkg/structure"
)

var log = logger.GetInstance()

func StemPhrase(ctx context.Context, phrase string, language string) ([]string, error) {

	isStopWordFunc, err := getFuncIsStopWord(language)
	if err != nil {
		log.Warn("Error while searching function for stopping word", "err", err)
		return nil, err
	}

	phrase = strings.Trim(phrase, " ")

	normalizedWords := hashset.New[string]()
	replacedPuncts := ReplacePuncts(phrase)

	words := strings.Fields(replacedPuncts)

	for _, word := range words {
		select {
		case <-ctx.Done():
			log.Info("Ctx done")
			return nil, fmt.Errorf("ctx is done, finishing normalizing")
		default:
		}
		if ctx.Err() != nil {
			log.Error("Context error", "err", err)
			return nil, fmt.Errorf("context error before processing: %w", ctx.Err())
		}

		if isStopWordFunc(word) {
			continue
		}

		if normalizedWords.Contains(word) {
			continue
		}

		stemmedPhrase, err := snowball.Stem(word, language, true) //NOTE: third parameter doesn't matter cause we skip stop words
		if err != nil {
			log.Warn("Error while searching function for stopping word", "err", err)
			return nil, err
		}
		normalizedWords.Add(stemmedPhrase)
	}

	return normalizedWords.ToSlice(), nil
}

func GetWords(_ context.Context, input string) ([]string, error) {
	return strings.Fields(input), nil
}

func ReplacePuncts(input string) string {

	var res strings.Builder
	for _, r := range input {
		if unicode.IsPunct(r) {
			res.WriteRune(' ')
			continue
		}
		res.WriteRune(unicode.ToLower(r))
	}
	return res.String()
}

// NOTE: можно сделать нормально, например через мапу, но это было подрезано по примеру snowball.Stem и переделывать немного лень:D
func getFuncIsStopWord(language string) (func(string) bool, error) {
	var f func(string) bool
	switch language {
	case "english":
		f = english.IsStopWord
	case "spanish":
		f = spanish.IsStopWord
	case "french":
		f = french.IsStopWord
	case "russian":
		f = russian.IsStopWord
	case "swedish":
		f = swedish.IsStopWord
	case "norwegian":
		f = norwegian.IsStopWord
	case "hungarian":
		f = hungarian.IsStopWord
	default:
		log.Warn("Language not found", "given language", language)
		err := fmt.Errorf("language not found")
		return nil, err
	}

	return f, nil
}
