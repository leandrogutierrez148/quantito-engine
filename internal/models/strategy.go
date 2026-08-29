package models

type Strategy struct {
	ID     int                        `db:"id" json:"id"`
	Name   string                     `db:"name" json:"name"`
	Class  string                     `db:"class" json:"class"`
	Params JSONColumn[map[string]any] `db:"params" json:"params"`
	Status string                     `db:"status" json:"status"`
}
