package service

import (
	"context"
	"strings"

	"yadro.com/course/words/pkg/logger"
	"yadro.com/course/words/pkg/words"
)

var log = logger.GetInstance()

const defaultLanguage = "english"

type NormalizeService interface {
	Normalize(ctx context.Context, phrase string, language string) ([]string, error)
	GetWords(ctx context.Context, phrase string) ([]string, error)
}

type normalizeSerive struct{}

func NewNormalizeService() NormalizeService {
	return &normalizeSerive{}
}

// NOTE: забавный выстрел в ногу самому себе был выполнен, но в целом так даже логичнее)
func (n *normalizeSerive) Normalize(ctx context.Context, phrase string, language string) ([]string, error) {
	log.Debug("Normalize method called")
	if strings.Trim(language, " ") == "" {
		language = defaultLanguage
	}
	return words.StemPhrase(ctx, phrase, language)
}

func (n *normalizeSerive) GetWords(ctx context.Context, phrase string) ([]string, error) {
	log.Debug("GetWords method called")
	return words.GetWords(ctx, phrase)
}
