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
		AssetID:      100,
		Type:         models.SignalBuy,
		Price:        150.50,
		Confidence:   0.95,
		StrategyName: "FiftyDayAvg",
	}

	expectedResult := models.SignalRecord{
		ID:           1,
		VariantID:    2,
		ExperimentID: 3,
		StrategyID:   4,
		AssetID:      100,
		Signal:       models.SignalBuy,
		Price:        150.50,
		Confidence:   0.95,
		CreatedAt:    time.Now(),
	}

	mockDb.EXPECT().
		QueryOne("CALL sp_save_signal (?, ?, ?, ?, ?)", mock.AnythingOfType("*models.SignalRecord"), []any{100, 2, "BUY", 150.50, 0.95}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*models.SignalRecord)
			*ptr = expectedResult
			return nil
		})

	result, err := repo.Save(100, 2, signal)

	assert.NoError(t, err)
	assert.Equal(t, &expectedResult, result)
}

func TestSignalsRepository_Save_Error(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	signal := &models.Signal{
		AssetID: 100,
		Type:    models.SignalBuy,
		Price:   150.50,
	}

	mockDb.EXPECT().
		QueryOne("CALL sp_save_signal (?, ?, ?, ?, ?)", mock.AnythingOfType("*models.SignalRecord"), []any{100, 2, "BUY", 150.50, 1.0}).
		Return(errors.New("db error"))

	result, err := repo.Save(100, 2, signal)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestSignalsRepository_Find_Success(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	now := time.Now()
	expectedResults := []*models.SignalRecord{
		{
			ID:      1,
			AssetID: 100,
			Signal:  models.SignalBuy,
		},
	}

	fromDate := models.NotNull(now.Add(-24 * time.Hour))
	toDate := models.NotNull(now)

	mockDb.EXPECT().
		QueryAll("CALL sp_find_signals (?, ?, ?)", mock.AnythingOfType("*[]*models.SignalRecord"), []any{100, fromDate.Val, toDate.Val}).
		RunAndReturn(func(query string, dest any, args ...any) error {
			ptr := dest.(*[]*models.SignalRecord)
			*ptr = expectedResults
			return nil
		})

	results, err := repo.Find(100, fromDate, toDate)

	assert.NoError(t, err)
	assert.Equal(t, expectedResults, results)
}

func TestSignalsRepository_Find_NullDates(t *testing.T) {
	mockDb := mocks.NewMockIDBHandler(t)
	repo := NewSignalsRepository(mockDb)

	fromDate := models.Null[time.Time]()
	toDate := models.Null[time.Time]()

	mockDb.EXPECT().
		QueryAll("CALL sp_find_signals (?, ?, ?)", mock.AnythingOfType("*[]*models.SignalRecord"), []any{100, nil, nil}).
		Return(nil)

	results, err := repo.Find(100, fromDate, toDate)

	assert.NoError(t, err)
	assert.Empty(t, results)
}
