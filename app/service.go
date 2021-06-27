package app

import (
	"context"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

const codeLength = 7

type Service interface {
	Create(ctx context.Context, url string) (string, error)
	FindURLByCode(ctx context.Context, code string) (string, error)
}

type service struct {
	repo   Repo
	logger *zap.SugaredLogger
}

func NewService(logger *zap.SugaredLogger, repo Repo) Service {
	return service{
		repo,
		logger,
	}
}

func (s service) Create(ctx context.Context, url string) (string, error) {
	return s.repo.Create(ctx, randomCode(), url)
}

func (s service) FindURLByCode(ctx context.Context, code string) (string, error) {
	return s.repo.FindURLByCode(ctx, code)
}

func randomCode() string {
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, codeLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
