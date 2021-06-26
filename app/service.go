package app

import (
	"context"
	"math/rand"

	"go.uber.org/zap"
)

const codeLength = 7

type Service interface {
	FindOrCreateCode(ctx context.Context, url string) (string, error)
}

type service struct {
	repo   Repo
	logger *zap.SugaredLogger
}

func NewService(repo Repo, logger *zap.SugaredLogger) Service {
	return service{
		repo,
		logger,
	}
}

func (s service) FindOrCreateCode(ctx context.Context, url string) (string, error) {

	return s.repo.FindOrCreateCode(ctx, randomCode(), url)
}

func randomCode() string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, codeLength)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
