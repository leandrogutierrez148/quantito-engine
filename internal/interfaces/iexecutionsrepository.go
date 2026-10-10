package interfaces

import (
	"quantito-engine/internal/models"
	"time"
)

type IExecutionsRepository interface {
	CreateExecution(experimentID int, triggerType, triggeredBy string, asOfTime time.Time) (*models.Execution, error)
	UpdateExecution(executionID int, status string, errorMessage models.Nullable[string]) (*models.Execution, error)
	CreateExecutionVariant(executionID, variantID int) (*models.ExecutionVariant, error)
	UpdateExecutionVariant(id int, status string, signalsCount, durationMs int, errorMessage models.Nullable[string]) (*models.ExecutionVariant, error)
	FindByID(executionID int) (*models.Execution, error)
	FindByExperiment(experimentID int, limit int) ([]*models.Execution, error)
}
