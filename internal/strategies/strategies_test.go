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

	_, err = factory.CreateStrategy("unknown")
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
			resp.Signals = []strategies.QuantitoSignalResponse{
				{
					Symbol:       req.Symbol,
					SignalType:   "BUY",
					Price:        req.MarketData[0]["close"].(float64),
					StrategyName: "fifty_day_avg",
					Confidence:   0.8,
				},
			}
		} else if diffPct < -threshPct {
			resp.Signals = []strategies.QuantitoSignalResponse{
				{
					Symbol:       req.Symbol,
					SignalType:   "SELL",
					Price:        req.MarketData[0]["close"].(float64),
					StrategyName: "fifty_day_avg",
					Confidence:   0.8,
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := infraestructures.NewAPIClient(server.URL, nil)
	strat := strategies.NewFiftyDayAvgStrategy(client)

	// Price > Avg -> Buy
	signal, err := strat.GenerateSignal(101, map[string]any{
		"current_price": 110.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalBuy, signal.Type)
	assert.Equal(t, 110.0, signal.Price)

	// Price < Avg -> Sell
	signal, err = strat.GenerateSignal(101, map[string]any{
		"current_price": 90.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalSell, signal.Type)

	// Price close to Avg -> Hold
	signal, err = strat.GenerateSignal(101, map[string]any{
		"current_price": 101.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalHold, signal.Type)
}
