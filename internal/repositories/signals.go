package repositories

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
	"time"
)

type signalsRepository struct {
	db interfaces.IDBHandler
}

func NewSignalsRepository(db interfaces.IDBHandler) interfaces.ISignalsRepository {
	return &signalsRepository{db}
}

func (r *signalsRepository) Save(assetID, variantID int, signal *models.Signal) (*models.SignalRecord, error) {
	var result models.SignalRecord

	confidence := signal.Confidence
	if confidence == 0 {
		confidence = 1.0
	}

	err := r.db.QueryOne(
		"CALL sp_save_signal (?, ?, ?, ?, ?)",
		&result,
		assetID,
		variantID,
		string(signal.Type),
		signal.Price,
		confidence,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *signalsRepository) Find(assetID int, fromDate, toDate models.Nullable[time.Time]) ([]*models.SignalRecord, error) {
	var results []*models.SignalRecord

	var fromVal any
	if fromDate.Valid {
		fromVal = fromDate.Val
	}

	var toVal any
	if toDate.Valid {
		toVal = toDate.Val
	}

	err := r.db.QueryAll(
		"CALL sp_find_signals (?, ?, ?)",
		&results,
		assetID,
		fromVal,
		toVal,
	)
	if err != nil {
		return nil, err
	}

	return results, nil
}
