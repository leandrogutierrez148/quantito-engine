package models

import "time"

type ExecutionTriggerType string

const (
	TriggerManual    ExecutionTriggerType = "MANUAL"
	TriggerScheduled ExecutionTriggerType = "SCHEDULED"
	TriggerBacktest  ExecutionTriggerType = "BACKTEST"
	TriggerWebhook   ExecutionTriggerType = "WEBHOOK"
)

type ExecutionStatus string

const (
	ExecutionPending        ExecutionStatus = "PENDING"
	ExecutionRunning        ExecutionStatus = "RUNNING"
	ExecutionCompleted      ExecutionStatus = "COMPLETED"
	ExecutionPartialSuccess ExecutionStatus = "PARTIAL_SUCCESS"
	ExecutionFailed         ExecutionStatus = "FAILED"
)

type Execution struct {
	ID           int                 `db:"id" json:"id"`
	ExperimentID int                 `db:"experiment_id" json:"experimentId"`
	TriggerType  string              `db:"trigger_type" json:"triggerType"`
	TriggeredBy  string              `db:"triggered_by" json:"triggeredBy"`
	AsOfTime     time.Time           `db:"as_of_time" json:"asOfTime"`
	Status       string              `db:"status" json:"status"`
	StartedAt    time.Time           `db:"started_at" json:"startedAt"`
	CompletedAt  Nullable[time.Time] `db:"completed_at" json:"completedAt,omitempty"`
	ErrorMessage Nullable[string]    `db:"error_message" json:"errorMessage,omitempty"`
	CreatedAt    time.Time           `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time           `db:"updated_at" json:"updatedAt"`
	Variants     []*ExecutionVariant `db:"-" json:"variants,omitempty"`
}
