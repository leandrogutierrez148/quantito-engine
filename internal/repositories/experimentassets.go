package repositories

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type experimentAssetsRepository struct {
	db interfaces.IDBHandler
}

func NewExperimentAssetsRepository(db interfaces.IDBHandler) interfaces.IExperimentAssetsRepository {
	return &experimentAssetsRepository{db: db}
}

// FindByExperiment returns all assets associated with an experiment.
func (r *experimentAssetsRepository) FindByExperiment(experimentID int) ([]*models.ExperimentAsset, error) {
	var assets []*models.ExperimentAsset
	err := r.db.QueryAll("CALL sp_find_experiment_assets (?)", &assets, experimentID)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

// FindExperimentsByAsset returns all active experiments that contain the given asset.
func (r *experimentAssetsRepository) FindExperimentsByAsset(assetID int) ([]*models.Experiment, error) {
	var experiments []*models.Experiment
	err := r.db.QueryAll("CALL sp_find_experiments_by_asset (?)", &experiments, assetID)
	if err != nil {
		return nil, err
	}
	return experiments, nil
}
