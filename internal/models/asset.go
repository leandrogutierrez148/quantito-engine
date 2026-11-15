package models

import "time"

type Asset struct {
	ID        int       `db:"id" json:"id"`
	Symbol    string    `db:"symbol" json:"symbol"`
	Name      string    `db:"name" json:"name"`
	AssetType string    `db:"asset_type" json:"assetType"`
	Exchange  Nullable[string] `db:"exchange" json:"exchange,omitempty"`
	Currency  string    `db:"currency" json:"currency"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
