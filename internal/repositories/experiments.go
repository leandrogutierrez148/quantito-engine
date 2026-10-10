package repositories

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type experimentsRepository struct {
	db interfaces.IDBHandler
}

func NewExperimentsRepository(db interfaces.IDBHandler) interfaces.IExperimentsRepository {
	return &experimentsRepository{db}
}

func (r *experimentsRepository) Find(includeInactive bool) ([]*models.Experiment, error) {
	var experiments []*models.Experiment
	err := r.db.QueryAll("CALL sp_find_experiments (?)", &experiments, includeInactive)
	if err != nil {
		return nil, err
	}
	return experiments, nil
}

func (r *experimentsRepository) FindByID(experimentID int) (*models.Experiment, error) {
	var experiment models.Experiment
	err := r.db.QueryOne("CALL sp_get_experiment_by_id (?)", &experiment, experimentID)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryAll("CALL sp_find_variants_by_experiment (?)", &experiment.Variants, experimentID)
	if err != nil {
		return nil, err
	}

	return &experiment, nil
}

func (r *experimentsRepository) FindByExperiment(experimentName string) (*models.Experiment, error) {
	var experiment models.Experiment
	err := r.db.QueryOne("CALL sp_get_experiment_by_name (?)", &experiment, experimentName)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryAll("CALL sp_find_variants_by_experiment (?)", &experiment.Variants, experiment.ID)
	if err != nil {
		return nil, err
	}

	return &experiment, nil
}
