package repositories

import (
	"errors"
	"testing"

	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExperimentAssetsRepositoryFindByExperiment(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewExperimentAssetsRepository(mockDb)

	expectedResult := []*models.ExperimentAsset{
		{
			ID:           1,
			ExperimentID: 1,
			AssetID:      101,
			Status:       "ACTIVE",
		},
	}

	mockDb.EXPECT().
		QueryAll("CALL sp_find_experiment_assets (?)", mock.AnythingOfType("*[]*models.ExperimentAsset"), []any{1}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.ExperimentAsset)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.FindByExperiment(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestExperimentAssetsRepositoryFindExperimentsByAsset(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewExperimentAssetsRepository(mockDb)

	expectedResult := []*models.Experiment{
		{
			ID:         1,
			Experiment: "Test Experiment",
			Status:     "ACTIVE",
		},
	}

	mockDb.EXPECT().
		QueryAll("CALL sp_find_experiments_by_asset (?)", mock.AnythingOfType("*[]*models.Experiment"), []any{101}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.Experiment)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.FindExperimentsByAsset(101)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)

	// Error case
	mockDb.EXPECT().
		QueryAll("CALL sp_find_experiments_by_asset (?)", mock.AnythingOfType("*[]*models.Experiment"), []any{999}).
		Return(errors.New("db error"))

	result, err = repo.FindExperimentsByAsset(999)
	assert.Error(t, err)
	assert.Nil(t, result)
}
