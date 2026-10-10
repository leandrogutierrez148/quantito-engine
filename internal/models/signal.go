package models

import "time"

type SignalType string

const (
	SignalBuy  SignalType = "BUY"
	SignalHold SignalType = "HOLD"
	SignalSell SignalType = "SELL"
)

// Signal represents an in-memory trading recommendation emitted by a strategy.
type Signal struct {
	Asset        string         `json:"asset"`
	Type         SignalType     `json:"type"`
	Price        float64        `json:"price"`
	Timestamp    time.Time      `json:"timestamp"`
	StrategyName string         `json:"strategyName"`
	Confidence   float64        `json:"confidence"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

func NewSignal(asset string, signalType SignalType, price float64, strategyName string) *Signal {
	return &Signal{
		Asset:        asset,
		Type:         signalType,
		Price:        price,
		Timestamp:    time.Now().UTC(),
		StrategyName: strategyName,
		Confidence:   1.0,
	}
}

// SignalRecord represents the persisted signal entry in the database.
type SignalRecord struct {
	ID                 int        `db:"id" json:"id"`
	ExecutionVariantID int        `db:"execution_variant_id" json:"executionVariantId"`
	ExecutionID        int        `db:"execution_id" json:"executionId,omitempty"`
	VariantID          int        `db:"variant_id" json:"variantId,omitempty"`
	VariantName        string     `db:"variant_name" json:"variantName,omitempty"`
	Asset              string     `db:"asset" json:"asset"`
	Signal             SignalType `db:"signal" json:"signal"`
	Price              float64    `db:"price" json:"price"`
	Confidence         float64    `db:"confidence" json:"confidence"`
	CreatedAt          time.Time  `db:"created_at" json:"createdAt"`
}
