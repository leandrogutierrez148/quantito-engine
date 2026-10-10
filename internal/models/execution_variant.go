package models

import "time"

type ExecutionVariantStatus string

const (
	ExecutionVariantPending   ExecutionVariantStatus = "PENDING"
	ExecutionVariantRunning   ExecutionVariantStatus = "RUNNING"
	ExecutionVariantCompleted ExecutionVariantStatus = "COMPLETED"
	ExecutionVariantFailed    ExecutionVariantStatus = "FAILED"
)

type ExecutionVariant struct {
	ID           int                 `db:"id" json:"id"`
	ExecutionID  int                 `db:"execution_id" json:"executionId"`
	VariantID    int                 `db:"variant_id" json:"variantId"`
	Status       string              `db:"status" json:"status"`
	SignalsCount int                 `db:"signals_count" json:"signalsCount"`
	DurationMs   int                 `db:"duration_ms" json:"durationMs"`
	ErrorMessage Nullable[string]    `db:"error_message" json:"errorMessage,omitempty"`
	StartedAt    Nullable[time.Time] `db:"started_at" json:"startedAt,omitempty"`
	CompletedAt  Nullable[time.Time] `db:"completed_at" json:"completedAt,omitempty"`
	CreatedAt    time.Time           `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time           `db:"updated_at" json:"updatedAt"`
}
