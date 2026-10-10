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

func (r *signalsRepository) Save(executionVariantID int, asset string, signal *models.Signal) (*models.SignalRecord, error) {
	var result models.SignalRecord

	confidence := signal.Confidence
	if confidence == 0 {
		confidence = 1.0
	}

	err := r.db.QueryOne(
		"CALL sp_save_signal (?, ?, ?, ?, ?)",
		&result,
		executionVariantID,
		asset,
		string(signal.Type),
		signal.Price,
		confidence,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *signalsRepository) FindByExecution(executionID int) ([]*models.SignalRecord, error) {
	var results []*models.SignalRecord
	err := r.db.QueryAll("CALL sp_find_signals_by_execution (?)", &results, executionID)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *signalsRepository) FindByAsset(asset string, fromDate, toDate models.Nullable[time.Time]) ([]*models.SignalRecord, error) {
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
		"CALL sp_find_signals_by_asset (?, ?, ?)",
		&results,
		asset,
		fromVal,
		toVal,
	)
	if err != nil {
		return nil, err
	}
	return results, nil
}
