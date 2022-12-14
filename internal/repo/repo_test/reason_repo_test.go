package repo_test

import (
	"context"
	"testing"

	"answer/internal/repo/config"
	"answer/internal/repo/reason"

	"github.com/stretchr/testify/assert"
)

func Test_reasonRepo_ListReasons(t *testing.T) {
	configRepo := config.NewConfigRepo(testDataSource)
	reasonRepo := reason.NewReasonRepo(configRepo)
	reasonItems, err := reasonRepo.ListReasons(context.TODO(), "question", "close")
	assert.NoError(t, err)
	assert.Equal(t, 4, len(reasonItems))
}
