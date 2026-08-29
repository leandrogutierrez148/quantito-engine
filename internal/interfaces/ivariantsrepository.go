package interfaces

import "quantito-engine/internal/models"

type IVariantsRepository interface {
	FindByExperiment(experimentID int) ([]*models.Variant, error)
}
