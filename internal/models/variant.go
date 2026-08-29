package models

type Variant struct {
	ID           int                        `db:"id" json:"id"`
	ExperimentID int                        `db:"experiment_id" json:"experimentId"`
	StrategyID   int                        `db:"strategy_id" json:"strategyId"`
	Variant      string                     `db:"variant" json:"variant"`
	Notes        Nullable[string]           `db:"notes" json:"notes"`
	StrategyName string                     `db:"strategy_name" json:"strategyName,omitempty"`
	Class        string                     `db:"class" json:"class,omitempty"`
	Params       JSONColumn[map[string]any] `db:"params" json:"params,omitempty"`
}
