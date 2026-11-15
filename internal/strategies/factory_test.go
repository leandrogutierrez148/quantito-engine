package strategies_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/models"
	"quantito-engine/internal/strategies"

	"github.com/stretchr/testify/assert"
)

func TestStrategyFactory(t *testing.T) {
	factory := strategies.NewStrategyFactory()

	strat, err := factory.CreateStrategy("fifty_day_avg")
	assert.NoError(t, err)
	assert.NotNil(t, strat)

	strat, err = factory.CreateStrategy("price_momentum")
	assert.NoError(t, err)
	assert.NotNil(t, strat)

	strat, err = factory.CreateStrategy("clockwork_hgb")
	assert.NoError(t, err)
	assert.NotNil(t, strat)

	_, err = factory.CreateStrategy("unknown")
	assert.Error(t, err)
}

func TestClockworkHGBStrategy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/strategies/clockwork_hgb/evaluate", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req strategies.QuantitoEvaluateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)

		predBps := 0.0
		price := 0.0
		if len(req.MarketData) > 0 {
			if p, ok := req.MarketData[0]["prediction_bps"].(float64); ok {
				predBps = p
			}
			if closePrice, ok := req.MarketData[0]["close"].(float64); ok {
				price = closePrice
			}
		}
		gateBps := req.Params["gate_bps"].(float64)
		modelVariant := req.Params["model_variant"].(string)

		resp := strategies.QuantitoEvaluateResponse{
			Strategy: "clockwork_hgb",
		}
		if predBps >= gateBps {
			resp.Signal = &strategies.QuantitoSignalResponse{
				Symbol:       req.Symbol,
				SignalType:   "BUY",
				Price:        price,
				StrategyName: "clockwork_hgb",
				Confidence:   0.85,
				Metadata: map[string]any{
					"model_variant": modelVariant,
				},
			}
		} else if predBps <= -gateBps {
			resp.Signal = &strategies.QuantitoSignalResponse{
				Symbol:       req.Symbol,
				SignalType:   "SELL",
				Price:        price,
				StrategyName: "clockwork_hgb",
				Confidence:   0.85,
				Metadata: map[string]any{
					"model_variant": modelVariant,
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := infraestructures.NewAPIClient("", nil)
	strat := strategies.NewClockworkHGBStrategy(client)

	// |prediction_bps| >= gate_bps and positive -> BUY
	signal, err := strat.GenerateSignal("BTCUSDT", map[string]any{
		"api_url":        server.URL,
		"current_price":  63500.0,
		"prediction_bps": 26.5,
		"gate_bps":       19.8,
		"model_variant":  "aug_fixed_postfit",
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalBuy, signal.Type)
	assert.Equal(t, 63500.0, signal.Price)
	assert.Equal(t, "BTCUSDT", signal.Asset)
	assert.Equal(t, "clockwork_hgb", signal.StrategyName)
	assert.Equal(t, 0.85, signal.Confidence)
	assert.Equal(t, "aug_fixed_postfit", signal.Metadata["model_variant"])

	// |prediction_bps| >= gate_bps and negative -> SELL
	signal, err = strat.GenerateSignal("BTCUSDT", map[string]any{
		"api_url":        server.URL,
		"current_price":  63500.0,
		"prediction_bps": -24.0,
		"gate_bps":       20.5,
		"model_variant":  "oct_refit_postfit",
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalSell, signal.Type)
	assert.Equal(t, "oct_refit_postfit", signal.Metadata["model_variant"])

	// |prediction_bps| < gate_bps -> HOLD
	signal, err = strat.GenerateSignal("BTCUSDT", map[string]any{
		"api_url":        server.URL,
		"current_price":  63500.0,
		"prediction_bps": 14.2,
		"gate_bps":       21.0,
		"model_variant":  "jan_refit_train90",
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalHold, signal.Type)
	assert.Equal(t, 63500.0, signal.Price)

	// Missing required parameters -> validation error via utils.MapAndValidateParams
	_, err = strat.GenerateSignal("BTCUSDT", map[string]any{
		"current_price": 63500.0,
	})
	assert.Error(t, err)
}

func TestFiftyDayAvgStrategy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/strategies/fifty_day_avg/evaluate", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var req strategies.QuantitoEvaluateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)

		diffPct := 0.0
		if len(req.MarketData) > 0 {
			closePrice := req.MarketData[0]["close"].(float64)
			avg50 := req.MarketData[0]["fiftyDayAveragePrice"].(float64)
			diffPct = ((closePrice - avg50) / avg50) * 100.0
		}
		threshPct := req.Params["threshold_pct"].(float64)

		resp := strategies.QuantitoEvaluateResponse{
			Strategy: "fifty_day_avg",
		}
		if diffPct > threshPct {
			resp.Signal = &strategies.QuantitoSignalResponse{
				Symbol:       req.Symbol,
				SignalType:   "BUY",
				Price:        req.MarketData[0]["close"].(float64),
				StrategyName: "fifty_day_avg",
				Confidence:   0.8,
			}
		} else if diffPct < -threshPct {
			resp.Signal = &strategies.QuantitoSignalResponse{
				Symbol:       req.Symbol,
				SignalType:   "SELL",
				Price:        req.MarketData[0]["close"].(float64),
				StrategyName: "fifty_day_avg",
				Confidence:   0.8,
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := infraestructures.NewAPIClient(server.URL, nil)
	strat := strategies.NewFiftyDayAvgStrategy(client)

	// Price > Avg -> Buy
	signal, err := strat.GenerateSignal("AAPL", map[string]any{
		"current_price": 110.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalBuy, signal.Type)
	assert.Equal(t, 110.0, signal.Price)
	assert.Equal(t, "AAPL", signal.Asset)

	// Price < Avg -> Sell
	signal, err = strat.GenerateSignal("AAPL", map[string]any{
		"current_price": 90.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalSell, signal.Type)
	assert.Equal(t, "AAPL", signal.Asset)

	// Price close to Avg -> Hold
	signal, err = strat.GenerateSignal("AAPL", map[string]any{
		"current_price": 101.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalHold, signal.Type)
	assert.Equal(t, "AAPL", signal.Asset)
}
