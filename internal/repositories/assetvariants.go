package repositories

import (
	"time"

	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type assetVariantsRepository struct {
	db interfaces.IDBHandler
}

func NewAssetVariantsRepository(db interfaces.IDBHandler) interfaces.IAssetVariantsRepository {
	return &assetVariantsRepository{db}
}

func (r *assetVariantsRepository) Save(assetID, variantID int, signal *models.Signal) (*models.AssetVariant, error) {
	var result models.AssetVariant
	err := r.db.QueryOne(
		"CALL sp_save_asset_variant (?, ?, ?, ?, ?)",
		&result,
		assetID,
		variantID,
		string(signal.Type),
		signal.Price,
		signal.Confidence,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *assetVariantsRepository) Find(assetID int, fromDate, toDate models.Nullable[time.Time]) ([]*models.AssetVariant, error) {
	var results []*models.AssetVariant

	var fromVal any
	if fromDate.Valid {
		fromVal = fromDate.Val
	}
	var toVal any
	if toDate.Valid {
		toVal = toDate.Val
	}

	err := r.db.QueryAll("CALL sp_find_asset_variants (?, ?, ?)", &results, assetID, fromVal, toVal)
	if err != nil {
		return nil, err
	}
	return results, nil
}
