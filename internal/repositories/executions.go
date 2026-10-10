package repositories

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
	"time"
)

type executionsRepository struct {
	db interfaces.IDBHandler
}

func NewExecutionsRepository(db interfaces.IDBHandler) interfaces.IExecutionsRepository {
	return &executionsRepository{db: db}
}

func (r *executionsRepository) CreateExecution(experimentID int, triggerType, triggeredBy string, asOfTime time.Time) (*models.Execution, error) {
	var exec models.Execution
	err := r.db.QueryOne(
		"CALL sp_create_execution (?, ?, ?, ?)",
		&exec,
		experimentID,
		triggerType,
		triggeredBy,
		asOfTime,
	)
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func (r *executionsRepository) UpdateExecution(executionID int, status string, errorMessage models.Nullable[string]) (*models.Execution, error) {
	var exec models.Execution
	var errVal any
	if errorMessage.Valid {
		errVal = errorMessage.Val
	}

	err := r.db.QueryOne(
		"CALL sp_update_execution (?, ?, ?)",
		&exec,
		executionID,
		status,
		errVal,
	)
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func (r *executionsRepository) CreateExecutionVariant(executionID, variantID int) (*models.ExecutionVariant, error) {
	var ev models.ExecutionVariant
	err := r.db.QueryOne(
		"CALL sp_create_execution_variant (?, ?)",
		&ev,
		executionID,
		variantID,
	)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *executionsRepository) UpdateExecutionVariant(id int, status string, signalsCount, durationMs int, errorMessage models.Nullable[string]) (*models.ExecutionVariant, error) {
	var ev models.ExecutionVariant
	var errVal any
	if errorMessage.Valid {
		errVal = errorMessage.Val
	}

	err := r.db.QueryOne(
		"CALL sp_update_execution_variant (?, ?, ?, ?, ?)",
		&ev,
		id,
		status,
		signalsCount,
		durationMs,
		errVal,
	)
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *executionsRepository) FindByID(executionID int) (*models.Execution, error) {
	var exec models.Execution
	err := r.db.QueryOne(
		"SELECT id, experiment_id, trigger_type, triggered_by, as_of_time, status, started_at, completed_at, error_message, created_at, updated_at FROM executions WHERE id = ?",
		&exec,
		executionID,
	)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryAll(
		"SELECT id, execution_id, variant_id, status, signals_count, duration_ms, error_message, started_at, completed_at, created_at, updated_at FROM execution_variants WHERE execution_id = ?",
		&exec.Variants,
		executionID,
	)
	if err != nil {
		return nil, err
	}

	return &exec, nil
}

func (r *executionsRepository) FindByExperiment(experimentID int, limit int) ([]*models.Execution, error) {
	if limit <= 0 {
		limit = 20
	}
	var executions []*models.Execution
	err := r.db.QueryAll(
		"SELECT id, experiment_id, trigger_type, triggered_by, as_of_time, status, started_at, completed_at, error_message, created_at, updated_at FROM executions WHERE experiment_id = ? ORDER BY id DESC LIMIT ?",
		&executions,
		experimentID,
		limit,
	)
	if err != nil {
		return nil, err
	}
	return executions, nil
}
