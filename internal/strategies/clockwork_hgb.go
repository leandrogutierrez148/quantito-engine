package strategies

import (
	"fmt"
	"net/url"
	"time"

	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
	"quantito-engine/internal/utils"
)

type ClockworkParams struct {
	APIURL        string           `json:"api_url"`
	CurrentPrice  float64          `json:"current_price" validate:"required,gt=0"`
	GateBps       float64          `json:"gate_bps" validate:"required,gt=0"`
	ModelVariant  string           `json:"model_variant" validate:"required"`
	PredictionBps float64          `json:"prediction_bps"`
	MarketData    []map[string]any `json:"market_data"`
}

type ClockworkHGBStrategy struct {
	apiClient interfaces.IAPIClient
}

func NewClockworkHGBStrategy(apiClient interfaces.IAPIClient) interfaces.IStrategy {
	return &ClockworkHGBStrategy{
		apiClient: apiClient,
	}
}

func (s *ClockworkHGBStrategy) GenerateSignal(asset string, params map[string]any) (*models.Signal, error) {
	mapped, err := utils.MapAndValidateParams[ClockworkParams](params)
	if err != nil {
		return nil, fmt.Errorf("invalid clockwork_hgb params: %w", err)
	}

	marketData := mapped.MarketData
	if len(marketData) == 0 {
		marketData = []map[string]any{
			{
				"symbol":         asset,
				"close":          mapped.CurrentPrice,
				"prediction_bps": mapped.PredictionBps,
			},
		}
	}

	req := QuantitoEvaluateRequest{
		Symbol: asset,
		Params: map[string]any{
			"model_variant": mapped.ModelVariant,
			"gate_bps":      mapped.GateBps,
		},
		MarketData: marketData,
	}

	endpoint := "/api/v1/strategies/clockwork_hgb/evaluate"
	if mapped.APIURL != "" {
		endpoint, err = url.JoinPath(mapped.APIURL, endpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid clockwork_hgb api_url: %w", err)
		}
	}

	var evalResp QuantitoEvaluateResponse
	err = s.apiClient.Post(
		endpoint,
		&evalResp,
		infraestructures.RequestJsonBody(req),
	)
	if err != nil {
		return nil, fmt.Errorf("quantito api error evaluating clockwork_hgb: %w", err)
	}

	if evalResp.Signal == nil {
		return models.NewSignal(asset, models.SignalHold, mapped.CurrentPrice, "clockwork_hgb"), nil
	}

	first := evalResp.Signal
	signalType := models.SignalType(first.SignalType)
	price := first.Price
	if price == 0 {
		price = mapped.CurrentPrice
	}

	signal := models.NewSignal(asset, signalType, price, "clockwork_hgb")
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
