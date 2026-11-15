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
		Variants: []*models.Variant{
			{
				ID:           variantID,
				ExperimentID: experimentID,
				StrategyID:   100,
				AssetID:      1,
				Asset:        asset,
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

	mockExecutionsRepo.EXPECT().CreateExecutionVariant(executionID, variantID, map[string]any{"threshold": 0.05}).Return(&models.ExecutionVariant{
		ID:          executionVariantID,
		ExecutionID: executionID,
		VariantID:   variantID,
		Params:      models.NewJsonColumn(map[string]any{"threshold": 0.05}),
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
		ID:       executionVariantID,
		Status:   "COMPLETED",
		Attempts: 1,
	}, nil)

	mockExecutionsRepo.EXPECT().UpdateExecution(executionID, "COMPLETED", models.Nullable[string]{}).Return(&models.Execution{
		ID:     executionID,
		Status: "COMPLETED",
	}, nil)

	err := service.Execute(experimentID)
	assert.NoError(t, err)
}

func TestExecutionService_Execute_VariantFailureWithAttempts(t *testing.T) {
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

	experimentID := 3
	asset := "BTCUSDT"
	variantID := 8
	executionID := 101
	executionVariantID := 201

	hydratedExperiment := &models.Experiment{
		ID:     experimentID,
		Name:   "Clockwork_BTC_2H_WalkForward",
		Status: "ACTIVE",
		Variants: []*models.Variant{
			{
				ID:           variantID,
				ExperimentID: experimentID,
				StrategyID:   3,
				AssetID:      8,
				Asset:        asset,
				Name:         "BTCUSDT - Aug Fixed Postfit",
				StrategyName: "Clockwork HGB Regressor",
				Class:        "clockwork_hgb",
				MaxAttempts:  3,
				Params:       models.NewJsonColumn(map[string]any{"current_price": 63500.0}),
			},
		},
	}

	mockExperimentsRepo.EXPECT().FindByID(experimentID).Return(hydratedExperiment, nil)
	mockExecutionsRepo.EXPECT().CreateExecution(experimentID, "MANUAL", "system", mock.AnythingOfType("time.Time")).Return(&models.Execution{
		ID:           executionID,
		ExperimentID: experimentID,
		Status:       "RUNNING",
	}, nil)

	mockExecutionsRepo.EXPECT().CreateExecutionVariant(executionID, variantID, map[string]any{"current_price": 63500.0}).Return(&models.ExecutionVariant{
		ID:          executionVariantID,
		ExecutionID: executionID,
		VariantID:   variantID,
		Params:      models.NewJsonColumn(map[string]any{"current_price": 63500.0}),
		Status:      "RUNNING",
	}, nil)

	totalAttempts := services.DefaultMaxAttempts
	mockStrategyFactory.EXPECT().CreateStrategy("clockwork_hgb").Return(mockStrategy, nil).Times(totalAttempts)
	mockExecutionMutex.EXPECT().Lock(experimentID, variantID, asset).Return().Times(totalAttempts)
	mockExecutionMutex.EXPECT().Unlock(experimentID, variantID, asset).Return().Times(totalAttempts)
	mockStrategy.EXPECT().GenerateSignal(asset, mock.AnythingOfType("map[string]interface {}")).Return(nil, errors.New("404 Not Found")).Times(totalAttempts)

	mockExecutionsRepo.EXPECT().UpdateExecutionVariant(
		executionVariantID,
		"FAILED",
		services.DefaultMaxAttempts,
		mock.AnythingOfType("int"),
		mock.MatchedBy(func(msg models.Nullable[string]) bool {
			return msg.Valid && msg.Val != ""
		}),
	).Return(&models.ExecutionVariant{
		ID:       executionVariantID,
		Status:   "FAILED",
		Attempts: services.DefaultMaxAttempts,
	}, nil)

	mockExecutionsRepo.EXPECT().UpdateExecution(
		executionID,
		"FAILED",
		models.NotNull("All variants failed during execution"),
	).Return(&models.Execution{
		ID:     executionID,
		Status: "FAILED",
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
