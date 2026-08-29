package models

import (
	"time"
)

type SignalType string

const (
	SignalBuy  SignalType = "BUY"
	SignalHold SignalType = "HOLD"
	SignalSell SignalType = "SELL"
)

type Signal struct {
	AssetID      int            `json:"assetId"`
	Type         SignalType     `json:"type"`
	Price        float64        `json:"price"`
	Timestamp    time.Time      `json:"timestamp"`
	StrategyName string         `json:"strategyName"`
	Confidence   float64        `json:"confidence"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

func NewSignal(assetID int, signalType SignalType, price float64, strategyName string) *Signal {
	return &Signal{
		AssetID:      assetID,
		Type:         signalType,
		Price:        price,
		Timestamp:    time.Now().UTC(),
		StrategyName: strategyName,
		Confidence:   1.0,
		Metadata:     make(map[string]any),
	}
}
