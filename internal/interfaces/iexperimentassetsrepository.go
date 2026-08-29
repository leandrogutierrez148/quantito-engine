package interfaces

import "quantito-engine/internal/models"

type IExperimentAssetsRepository interface {
	FindByExperiment(experimentID int) ([]*models.ExperimentAsset, error)
	FindExperimentsByAsset(assetID int) ([]*models.Experiment, error)
}
