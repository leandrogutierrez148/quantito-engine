package models

import "time"

type Experiment struct {
	ID            int                        `db:"id" json:"id"`
	Name          string                     `db:"name" json:"name"`
	Description   Nullable[string]           `db:"description" json:"description"`
	Status        string                     `db:"status" json:"status"`
	ExecutionMode string                     `db:"execution_mode" json:"executionMode"`
	StartDate     Nullable[time.Time]        `db:"start_date" json:"startDate,omitempty"`
	EndDate       Nullable[time.Time]        `db:"end_date" json:"endDate,omitempty"`
	Frequency     string                     `db:"frequency" json:"frequency"`
	MaxCapital    float64                    `db:"max_capital" json:"maxCapital"`
	Currency      string                     `db:"currency" json:"currency"`
	Assets        JSONColumn[[]string]       `db:"assets" json:"assets"`
	Params        JSONColumn[map[string]any] `db:"params" json:"params,omitempty"`
	Notes         Nullable[string]           `db:"notes" json:"notes"`
	Variants      []*Variant                 `db:"-" json:"variants,omitempty"`
}

func (e *Experiment) IsActive() bool {
	return e != nil && (e.Status == "ACTIVE" || e.Status == "A")
}

func (e *Experiment) GetParam(key string) (any, bool) {
	if e == nil || e.Params.Get() == nil {
		return nil, false
	}
	val, ok := (*e.Params.Get())[key]
	return val, ok
}

func (e *Experiment) GetAssets() []string {
	if e == nil || e.Assets.Get() == nil {
		return nil
	}
	return *e.Assets.Get()
}
