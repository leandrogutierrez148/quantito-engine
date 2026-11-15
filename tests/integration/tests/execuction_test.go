package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"quantito-engine/internal/controllers"
	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"
	"quantito-engine/internal/repositories"
	"quantito-engine/internal/services"
	"quantito-engine/internal/sinks"
	"quantito-engine/internal/strategies"
	testdb "quantito-engine/tests/integration/db"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupIntegrationApp(t *testing.T) (*echo.Echo, interfaces.IExperimentsRepository, interfaces.IExecutionsRepository, interfaces.ISignalsRepository, *mocks.MockISignalSink, interfaces.IDBHandler) {
	dbHandler := testdb.NewTestDBHandler(t)

	experimentsRepo := repositories.NewExperimentsRepository(dbHandler)
	executionsRepo := repositories.NewExecutionsRepository(dbHandler)
	signalsRepo := repositories.NewSignalsRepository(dbHandler)

	strategyFactory := strategies.NewStrategyFactory()
	mockSink := mocks.NewMockISignalSink(t)
	consoleSink := sinks.NewConsoleSink()
	signalDispatcher := sinks.NewSignalDispatcher(consoleSink, mockSink)
	executionMutex := infraestructures.NewExecutionMutex()

	executionService := services.NewExecutionService(
		experimentsRepo,
		executionsRepo,
		signalsRepo,
		strategyFactory,
		signalDispatcher,
		executionMutex,
	)

	controller := controllers.NewExecutionController(executionService)
	router := echo.New()
	router.POST("/execute/:id", controller.Execute)

	return router, experimentsRepo, executionsRepo, signalsRepo, mockSink, dbHandler
}

func TestIntegration_ExecuteExperiment_Success(t *testing.T) {
	router, _, executionsRepo, signalsRepo, mockSink, dbHandler := setupIntegrationApp(t)
	testdb.CleanSignalsTable(t, dbHandler)

	mockSink.EXPECT().
		Send(mock.MatchedBy(func(signal *models.Signal) bool {
			return signal != nil &&
				(signal.Asset == "AAPL" || signal.Asset == "MSFT" || signal.Asset == "NVDA") &&
				(signal.Type == models.SignalBuy || signal.Type == models.SignalHold || signal.Type == models.SignalSell) &&
				signal.Price > 0.0 &&
				signal.StrategyName != "" &&
				!signal.Timestamp.IsZero()
		})).
		Return(nil).
		Times(6)

	// 1. Perform HTTP request to execute Experiment #1
	// Experiment 1 has 6 Variants (AAPL, MSFT, NVDA across 50-day SMA and Price Momentum) in MySQL
	req := httptest.NewRequest(http.MethodPost, "/execute/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// 2. Assert HTTP Status is 200 OK
	require.Equal(t, http.StatusOK, rec.Code)

	// 3. Verify execution and execution_variants persistence in MySQL
	executions, err := executionsRepo.FindByExperiment(1, 10)
	require.NoError(t, err)
	require.Len(t, executions, 1, "Expected 1 execution record for Experiment 1")
	assert.Equal(t, string(models.ExecutionCompleted), executions[0].Status)

	execDetail, err := executionsRepo.FindByID(executions[0].ID)
	require.NoError(t, err)
	require.Len(t, execDetail.Variants, 6, "Expected 6 execution_variants records (1 per variant)")
	for _, ev := range execDetail.Variants {
		assert.Equal(t, string(models.ExecutionVariantCompleted), ev.Status)
		assert.Equal(t, 1, ev.Attempts, "Expected 1 attempt for successful execution variant")
		assert.NotNil(t, ev.Params.Get(), "Expected snapshot params persisted in execution_variants")
	}

	// 4. Verify historical persistence in real MySQL database via signalsRepo
	noDate := models.Nullable[time.Time]{}
	persistedSignalsAAPL, err := signalsRepo.FindByAsset("AAPL", noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignalsAAPL, 2, "Expected 2 records persisted in signals table for AAPL (SMA 50 + Price Momentum)")

	persistedSignalsMSFT, err := signalsRepo.FindByAsset("MSFT", noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignalsMSFT, 2, "Expected 2 records persisted in signals table for MSFT (SMA 50 + Price Momentum)")

	persistedSignalsNVDA, err := signalsRepo.FindByAsset("NVDA", noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignalsNVDA, 2, "Expected 2 records persisted in signals table for NVDA (SMA 50 + Price Momentum)")
}

func TestIntegration_ExecuteExperiment_QuantitoRouteNotFound_FailedVariantsAndAttempts(t *testing.T) {
	router, _, executionsRepo, signalsRepo, mockSink, dbHandler := setupIntegrationApp(t)
	testdb.CleanSignalsTable(t, dbHandler)

	// 1. Perform HTTP request to execute Experiment #3 (Clockwork_BTC_2H_WalkForward)
	// Experiment 3 has 3 ACTIVE variants using strategy 'clockwork_hgb', which is not registered in quantito (returns 404)
	req := httptest.NewRequest(http.MethodPost, "/execute/3", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// 2. Assert HTTP Status is 200 OK and no signals were dispatched to sinks
	require.Equal(t, http.StatusOK, rec.Code)
	mockSink.AssertNotCalled(t, "Send", mock.Anything)

	// 3. Verify execution record was created and marked as FAILED
	executions, err := executionsRepo.FindByExperiment(3, 10)
	require.NoError(t, err)
	require.Len(t, executions, 1, "Expected 1 execution record created for Experiment 3")
	assert.Equal(t, string(models.ExecutionFailed), executions[0].Status)
	assert.True(t, executions[0].ErrorMessage.Valid, "Expected error message on failed execution")
	assert.Equal(t, "All variants failed during execution", executions[0].ErrorMessage.Val)

	// 4. Verify all 3 execution_variants were created, marked as FAILED, and recorded all attempts performed
	execDetail, err := executionsRepo.FindByID(executions[0].ID)
	require.NoError(t, err)
	require.Len(t, execDetail.Variants, 3, "Expected 3 execution_variants records for Experiment 3")
	for _, ev := range execDetail.Variants {
		assert.Equal(t, string(models.ExecutionVariantFailed), ev.Status)
		assert.GreaterOrEqual(t, ev.Attempts, 1, "Expected at least 1 attempt performed")
		assert.Equal(t, services.DefaultMaxAttempts, ev.Attempts, "Expected execution_variant to record all attempts performed")
		assert.True(t, ev.ErrorMessage.Valid, "Expected error message persisted in failed execution_variant")
		assert.Contains(t, ev.ErrorMessage.Val, "quantito api error evaluating clockwork_hgb")
		assert.NotNil(t, ev.Params.Get(), "Expected snapshot params persisted in execution_variants")
	}

	// 5. Verify no signals were created for this execution or asset
	execSignals, err := signalsRepo.FindByExecution(executions[0].ID)
	require.NoError(t, err)
	assert.Empty(t, execSignals, "No signals should be created when variant evaluation fails")

	noDate := models.Nullable[time.Time]{}
	persistedSignalsBTC, err := signalsRepo.FindByAsset("BTCUSDT", noDate, noDate)
	require.NoError(t, err)
	assert.Empty(t, persistedSignalsBTC, "No signals should be persisted for BTCUSDT when variants fail")
}

func TestIntegration_ExecuteExperiment_Inactive(t *testing.T) {
	router, _, executionsRepo, signalsRepo, mockSink, dbHandler := setupIntegrationApp(t)
	testdb.CleanSignalsTable(t, dbHandler)

	// Execute Experiment #2 which is PAUSED in MySQL seed data
	req := httptest.NewRequest(http.MethodPost, "/execute/2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should respond 200 OK and gracefully skip execution
	assert.Equal(t, http.StatusOK, rec.Code)
	mockSink.AssertNotCalled(t, "Send", mock.Anything)

	executions, err := executionsRepo.FindByExperiment(2, 10)
	require.NoError(t, err)
	assert.Empty(t, executions, "No execution or execution_variants should be created for inactive experiment")

	noDate := models.Nullable[time.Time]{}
	persistedSignals, err := signalsRepo.FindByAsset("BTC-USD", noDate, noDate)
	require.NoError(t, err)
	assert.Empty(t, persistedSignals, "No records should be persisted for inactive experiment")
}

func TestIntegration_ExecuteExperiment_NotFound(t *testing.T) {
	router, _, executionsRepo, _, mockSink, dbHandler := setupIntegrationApp(t)
	testdb.CleanSignalsTable(t, dbHandler)

	// Execute non-existent experiment
	req := httptest.NewRequest(http.MethodPost, "/execute/999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "error finding experiment 999")
	mockSink.AssertNotCalled(t, "Send", mock.Anything)

	executions, err := executionsRepo.FindByExperiment(999, 10)
	require.NoError(t, err)
	assert.Empty(t, executions, "No execution should be created for non-existent experiment")
}

func TestIntegration_ExecuteExperiment_InvalidIDParam(t *testing.T) {
	router, _, _, _, _, _ := setupIntegrationApp(t)

	// Execute with non-integer parameter
	req := httptest.NewRequest(http.MethodPost, "/execute/invalid-param", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid ID parameter")
}
