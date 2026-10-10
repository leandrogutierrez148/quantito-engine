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
	mockExecutionsRepo := mocks.NewMockIExecutionsRepository(t)
	mockSignalsRepo := mocks.NewMockISignalsRepository(t)
	mockStrategyFactory := mocks.NewMockIStrategyFactory(t)
	mockSignalDispatcher := mocks.NewMockISignalDispatcher(t)
	mockExecutionMutex := mocks.NewMockIExecutionMutex(t)
	mockStrategy := mocks.NewMockIStrategy(t)

	service := services.NewExecutionService(
		mockExperimentsRepo,
		mockExecutionsRepo,
		mockSignalsRepo,
		mockStrategyFactory,
		mockSignalDispatcher,
		mockExecutionMutex,
	)

	experimentID := 1
	asset := "AAPL"
	variantID := 10
	executionID := 100
	executionVariantID := 200

	hydratedExperiment := &models.Experiment{
		ID:     experimentID,
		Name:   "Tech_Momentum",
		Status: "ACTIVE",
		Assets: models.NewJsonColumn([]string{asset}),
		Variants: []*models.Variant{
			{
				ID:           variantID,
				ExperimentID: experimentID,
				StrategyID:   100,
				Name:         "Variant A",
				StrategyName: "Fifty Day Avg",
				Class:        "fifty_day_avg",
				Params:       models.NewJsonColumn(map[string]any{"threshold": 0.05}),
			},
		},
	}

	mockExperimentsRepo.EXPECT().FindByID(experimentID).Return(hydratedExperiment, nil)
	mockExecutionsRepo.EXPECT().CreateExecution(experimentID, "MANUAL", "system", mock.AnythingOfType("time.Time")).Return(&models.Execution{
		ID:           executionID,
		ExperimentID: experimentID,
		Status:       "RUNNING",
	}, nil)

	mockExecutionsRepo.EXPECT().CreateExecutionVariant(executionID, variantID).Return(&models.ExecutionVariant{
		ID:          executionVariantID,
		ExecutionID: executionID,
		VariantID:   variantID,
		Status:      "RUNNING",
	}, nil)

	mockExecutionMutex.EXPECT().Lock(experimentID, variantID, asset).Return()
	mockExecutionMutex.EXPECT().Unlock(experimentID, variantID, asset).Return()

	mockStrategyFactory.EXPECT().CreateStrategy("fifty_day_avg").Return(mockStrategy, nil)

	expectedSignal := models.NewSignal(asset, models.SignalBuy, 150.0, "fifty_day_avg")
	mockStrategy.EXPECT().GenerateSignal(asset, mock.AnythingOfType("map[string]interface {}")).Return(expectedSignal, nil)

	mockSignalsRepo.EXPECT().Save(executionVariantID, asset, expectedSignal).Return(&models.SignalRecord{
		ID:                 55,
		ExecutionVariantID: executionVariantID,
		Asset:              asset,
		Signal:             models.SignalBuy,
		Price:              150.0,
	}, nil)

	mockSignalDispatcher.EXPECT().Dispatch(expectedSignal).Return(nil)

	mockExecutionsRepo.EXPECT().UpdateExecutionVariant(executionVariantID, "COMPLETED", 1, mock.AnythingOfType("int"), models.Nullable[string]{}).Return(&models.ExecutionVariant{
		ID:     executionVariantID,
		Status: "COMPLETED",
	}, nil)

	mockExecutionsRepo.EXPECT().UpdateExecution(executionID, "COMPLETED", models.Nullable[string]{}).Return(&models.Execution{
		ID:     executionID,
		Status: "COMPLETED",
	}, nil)

	err := service.Execute(experimentID)
	assert.NoError(t, err)
}

func TestExecutionService_Execute_InactiveExperiment(t *testing.T) {
	mockExperimentsRepo := mocks.NewMockIExperimentsRepository(t)
	mockExecutionsRepo := mocks.NewMockIExecutionsRepository(t)
	mockSignalsRepo := mocks.NewMockISignalsRepository(t)
	mockStrategyFactory := mocks.NewMockIStrategyFactory(t)
	mockSignalDispatcher := mocks.NewMockISignalDispatcher(t)
	mockExecutionMutex := mocks.NewMockIExecutionMutex(t)

	service := services.NewExecutionService(
		mockExperimentsRepo,
		mockExecutionsRepo,
		mockSignalsRepo,
		mockStrategyFactory,
		mockSignalDispatcher,
		mockExecutionMutex,
	)

	experimentID := 999
	mockExperimentsRepo.EXPECT().FindByID(experimentID).Return(&models.Experiment{
		ID:     experimentID,
		Name:   "Paused_Exp",
		Status: "PAUSED",
	}, nil)

	err := service.Execute(experimentID)
	assert.NoError(t, err)
}

func TestExecutionService_Execute_ExperimentNotFoundOrError(t *testing.T) {
	mockExperimentsRepo := mocks.NewMockIExperimentsRepository(t)
	mockExecutionsRepo := mocks.NewMockIExecutionsRepository(t)
	mockSignalsRepo := mocks.NewMockISignalsRepository(t)
	mockStrategyFactory := mocks.NewMockIStrategyFactory(t)
	mockSignalDispatcher := mocks.NewMockISignalDispatcher(t)
	mockExecutionMutex := mocks.NewMockIExecutionMutex(t)

	service := services.NewExecutionService(
		mockExperimentsRepo,
		mockExecutionsRepo,
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
