package tests

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"quantito-engine/internal/controllers"
	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
	"quantito-engine/internal/repositories"
	"quantito-engine/internal/services"
	"quantito-engine/internal/sinks"
	"quantito-engine/internal/strategies"
	testdb "quantito-engine/tests/integration/db"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Capturing Signal Sink for Integration Assertions
type capturingSignalSink struct {
	mu      sync.Mutex
	signals []*models.Signal
}

func newCapturingSignalSink() *capturingSignalSink {
	return &capturingSignalSink{}
}

func (s *capturingSignalSink) Name() string {
	return "CapturingSignalSink"
}

func (s *capturingSignalSink) Send(signal *models.Signal) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.signals = append(s.signals, signal)
	return nil
}

func (s *capturingSignalSink) GetSignals() []*models.Signal {
	s.mu.Lock()
	defer s.mu.Unlock()
	copied := make([]*models.Signal, len(s.signals))
	copy(copied, s.signals)
	return copied
}

func setupIntegrationApp(t *testing.T) (*echo.Echo, interfaces.IExperimentsRepository, interfaces.IExecutionsRepository, interfaces.ISignalsRepository, *capturingSignalSink, interfaces.IDBHandler) {
	dbHandler := testdb.NewTestDBHandler(t)

	experimentsRepo := repositories.NewExperimentsRepository(dbHandler)
	executionsRepo := repositories.NewExecutionsRepository(dbHandler)
	signalsRepo := repositories.NewSignalsRepository(dbHandler)

	strategyFactory := strategies.NewStrategyFactory()
	capturingSink := newCapturingSignalSink()
	consoleSink := sinks.NewConsoleSink()
	signalDispatcher := sinks.NewSignalDispatcher(consoleSink, capturingSink)
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

	return router, experimentsRepo, executionsRepo, signalsRepo, capturingSink, dbHandler
}

func TestIntegration_ExecuteExperiment_Success(t *testing.T) {
	router, _, _, signalsRepo, capturingSink, dbHandler := setupIntegrationApp(t)
	testdb.CleanSignalsTable(t, dbHandler)

	// 1. Perform HTTP request to execute Experiment #1
	// Experiment 1 has 3 Assets ("AAPL", "MSFT", "NVDA") and 2 Variants (50-day SMA, Price Momentum) in MySQL
	req := httptest.NewRequest(http.MethodPost, "/execute/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// 2. Assert HTTP Status is 200 OK
	require.Equal(t, http.StatusOK, rec.Code)

	// 3. Verify that 6 evaluations occurred (3 assets * 2 strategies in contraposition)
	signals := capturingSink.GetSignals()
	require.Len(t, signals, 6, "Expected 6 signals generated (3 assets x 2 strategies)")

	// Group signals by Asset
	signalsByAsset := make(map[string][]*models.Signal)
	for _, s := range signals {
		signalsByAsset[s.Asset] = append(signalsByAsset[s.Asset], s)
		assert.Contains(t, []models.SignalType{models.SignalBuy, models.SignalHold, models.SignalSell}, s.Type)
		assert.Greater(t, s.Price, 0.0)
		assert.NotEmpty(t, s.StrategyName)
		assert.False(t, s.Timestamp.IsZero())
	}

	// Verify all 3 assets ("AAPL", "MSFT", "NVDA") were evaluated with 2 signals each
	assert.Len(t, signalsByAsset["AAPL"], 2, "AAPL should have 2 signals (one per strategy)")
	assert.Len(t, signalsByAsset["MSFT"], 2, "MSFT should have 2 signals (one per strategy)")
	assert.Len(t, signalsByAsset["NVDA"], 2, "NVDA should have 2 signals (one per strategy)")

	// 4. Verify historical persistence in real MySQL database via signalsRepo
	noDate := models.Nullable[time.Time]{}
	persistedSignalsAAPL, err := signalsRepo.FindByAsset("AAPL", noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignalsAAPL, 2, "Expected 2 records persisted in signals table for AAPL")

	persistedSignalsMSFT, err := signalsRepo.FindByAsset("MSFT", noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignalsMSFT, 2, "Expected 2 records persisted in signals table for MSFT")

	persistedSignalsNVDA, err := signalsRepo.FindByAsset("NVDA", noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignalsNVDA, 2, "Expected 2 records persisted in signals table for NVDA")
}

func TestIntegration_ExecuteExperiment_Inactive(t *testing.T) {
	router, _, _, signalsRepo, capturingSink, _ := setupIntegrationApp(t)

	// Execute Experiment #2 which is PAUSED in MySQL seed data
	req := httptest.NewRequest(http.MethodPost, "/execute/2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should respond 200 OK and gracefully skip execution
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, capturingSink.GetSignals(), "No signals should be emitted for inactive experiment")

	noDate := models.Nullable[time.Time]{}
	persistedSignals, err := signalsRepo.FindByAsset("BTC-USD", noDate, noDate)
	require.NoError(t, err)
	assert.Empty(t, persistedSignals, "No records should be persisted for inactive experiment")
}

func TestIntegration_ExecuteExperiment_NotFound(t *testing.T) {
	router, _, _, _, _, _ := setupIntegrationApp(t)

	// Execute non-existent experiment
	req := httptest.NewRequest(http.MethodPost, "/execute/999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "error finding experiment 999")
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
