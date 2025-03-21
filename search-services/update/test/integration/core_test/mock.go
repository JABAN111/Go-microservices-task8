package core_test

import (
	"context"

	"github.com/stretchr/testify/mock"
	"yadro.com/course/update/core"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Add(ctx context.Context, comics core.Comics) error {
	args := m.Called(ctx, comics)
	return args.Error(0)
}

func (m *MockDB) UpdateStats(ctx context.Context, cntUniqueWords int, comicsInTotal int) error {
	args := m.Called(ctx, cntUniqueWords, comicsInTotal)
	return args.Error(0)
}

func (m *MockDB) AddAllComics(ctx context.Context, comics <-chan core.Comics) error {
	args := m.Called(ctx, comics)
	return args.Error(0)
}

func (m *MockDB) AddWordStats(ctx context.Context, wordsList []string) error {
	args := m.Called(ctx, wordsList)
	return args.Error(0)
}

func (m *MockDB) DbStats(ctx context.Context) (core.DBStats, error) {
	args := m.Called(ctx)
	return args.Get(0).(core.DBStats), args.Error(1)
}

func (m *MockDB) IDs(ctx context.Context) ([]int, error) {
	args := m.Called(ctx)
	return args.Get(0).([]int), args.Error(1)
}

type MockXKCD struct {
	mock.Mock
}

func (m *MockXKCD) Get(ctx context.Context, id int) (core.XKCDInfo, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(core.XKCDInfo), args.Error(1)
}

func (m *MockDB) Drop(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockXKCD) LastID(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

type MockWords struct {
	mock.Mock
}

func (m *MockWords) Norm(ctx context.Context, phrase string) ([]string, error) {
	args := m.Called(ctx, phrase)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockWords) GetWords(ctx context.Context, phrase string) ([]string, error) {
	args := m.Called(ctx, phrase)
	return args.Get(0).([]string), args.Error(1)
}
