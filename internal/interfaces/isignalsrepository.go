package interfaces

import (
	"quantito-engine/internal/models"
	"time"
)

type ISignalsRepository interface {
	Save(assetID, variantID int, signal *models.Signal) (*models.SignalRecord, error)
	Find(assetID int, fromDate, toDate models.Nullable[time.Time]) ([]*models.SignalRecord, error)
}
