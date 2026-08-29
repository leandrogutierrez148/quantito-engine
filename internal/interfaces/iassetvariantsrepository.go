package interfaces

import (
	"time"

	"quantito-engine/internal/models"
)

type IAssetVariantsRepository interface {
	Save(assetID, variantID int, signal *models.Signal) (*models.AssetVariant, error)
	Find(assetID int, fromDate, toDate models.Nullable[time.Time]) ([]*models.AssetVariant, error)
}
