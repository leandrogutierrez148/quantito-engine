package repositories

import (
	"errors"
	"testing"
	"time"

	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignalsRepository_Save_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	signal := &models.Signal{
		Asset:        "AAPL",
		Type:         models.SignalBuy,
		Price:        150.50,
		Confidence:   0.95,
		StrategyName: "FiftyDayAvg",
	}

	expectedResult := models.SignalRecord{
		ID:                 1,
		ExecutionVariantID: 10,
		Asset:              "AAPL",
		Signal:             models.SignalBuy,
		Price:              150.50,
		Confidence:         0.95,
		CreatedAt:          time.Now(),
	}

	mockDb.EXPECT().
		QueryOne("CALL sp_save_signal (?, ?, ?, ?, ?)", mock.AnythingOfType("*models.SignalRecord"), []any{10, "AAPL", "BUY", 150.50, 0.95}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*models.SignalRecord)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.Save(10, "AAPL", signal)

	assert.NoError(t, err)
	assert.Equal(t, &expectedResult, result)
}

func TestSignalsRepository_Save_Error(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	signal := &models.Signal{
		Asset: "AAPL",
		Type:  models.SignalBuy,
		Price: 150.50,
	}

	mockDb.EXPECT().
		QueryOne("CALL sp_save_signal (?, ?, ?, ?, ?)", mock.AnythingOfType("*models.SignalRecord"), []any{10, "AAPL", "BUY", 150.50, 1.0}).
		Return(errors.New("db error"))

	result, err := repo.Save(10, "AAPL", signal)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSignalsRepository_FindByExecution_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	expectedResults := []*models.SignalRecord{
		{
			ID:                 1,
			ExecutionVariantID: 10,
			Asset:              "AAPL",
			Signal:             models.SignalBuy,
		},
	}

	mockDb.EXPECT().
		QueryAll("CALL sp_find_signals_by_execution (?)", mock.AnythingOfType("*[]*models.SignalRecord"), []any{100}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.SignalRecord)
			*ptr = expectedResults
			return nil
		})

	results, err := repo.FindByExecution(100)

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, results)
}

func TestSignalsRepository_FindByAsset_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	now := time.Now()
	expectedResults := []*models.SignalRecord{
		{
			ID:                 1,
			ExecutionVariantID: 10,
			Asset:              "AAPL",
			Signal:             models.SignalBuy,
		},
	}

	fromDate := models.NotNull(now.Add(-24 * time.Hour))
	toDate := models.NotNull(now)

	mockDb.EXPECT().
		QueryAll("CALL sp_find_signals_by_asset (?, ?, ?)", mock.AnythingOfType("*[]*models.SignalRecord"), []any{"AAPL", fromDate.Val, toDate.Val}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.SignalRecord)
			*ptr = expectedResults
			return nil
		})

	results, err := repo.FindByAsset("AAPL", fromDate, toDate)

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, results)
}
