package strategies

import (
	"fmt"
	"os"
	"time"

	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type QuantitoEvaluateRequest struct {
	Symbol     string           `json:"symbol,omitempty"`
	Params     map[string]any   `json:"params"`
	MarketData []map[string]any `json:"market_data"`
}

type QuantitoSignalResponse struct {
	Symbol       string         `json:"symbol"`
	SignalType   string         `json:"signal_type"`
	Price        float64        `json:"price"`
	Timestamp    string         `json:"timestamp"`
	StrategyName string         `json:"strategy_name"`
	Confidence   float64        `json:"confidence"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

type QuantitoEvaluateResponse struct {
	Strategy string                   `json:"strategy"`
	Signals  []QuantitoSignalResponse `json:"signals"`
}

type FiftyDayAvgStrategy struct {
	apiClient interfaces.IAPIClient
}

func NewFiftyDayAvgStrategy(apiClient ...interfaces.IAPIClient) interfaces.IStrategy {
	if len(apiClient) > 0 && apiClient[0] != nil {
		return &FiftyDayAvgStrategy{apiClient: apiClient[0]}
	}

	baseUrl := os.Getenv("QUANTITO_API_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:8081"
	}

	return &FiftyDayAvgStrategy{
		apiClient: infraestructures.NewAPIClient(baseUrl, nil),
	}
}

func getFloatParam(params map[string]any, keys ...string) (float64, bool) {
	for _, key := range keys {
		if val, exists := params[key]; exists {
			switch v := val.(type) {
			case float64:
				return v, true
			case float32:
				return float64(v), true
			case int:
				return float64(v), true
			case int64:
				return float64(v), true
			}
		}
	}
	return 0, false
}

func (s *FiftyDayAvgStrategy) GenerateSignal(asset string, params map[string]any) (*models.Signal, error) {
	currentPrice, ok := getFloatParam(params, "current_price", "close", "lastTradePrice", "price")
	if !ok {
		currentPrice = 100.0
	}

	fiftyAvg, ok := getFloatParam(params, "fifty_day_avg", "fiftyDayAveragePrice", "avg_50")
	if !ok {
		fiftyAvg = 100.0
	}

	var thresholdPct float64
	if tPct, ok := getFloatParam(params, "threshold_pct"); ok {
		thresholdPct = tPct
	} else if t, ok := getFloatParam(params, "threshold"); ok {
		if t <= 1.0 && t > 0.0 {
			thresholdPct = t * 100.0
		} else {
			thresholdPct = t
		}
	}

	signalOn := "both"
	if sOn, ok := params["signal_on"].(string); ok && sOn != "" {
		signalOn = sOn
	}

	symbol := asset
	if sym, ok := params["symbol"].(string); ok && sym != "" {
		symbol = sym
	}

	req := QuantitoEvaluateRequest{
		Symbol: symbol,
		Params: map[string]any{
			"threshold_pct": thresholdPct,
			"signal_on":     signalOn,
		},
		MarketData: []map[string]any{
			{
				"symbol":               symbol,
				"close":                currentPrice,
				"fiftyDayAveragePrice": fiftyAvg,
			},
		},
	}

	var evalResp QuantitoEvaluateResponse
	err := s.apiClient.Post(
		"/api/v1/strategies/fifty_day_avg/evaluate",
		&evalResp,
		infraestructures.RequestJsonBody(req),
	)
	if err != nil {
		return nil, fmt.Errorf("quantito api error evaluating fifty_day_avg: %w", err)
	}

	if len(evalResp.Signals) == 0 {
		return models.NewSignal(asset, models.SignalHold, currentPrice, "fifty_day_avg"), nil
	}

	first := evalResp.Signals[0]
	signalType := models.SignalType(first.SignalType)
	price := first.Price
	if price == 0 {
		price = currentPrice
	}

	signal := models.NewSignal(asset, signalType, price, "fifty_day_avg")
	if first.Confidence > 0 {
		signal.Confidence = first.Confidence
	}
	if len(first.Metadata) > 0 {
		signal.Metadata = first.Metadata
	}
	if first.Timestamp != "" {
		if parsedTime, parseErr := time.Parse(time.RFC3339, first.Timestamp); parseErr == nil {
			signal.Timestamp = parsedTime
		}
	}

	return signal, nil
}
