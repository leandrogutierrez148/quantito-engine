package interfaces

import "quantito-engine/internal/models"

type IExperimentsRepository interface {
	Find(includeInactive bool) ([]*models.Experiment, error)
	FindByID(experimentID int) (*models.Experiment, error)
	FindByExperiment(experiment string) (*models.Experiment, error)
}
