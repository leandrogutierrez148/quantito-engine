package strategies

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type FiftyDayAvgStrategy struct{}

func NewFiftyDayAvgStrategy() interfaces.IStrategy {
	return &FiftyDayAvgStrategy{}
}

func (s *FiftyDayAvgStrategy) GenerateSignal(assetID int, params map[string]any) (*models.Signal, error) {
	// Example calculation based on strategy parameters
	// Default price and signal type
	currentPrice := 100.0
	if priceParam, ok := params["current_price"].(float64); ok {
		currentPrice = priceParam
	}

	threshold := 0.0
	if threshParam, ok := params["threshold"].(float64); ok {
		threshold = threshParam
	}

	signalType := models.SignalHold
	if fiftyAvg, ok := params["fifty_day_avg"].(float64); ok {
		diff := (currentPrice - fiftyAvg) / fiftyAvg
		if diff > threshold {
			signalType = models.SignalBuy
		} else if diff < -threshold {
			signalType = models.SignalSell
		}
	}

	signal := models.NewSignal(assetID, signalType, currentPrice, "fifty_day_avg")
	return signal, nil
}
