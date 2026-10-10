package models

type Strategy struct {
	ID               int                        `db:"id" json:"id"`
	Name             string                     `db:"name" json:"name"`
	Class            string                     `db:"class" json:"class"`
	ParamDefinitions JSONColumn[map[string]any] `db:"param_definitions" json:"paramDefinitions"`
	Status           string                     `db:"status" json:"status"`
}
