package repositories

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type variantsRepository struct {
	db interfaces.IDBHandler
}

func NewVariantsRepository(db interfaces.IDBHandler) interfaces.IVariantsRepository {
	return &variantsRepository{db}
}

func (r *variantsRepository) FindByExperiment(experimentID int) ([]*models.Variant, error) {
	var variants []*models.Variant
	err := r.db.QueryAll("CALL sp_find_variants_by_experiment (?)", &variants, experimentID)
	if err != nil {
		return nil, err
	}
	return variants, nil
}
