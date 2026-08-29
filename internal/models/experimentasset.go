package models

import "time"

type ExperimentAsset struct {
	ID           int       `db:"id" json:"id"`
	ExperimentID int       `db:"experiment_id" json:"experimentId"`
	AssetID      int       `db:"asset_id" json:"assetId"`
	Status       string    `db:"status" json:"status"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
}
