package models

type Variant struct {
	ID                   int                        `db:"id" json:"id"`
	ExperimentID         int                        `db:"experiment_id" json:"experimentId"`
	StrategyID           int                        `db:"strategy_id" json:"strategyId"`
	Name                 string                     `db:"name" json:"name"`
	Params               JSONColumn[map[string]any] `db:"params" json:"params"`
	CapitalAllocationPct float64                    `db:"capital_allocation_pct" json:"capitalAllocationPct"`
	Status               string                     `db:"status" json:"status"`
	Notes                Nullable[string]           `db:"notes" json:"notes"`
	StrategyName         string                     `db:"strategy_name" json:"strategyName,omitempty"`
	Class                string                     `db:"class" json:"class,omitempty"`
	ParamDefinitions     JSONColumn[map[string]any] `db:"param_definitions" json:"paramDefinitions,omitempty"`
}
