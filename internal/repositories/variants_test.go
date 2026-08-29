package repositories

import (
	"errors"
	"testing"

	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVariantsRepository_FindByExperiment_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewVariantsRepository(mockDb)

	expectedResult := []*models.Variant{
		{
			ID:           1,
			ExperimentID: 1,
			StrategyID:   10,
			Variant:      "Variant A",
			StrategyName: "Fifty Day Avg",
			Class:        "fifty_day_avg",
		},
	}

	mockDb.EXPECT().
		QueryAll("CALL sp_find_variants_by_experiment (?)", mock.AnythingOfType("*[]*models.Variant"), []any{1}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.Variant)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.FindByExperiment(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, result)
}

func TestVariantsRepository_FindByExperiment_Error(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewVariantsRepository(mockDb)

	mockDb.EXPECT().
		QueryAll("CALL sp_find_variants_by_experiment (?)", mock.AnythingOfType("*[]*models.Variant"), []any{42}).
		Return(errors.New("db error"))

	result, err := repo.FindByExperiment(42)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "db error", err.Error())
}
