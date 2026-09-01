package services_test

import (
	"errors"
	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"
	"quantito-engine/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExecutionService_Execute_Success(t *testing.T) {
	mockExperimentsRepo := mocks.NewMockIExperimentsRepository(t)
	mockSignalsRepo := mocks.NewMockISignalsRepository(t)
	mockStrategyFactory := mocks.NewMockIStrategyFactory(t)
	mockSignalDispatcher := mocks.NewMockISignalDispatcher(t)
	mockExecutionMutex := mocks.NewMockIExecutionMutex(t)
	mockStrategy := mocks.NewMockIStrategy(t)

	service := services.NewExecutionService(
		mockExperimentsRepo,
		mockSignalsRepo,
		mockStrategyFactory,
		mockSignalDispatcher,
		mockExecutionMutex,
	)

	experimentID := 1
	assetID := 101
	variantID := 10

	hydratedExperiment := &models.Experiment{
		ID:         experimentID,
		Experiment: "Tech_Momentum",
		Status:     "ACTIVE",
		Assets: []*models.ExperimentAsset{
			{ID: 1, ExperimentID: experimentID, AssetID: assetID, Status: "ACTIVE"},
		},
		Variants: []*models.Variant{
			{
				ID:           variantID,
				ExperimentID: experimentID,
				StrategyID:   100,
				Variant:      "A",
				StrategyName: "Fifty Day Avg",
				Class:        "fifty_day_avg",
				Params:       models.NewJsonColumn(map[string]any{"threshold": 0.05}),
			},
		},
	}

	mockExperimentsRepo.EXPECT().FindByID(experimentID).Return(hydratedExperiment, nil)
	mockExecutionMutex.EXPECT().Lock(experimentID, variantID, assetID).Return()
	mockExecutionMutex.EXPECT().Unlock(experimentID, variantID, assetID).Return()

	mockStrategyFactory.EXPECT().CreateStrategy("fifty_day_avg").Return(mockStrategy, nil)

	expectedSignal := models.NewSignal(assetID, models.SignalBuy, 150.0, "fifty_day_avg")
	mockStrategy.EXPECT().GenerateSignal(assetID, mock.AnythingOfType("map[string]interface {}")).Return(expectedSignal, nil)

	mockSignalsRepo.EXPECT().Save(assetID, variantID, expectedSignal).Return(&models.SignalRecord{
		ID:           55,
		VariantID:    variantID,
		ExperimentID: experimentID,
		AssetID:      assetID,
		Signal:       models.SignalBuy,
		Price:        150.0,
	}, nil)

	mockSignalDispatcher.EXPECT().Dispatch(expectedSignal).Return(nil)

	err := service.Execute(experimentID)
	assert.NoError(t, err)
}

func TestExecutionService_Execute_InactiveExperiment(t *testing.T) {
	mockExperimentsRepo := mocks.NewMockIExperimentsRepository(t)
	mockSignalsRepo := mocks.NewMockISignalsRepository(t)
	mockStrategyFactory := mocks.NewMockIStrategyFactory(t)
	mockSignalDispatcher := mocks.NewMockISignalDispatcher(t)
	mockExecutionMutex := mocks.NewMockIExecutionMutex(t)

	service := services.NewExecutionService(
		mockExperimentsRepo,
		mockSignalsRepo,
		mockStrategyFactory,
		mockSignalDispatcher,
		mockExecutionMutex,
	)

	experimentID := 999
	mockExperimentsRepo.EXPECT().FindByID(experimentID).Return(&models.Experiment{
		ID:         experimentID,
		Experiment: "Paused_Exp",
		Status:     "PAUSED",
	}, nil)

	err := service.Execute(experimentID)
	assert.NoError(t, err)
}

func TestExecutionService_Execute_ExperimentNotFoundOrError(t *testing.T) {
	mockExperimentsRepo := mocks.NewMockIExperimentsRepository(t)
	mockSignalsRepo := mocks.NewMockISignalsRepository(t)
	mockStrategyFactory := mocks.NewMockIStrategyFactory(t)
	mockSignalDispatcher := mocks.NewMockISignalDispatcher(t)
	mockExecutionMutex := mocks.NewMockIExecutionMutex(t)

	service := services.NewExecutionService(
		mockExperimentsRepo,
		mockSignalsRepo,
		mockStrategyFactory,
		mockSignalDispatcher,
		mockExecutionMutex,
	)

	experimentID := 404
	mockExperimentsRepo.EXPECT().FindByID(experimentID).Return(nil, errors.New("db error"))

	err := service.Execute(experimentID)
	assert.Error(t, err)
}
