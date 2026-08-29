package models

import "time"

type AssetVariant struct {
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
