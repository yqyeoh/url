package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/yqyeoh/url/app"
	"github.com/yqyeoh/url/mocks"
	"go.uber.org/zap"
)

func getLogger() *zap.SugaredLogger {
	logger, _ := zap.NewDevelopment()
	log := logger.Sugar()
	return log
}
func TestAppService_Create(t *testing.T) {
	t.Run("Returns url string when succeeded", func(t *testing.T) {
		logger := getLogger()

		repo := &mocks.AppRepo{}
		repo.On("Create", mock.Anything, mock.Anything, mock.Anything).Return("abc", nil)

		service := app.NewService(logger, repo)
		url, err := service.Create(context.Background(), "xyz")

		require.NoError(t, err)
		require.Equal(t, url, "abc")
	})
}
