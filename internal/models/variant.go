package models

type Variant struct {
	ID               int                        `db:"id" json:"id"`
	ExperimentID     int                        `db:"experiment_id" json:"experimentId"`
	StrategyID       int                        `db:"strategy_id" json:"strategyId"`
	AssetID          int                        `db:"asset_id" json:"assetId"`
	Asset            string                     `db:"asset" json:"asset,omitempty"`
	Name             string                     `db:"name" json:"name"`
	Params           JSONColumn[map[string]any] `db:"params" json:"params"`
	MaxCapitalPct    float64                    `db:"max_capital_pct" json:"maxCapitalPct"`
	MaxAttempts      int                        `db:"max_attempts" json:"maxAttempts"`
	Status           string                     `db:"status" json:"status"`
	Notes            Nullable[string]           `db:"notes" json:"notes"`
	StrategyName     string                     `db:"strategy_name" json:"strategyName,omitempty"`
	Class            string                     `db:"class" json:"class,omitempty"`
	ParamDefinitions JSONColumn[map[string]any] `db:"param_definitions" json:"paramDefinitions,omitempty"`
}
