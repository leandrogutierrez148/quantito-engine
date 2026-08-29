package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"quantito-engine/internal/controllers"
	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/models"
	"quantito-engine/internal/services"
	"quantito-engine/internal/sinks"
	"quantito-engine/internal/strategies"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// In-Memory Repository Implementations for Integration Testing
type inMemoryExperimentsRepo struct {
	experiments map[int]*models.Experiment
}

func (m *inMemoryExperimentsRepo) Find(includeInactive bool) ([]*models.Experiment, error) {
	var list []*models.Experiment
	for _, exp := range m.experiments {
		if includeInactive || exp.IsActive() {
			list = append(list, exp)
		}
	}
	return list, nil
}

func (m *inMemoryExperimentsRepo) FindByID(id int) (*models.Experiment, error) {
	exp, exists := m.experiments[id]
	if !exists {
		return nil, errors.New("experiment not found")
	}
	return exp, nil
}

func (m *inMemoryExperimentsRepo) FindByExperiment(name string) (*models.Experiment, error) {
	for _, exp := range m.experiments {
		if exp.Experiment == name {
			return exp, nil
		}
	}
	return nil, errors.New("experiment not found")
}

type inMemoryAssetVariantsRepo struct {
	mu      sync.Mutex
	records []*models.AssetVariant
}

func (m *inMemoryAssetVariantsRepo) Save(assetID, variantID int, signal *models.Signal) (*models.AssetVariant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rec := &models.AssetVariant{
		ID:         len(m.records) + 1,
		AssetID:    assetID,
		VariantID:  variantID,
		Signal:     signal.Type,
		Price:      signal.Price,
		Confidence: signal.Confidence,
		CreatedAt:  signal.Timestamp,
	}
	m.records = append(m.records, rec)
	return rec, nil
}

func (m *inMemoryAssetVariantsRepo) Find(assetID int, fromDate, toDate models.Nullable[time.Time]) ([]*models.AssetVariant, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var res []*models.AssetVariant
	for _, r := range m.records {
		if r.AssetID == assetID {
			res = append(res, r)
		}
	}
	return res, nil
}

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

func setupIntegrationApp(t *testing.T) (*echo.Echo, *inMemoryExperimentsRepo, *inMemoryAssetVariantsRepo, *capturingSignalSink) {
	// Seed Data: Experiment with 2 Assets and 2 Competing Strategy Variants (Contraposition)
	exp1 := &models.Experiment{
		ID:          1,
		Experiment:  "Tech_Momentum_vs_SMA",
		Description: models.NotNull("Test comparing 50-day moving average and price momentum on Tech assets"),
		Params:      models.NewJsonColumn(map[string]any{"benchmark": "SPY", "execution_mode": "PAPER"}),
		Status:      "ACTIVE",
		Assets: []*models.ExperimentAsset{
			{ID: 10, ExperimentID: 1, AssetID: 101, Status: "ACTIVE"}, // AAPL
			{ID: 11, ExperimentID: 1, AssetID: 102, Status: "ACTIVE"}, // MSFT
		},
		Variants: []*models.Variant{
			{
				ID:           1,
				ExperimentID: 1,
				StrategyID:   10,
				Variant:      "Variant A (SMA 50)",
				StrategyName: "SMA 50 Strategy",
				Class:        "fifty_day_avg",
			},
			{
				ID:           2,
				ExperimentID: 1,
				StrategyID:   20,
				Variant:      "Variant B (Momentum)",
				StrategyName: "Momentum Strategy",
				Class:        "price_momentum",
			},
		},
	}

	expInactive := &models.Experiment{
		ID:          2,
		Experiment:  "Paused_Study",
		Status:      "PAUSED",
		Assets: []*models.ExperimentAsset{
			{ID: 20, ExperimentID: 2, AssetID: 201, Status: "ACTIVE"},
		},
		Variants: []*models.Variant{
			{
				ID:           3,
				ExperimentID: 2,
				StrategyID:   10,
				Variant:      "Variant A",
				Class:        "fifty_day_avg",
			},
		},
	}

	expRepo := &inMemoryExperimentsRepo{
		experiments: map[int]*models.Experiment{
			1: exp1,
			2: expInactive,
		},
	}
	assetVariantsRepo := &inMemoryAssetVariantsRepo{}

	// Real Factories, Sinks, and Mutex
	strategyFactory := strategies.NewStrategyFactory()
	capturingSink := newCapturingSignalSink()
	consoleSink := sinks.NewConsoleSink()
	signalDispatcher := sinks.NewSignalDispatcher(consoleSink, capturingSink)
	executionMutex := infraestructures.NewExecutionMutex()

	// Real ExecutionService
	executionService := services.NewExecutionService(
		expRepo,
		assetVariantsRepo,
		strategyFactory,
		signalDispatcher,
		executionMutex,
	)

	// Real Controller and Echo Router
	controller := controllers.NewExecutionController(executionService)
	router := echo.New()
	router.POST("/execute/:id", controller.Execute)

	return router, expRepo, assetVariantsRepo, capturingSink
}

func TestIntegration_ExecuteExperiment_Success(t *testing.T) {
	router, _, assetVariantsRepo, capturingSink := setupIntegrationApp(t)

	// 1. Perform HTTP request to execute Experiment #1
	req := httptest.NewRequest(http.MethodPost, "/execute/1", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// 2. Assert HTTP Status is 200 OK
	require.Equal(t, http.StatusOK, rec.Code)

	// 3. Verify that all 4 evaluations occurred (2 assets * 2 strategies in contraposition)
	signals := capturingSink.GetSignals()
	require.Len(t, signals, 4, "Expected 4 signals generated (2 assets x 2 strategies)")

	// Group signals by AssetID
	signalsByAsset := make(map[int][]*models.Signal)
	for _, s := range signals {
		signalsByAsset[s.AssetID] = append(signalsByAsset[s.AssetID], s)
		assert.Contains(t, []models.SignalType{models.SignalBuy, models.SignalHold, models.SignalSell}, s.Type)
		assert.Greater(t, s.Price, 0.0)
		assert.NotEmpty(t, s.StrategyName)
		assert.False(t, s.Timestamp.IsZero())
	}

	// Verify both assets (101 and 102) were evaluated
	assert.Len(t, signalsByAsset[101], 2, "Asset 101 should have 2 signals (one per strategy)")
	assert.Len(t, signalsByAsset[102], 2, "Asset 102 should have 2 signals (one per strategy)")

	// 4. Verify historical persistence in assetVariantsRepo
	assert.Len(t, assetVariantsRepo.records, 4, "Expected 4 records persisted in asset_variants")
}

func TestIntegration_ExecuteExperiment_Inactive(t *testing.T) {
	router, _, assetVariantsRepo, capturingSink := setupIntegrationApp(t)

	// Execute Experiment #2 which is PAUSED
	req := httptest.NewRequest(http.MethodPost, "/execute/2", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	// Should respond 200 OK and gracefully skip execution
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, capturingSink.GetSignals(), "No signals should be emitted for inactive experiment")
	assert.Empty(t, assetVariantsRepo.records, "No records should be persisted for inactive experiment")
}

func TestIntegration_ExecuteExperiment_NotFound(t *testing.T) {
	router, _, _, _ := setupIntegrationApp(t)

	// Execute non-existent experiment
	req := httptest.NewRequest(http.MethodPost, "/execute/999", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.Contains(t, rec.Body.String(), "error finding experiment 999")
}

func TestIntegration_ExecuteExperiment_InvalidIDParam(t *testing.T) {
	router, _, _, _ := setupIntegrationApp(t)

	// Execute with non-integer parameter
	req := httptest.NewRequest(http.MethodPost, "/execute/invalid-param", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid ID parameter")
}
