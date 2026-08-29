package interfaces

import (
	"quantito-engine/internal/models"
)

type IStrategy interface {
	GenerateSignal(assetID int, params map[string]any) (*models.Signal, error)
}
