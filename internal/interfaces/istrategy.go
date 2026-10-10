package interfaces

import (
	"quantito-engine/internal/models"
)

type IStrategy interface {
	GenerateSignal(asset string, params map[string]any) (*models.Signal, error)
}
