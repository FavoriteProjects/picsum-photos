package health_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/DMarby/picsum-photos/internal/health"
	"github.com/DMarby/picsum-photos/internal/logger"
	"go.uber.org/zap"

	fileDatabase "github.com/DMarby/picsum-photos/internal/database/file"
	mockDatabase "github.com/DMarby/picsum-photos/internal/database/mock"

	fileStorage "github.com/DMarby/picsum-photos/internal/storage/file"
	mockStorage "github.com/DMarby/picsum-photos/internal/storage/mock"

	memoryCache "github.com/DMarby/picsum-photos/internal/cache/memory"
	mockCache "github.com/DMarby/picsum-photos/internal/cache/mock"
)

func TestHealth(t *testing.T) {
	log := logger.New(zap.ErrorLevel)
	defer log.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage, _ := fileStorage.New("../../test/fixtures/file")
	db, _ := fileDatabase.New("../../test/fixtures/file/metadata.json")
	cache := memoryCache.New()

	checker := &health.Checker{Ctx: ctx, Storage: storage, Database: db, Cache: cache, Log: log}
	mockStorageChecker := &health.Checker{Ctx: ctx, Storage: &mockStorage.Provider{}, Database: db, Cache: cache, Log: log}
	mockDatabaseChecker := &health.Checker{Ctx: ctx, Storage: storage, Database: &mockDatabase.Provider{}, Cache: cache, Log: log}
	mockCacheChecker := &health.Checker{Ctx: ctx, Storage: storage, Database: db, Cache: &mockCache.Provider{}, Log: log}

	dbOnlyChecker := &health.Checker{Ctx: ctx, Database: db, Log: log}
	mockDbOnlyChecker := &health.Checker{Ctx: ctx, Database: &mockDatabase.Provider{}, Log: log}

	tests := []struct {
		Name           string
		ExpectedStatus health.Status
		Checker        *health.Checker
	}{
		{
			Name: "runs checks and returns correct status",
			ExpectedStatus: health.Status{
				Healthy:  true,
				Cache:    "healthy",
				Database: "healthy",
				Storage:  "healthy",
			},
			Checker: checker,
		},
		{
			Name: "runs checks and returns correct status with broken storage",
			ExpectedStatus: health.Status{
				Healthy:  false,
				Cache:    "healthy",
				Database: "healthy",
				Storage:  "unhealthy",
			},
			Checker: mockStorageChecker,
		},
		{
			Name: "runs checks and returns correct status with broken database",
			ExpectedStatus: health.Status{
				Healthy:  false,
				Cache:    "healthy",
				Database: "unhealthy",
				Storage:  "unknown",
			},
			Checker: mockDatabaseChecker,
		},
		{
			Name: "runs checks and returns correct status with broken cache",
			ExpectedStatus: health.Status{
				Healthy:  false,
				Cache:    "unhealthy",
				Database: "healthy",
				Storage:  "healthy",
			},
			Checker: mockCacheChecker,
		},
		{
			Name: "runs checks and returns correct status with only a database",
			ExpectedStatus: health.Status{
				Healthy:  true,
				Database: "healthy",
			},
			Checker: dbOnlyChecker,
		},
		{
			Name: "runs checks and returns correct status with only a broken database",
			ExpectedStatus: health.Status{
				Healthy:  false,
				Database: "unhealthy",
			},
			Checker: mockDbOnlyChecker,
		},
	}

	for _, test := range tests {
		test.Checker.Run()
		status := test.Checker.Status()

		if !reflect.DeepEqual(status, test.ExpectedStatus) {
			t.Errorf("%s: wrong status %+v", test.Name, status)
		}
	}

	t.Run("checker runs and returns correct status", func(t *testing.T) {

	})
}
