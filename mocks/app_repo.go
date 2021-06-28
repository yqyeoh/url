package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type AppRepo struct {
	mock.Mock
}

func (m *AppRepo) Create(ctx context.Context, randomCode, url string) (string, error) {
	args := m.Called(ctx, randomCode, url)
	return args.String(0), args.Error(1)
}

func (m *AppRepo) FindURLByCode(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}
