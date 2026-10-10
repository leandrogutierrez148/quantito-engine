package strategies

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type PriceMomentumStrategy struct{}

func NewPriceMomentumStrategy() interfaces.IStrategy {
	return &PriceMomentumStrategy{}
}

func (s *PriceMomentumStrategy) GenerateSignal(asset string, params map[string]any) (*models.Signal, error) {
	currentPrice := 100.0
	if priceParam, ok := params["current_price"].(float64); ok {
		currentPrice = priceParam
	}

	buyThreshold := 0.05
	if bt, ok := params["buy_threshold"].(float64); ok {
		buyThreshold = bt
	}

	sellThreshold := -0.05
	if st, ok := params["sell_threshold"].(float64); ok {
		sellThreshold = st
	}

	signalType := models.SignalHold
	if change, ok := params["pct_change"].(float64); ok {
		if change >= buyThreshold {
			signalType = models.SignalBuy
		} else if change <= sellThreshold {
			signalType = models.SignalSell
		}
	}

	signal := models.NewSignal(asset, signalType, currentPrice, "price_momentum")
	return signal, nil
}
