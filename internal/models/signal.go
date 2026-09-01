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
	}
}

// SignalRecord represents the persisted signal entry in the database.
type SignalRecord struct {
	ID           int                 `db:"id" json:"id"`
	VariantID    int                 `db:"variant_id" json:"variantId"`
	ExperimentID int                 `db:"experiment_id" json:"experimentId"`
	StrategyID   int                 `db:"strategy_id" json:"strategyId"`
	AssetID      int                 `db:"asset_id" json:"assetId"`
	Signal       SignalType          `db:"signal" json:"signal"`
	Price        float64             `db:"price" json:"price"`
	Confidence   float64             `db:"confidence" json:"confidence"`
	CreatedAt    time.Time           `db:"created_at" json:"createdAt"`
	Completed    Nullable[time.Time] `db:"completed" json:"completed,omitempty"`
}
