package models

type Experiment struct {
	ID          int                        `db:"id" json:"id"`
	Experiment  string                     `db:"experiment" json:"experiment"`
	Description Nullable[string]           `db:"description" json:"description"`
	Params      JSONColumn[map[string]any] `db:"params" json:"params"`
	Status      string                     `db:"status" json:"status"`
	Notes       Nullable[string]           `db:"notes" json:"notes"`
	Assets      []*ExperimentAsset         `db:"-" json:"assets,omitempty"`
	Variants    []*Variant                 `db:"-" json:"variants,omitempty"`
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
