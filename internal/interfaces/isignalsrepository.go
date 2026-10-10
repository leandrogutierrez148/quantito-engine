package interfaces

import (
	"quantito-engine/internal/models"
	"time"
)

type ISignalsRepository interface {
	Save(executionVariantID int, asset string, signal *models.Signal) (*models.SignalRecord, error)
	FindByExecution(executionID int) ([]*models.SignalRecord, error)
	FindByAsset(asset string, fromDate, toDate models.Nullable[time.Time]) ([]*models.SignalRecord, error)
}
