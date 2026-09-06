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

func setupIntegrationApp(t *testing.T) (*echo.Echo, interfaces.IExperimentsRepository, interfaces.ISignalsRepository, *capturingSignalSink, interfaces.IDBHandler) {
	dbHandler := testdb.NewTestDBHandler(t)
	testdb.CleanSignalsTable(t, dbHandler)

	experimentsRepo := repositories.NewExperimentsRepository(dbHandler)
	signalsRepo := repositories.NewSignalsRepository(dbHandler)

	strategyFactory := strategies.NewStrategyFactory()
	capturingSink := newCapturingSignalSink()
	consoleSink := sinks.NewConsoleSink()
	signalDispatcher := sinks.NewSignalDispatcher(consoleSink, capturingSink)
	executionMutex := infraestructures.NewExecutionMutex()

	executionService := services.NewExecutionService(
		experimentsRepo,
		signalsRepo,
		strategyFactory,
		signalDispatcher,
		executionMutex,
	)

	controller := controllers.NewExecutionController(executionService)
	router := echo.New()
	router.POST("/execute/:id", controller.Execute)

	return router, experimentsRepo, signalsRepo, capturingSink, dbHandler
}

func TestIntegration_ExecuteExperiment_Success(t *testing.T) {
	router, _, signalsRepo, capturingSink, dbHandler := setupIntegrationApp(t)
	defer testdb.CleanSignalsTable(t, dbHandler)

	// 1. Perform HTTP request to execute Experiment #1
	// Experiment 1 has 3 Assets (101, 102, 103) and 2 Variants (50-day SMA, Price Momentum) in MySQL
	req := httptest.NewRequest(http.MethodPost, "/execute/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// 2. Assert HTTP Status is 200 OK
	require.Equal(t, http.StatusOK, rec.Code)

	// 3. Verify that 6 evaluations occurred (3 assets * 2 strategies in contraposition)
	signals := capturingSink.GetSignals()
	require.Len(t, signals, 6, "Expected 6 signals generated (3 assets x 2 strategies)")

	// Group signals by AssetID
	signalsByAsset := make(map[int][]*models.Signal)
	for _, s := range signals {
		signalsByAsset[s.AssetID] = append(signalsByAsset[s.AssetID], s)
		assert.Contains(t, []models.SignalType{models.SignalBuy, models.SignalHold, models.SignalSell}, s.Type)
		assert.Greater(t, s.Price, 0.0)
		assert.NotEmpty(t, s.StrategyName)
		assert.False(t, s.Timestamp.IsZero())
	}

	// Verify all 3 assets (101, 102, 103) were evaluated with 2 signals each
	assert.Len(t, signalsByAsset[101], 2, "Asset 101 should have 2 signals (one per strategy)")
	assert.Len(t, signalsByAsset[102], 2, "Asset 102 should have 2 signals (one per strategy)")
	assert.Len(t, signalsByAsset[103], 2, "Asset 103 should have 2 signals (one per strategy)")

	// 4. Verify historical persistence in real MySQL database via signalsRepo
	noDate := models.Nullable[time.Time]{}
	persistedSignals101, err := signalsRepo.Find(101, noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignals101, 2, "Expected 2 records persisted in signals table for Asset 101")

	persistedSignals102, err := signalsRepo.Find(102, noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignals102, 2, "Expected 2 records persisted in signals table for Asset 102")

	persistedSignals103, err := signalsRepo.Find(103, noDate, noDate)
	require.NoError(t, err)
	assert.Len(t, persistedSignals103, 2, "Expected 2 records persisted in signals table for Asset 103")
}

func TestIntegration_ExecuteExperiment_Inactive(t *testing.T) {
	router, _, signalsRepo, capturingSink, dbHandler := setupIntegrationApp(t)
	defer testdb.CleanSignalsTable(t, dbHandler)

	// Execute Experiment #2 which is PAUSED in MySQL seed data
	req := httptest.NewRequest(http.MethodPost, "/execute/2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should respond 200 OK and gracefully skip execution
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, capturingSink.GetSignals(), "No signals should be emitted for inactive experiment")

	noDate := models.Nullable[time.Time]{}
	persistedSignals, err := signalsRepo.Find(201, noDate, noDate)
	require.NoError(t, err)
	assert.Empty(t, persistedSignals, "No records should be persisted for inactive experiment")
}

func TestIntegration_ExecuteExperiment_NotFound(t *testing.T) {
	router, _, _, _, dbHandler := setupIntegrationApp(t)
	defer testdb.CleanSignalsTable(t, dbHandler)

	// Execute non-existent experiment
	req := httptest.NewRequest(http.MethodPost, "/execute/999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "error finding experiment 999")
}

func TestIntegration_ExecuteExperiment_InvalidIDParam(t *testing.T) {
	router, _, _, _, dbHandler := setupIntegrationApp(t)
	defer testdb.CleanSignalsTable(t, dbHandler)

	// Execute with non-integer parameter
	req := httptest.NewRequest(http.MethodPost, "/execute/invalid-param", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid ID parameter")
}
