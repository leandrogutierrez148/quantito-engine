package repositories

import (
	"errors"
	"testing"

	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExperimentsRepository_Find_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewExperimentsRepository(mockDb)

	expectedExperiments := []*models.Experiment{
		{
			ID:     1,
			Name:   "Tech Momentum",
			Status: "ACTIVE",
			Assets: models.NewJsonColumn([]string{"AAPL", "MSFT"}),
		},
	}

	mockDb.EXPECT().
		QueryAll("CALL sp_find_experiments (?)", mock.AnythingOfType("*[]*models.Experiment"), []any{false}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.Experiment)
			*ptr = expectedExperiments
			return nil
		})

	result, err := repo.Find(false)

	assert.NoError(t, err)
	assert.Equal(t, expectedExperiments, result)
}

func TestExperimentsRepository_Find_Error(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewExperimentsRepository(mockDb)

	mockDb.EXPECT().
		QueryAll("CALL sp_find_experiments (?)", mock.AnythingOfType("*[]*models.Experiment"), []any{true}).
		Return(errors.New("db error"))

	result, err := repo.Find(true)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestExperimentsRepository_FindByID_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewExperimentsRepository(mockDb)

	baseExp := models.Experiment{
		ID:     1,
		Name:   "Tech Momentum",
		Status: "ACTIVE",
		Assets: models.NewJsonColumn([]string{"AAPL", "MSFT"}),
	}

	expectedVariants := []*models.Variant{
		{
			ID:           1,
			ExperimentID: 1,
			StrategyID:   10,
			Name:         "Variant A",
			StrategyName: "Fifty Day Avg",
			Class:        "fifty_day_avg",
		},
	}

	mockDb.EXPECT().
		QueryOne("CALL sp_get_experiment_by_id (?)", mock.AnythingOfType("*models.Experiment"), []any{1}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*models.Experiment)
			*ptr = baseExp
			return nil
		})

	mockDb.EXPECT().
		QueryAll("CALL sp_find_variants_by_experiment (?)", mock.AnythingOfType("*[]*models.Variant"), []any{1}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.Variant)
			*ptr = expectedVariants
			return nil
		})

	result, err := repo.FindByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, []string{"AAPL", "MSFT"}, result.GetAssets())
	assert.Equal(t, expectedVariants, result.Variants)
}

func TestExperimentsRepository_FindByExperiment_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewExperimentsRepository(mockDb)

	baseExp := models.Experiment{
		ID:     1,
		Name:   "Tech Momentum",
		Status: "ACTIVE",
		Assets: models.NewJsonColumn([]string{"AAPL"}),
	}

	expectedVariants := []*models.Variant{
		{
			ID:           1,
			ExperimentID: 1,
			StrategyID:   10,
			Name:         "Variant A",
			StrategyName: "Fifty Day Avg",
			Class:        "fifty_day_avg",
		},
	}

	mockDb.EXPECT().
		QueryOne("CALL sp_get_experiment_by_name (?)", mock.AnythingOfType("*models.Experiment"), []any{"Tech Momentum"}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*models.Experiment)
			*ptr = baseExp
			return nil
		})

	mockDb.EXPECT().
		QueryAll("CALL sp_find_variants_by_experiment (?)", mock.AnythingOfType("*[]*models.Variant"), []any{1}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.Variant)
			*ptr = expectedVariants
			return nil
		})

	result, err := repo.FindByExperiment("Tech Momentum")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, expectedVariants, result.Variants)
}
