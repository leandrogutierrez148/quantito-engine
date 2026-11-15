package services

import (
	"fmt"
	"time"

	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"

	"github.com/rs/zerolog/log"
)

// DefaultMaxAttempts defines the default maximum number of evaluation attempts per variant.
const DefaultMaxAttempts = 3

// ExecutionService orchestrates experiment executions across their configured assets and competing strategy variants.
type ExecutionService struct {
	experimentsRepo  interfaces.IExperimentsRepository
	executionsRepo   interfaces.IExecutionsRepository
	signalsRepo      interfaces.ISignalsRepository
	strategyFactory  interfaces.IStrategyFactory
	signalDispatcher interfaces.ISignalDispatcher
	executionMutex   interfaces.IExecutionMutex
}

// NewExecutionService creates a new ExecutionService with all required dependencies.
func NewExecutionService(
	experimentsRepo interfaces.IExperimentsRepository,
	executionsRepo interfaces.IExecutionsRepository,
	signalsRepo interfaces.ISignalsRepository,
	strategyFactory interfaces.IStrategyFactory,
	signalDispatcher interfaces.ISignalDispatcher,
	executionMutex interfaces.IExecutionMutex,
) interfaces.IExecutionService {
	return &ExecutionService{
		experimentsRepo:  experimentsRepo,
		executionsRepo:   executionsRepo,
		signalsRepo:      signalsRepo,
		strategyFactory:  strategyFactory,
		signalDispatcher: signalDispatcher,
		executionMutex:   executionMutex,
	}
}

// Execute loads an experiment by its ID, validates its active state, and orchestrates its execution.
func (e *ExecutionService) Execute(id int) error {
	return e.ExecuteWithDetails(id, string(models.TriggerManual), "system", time.Now().UTC())
}

// ExecuteWithDetails runs an experiment with explicit trigger metadata and asOfTime.
func (e *ExecutionService) ExecuteWithDetails(id int, triggerType, triggeredBy string, asOfTime time.Time) error {
	experiment, err := e.experimentsRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding experiment %d: %w", id, err)
	}

	if experiment == nil || !experiment.IsActive() {
		log.Warn().Int("id", id).Msg("Experiment is inactive or not found")
		return nil
	}

	if len(experiment.Variants) == 0 {
		log.Warn().Int("experimentId", experiment.ID).Msg("No variants configured for experiment")
		return nil
	}

	if asOfTime.IsZero() {
		asOfTime = time.Now().UTC()
	}

	execution, err := e.executionsRepo.CreateExecution(experiment.ID, triggerType, triggeredBy, asOfTime)
	if err != nil {
		return fmt.Errorf("error creating execution for experiment %d: %w", experiment.ID, err)
	}

	return e.runExecution(experiment, execution)
}

func (e *ExecutionService) runExecution(experiment *models.Experiment, execution *models.Execution) error {
	log.Info().
		Int("experimentId", experiment.ID).
		Int("executionId", execution.ID).
		Str("experiment", experiment.Name).
		Msg("Starting execution for experiment")

	totalVariants := len(experiment.Variants)
	completedCount := 0
	failedCount := 0

	for _, variant := range experiment.Variants {
		snapshotParams := e.mergeVariantParams(experiment, variant)

		ev, err := e.executionsRepo.CreateExecutionVariant(execution.ID, variant.ID, snapshotParams)
		if err != nil {
			log.Error().Err(err).Int("executionId", execution.ID).Int("variantId", variant.ID).Msg("Error creating execution variant record")
			failedCount++
			continue
		}

		maxAttempts := e.getMaxAttempts(variant)
		attempts := 0
		start := time.Now()
		var evalErr error
		for attempts < maxAttempts {
			attempts++
			if attempts > 1 {
				log.Warn().
					Int("executionId", execution.ID).
					Int("variantId", variant.ID).
					Int("attempt", attempts).
					Int("maxAttempts", maxAttempts).
					Msg("Retrying variant evaluation")
			}
			evalErr = e.evaluateVariant(ev.ID, experiment, variant, snapshotParams)
			if evalErr == nil {
				break
			}
		}
		durationMs := int(time.Since(start).Milliseconds())

		if evalErr != nil {
			log.Error().Err(evalErr).Int("executionId", execution.ID).Int("variantId", variant.ID).Int("attempts", attempts).Msg("Error evaluating variant")
			errStr := evalErr.Error()
			_, _ = e.executionsRepo.UpdateExecutionVariant(ev.ID, string(models.ExecutionVariantFailed), attempts, durationMs, models.NotNull(errStr))
			failedCount++
		} else {
			_, _ = e.executionsRepo.UpdateExecutionVariant(ev.ID, string(models.ExecutionVariantCompleted), attempts, durationMs, models.Nullable[string]{})
			completedCount++
		}
	}

	finalStatus := string(models.ExecutionCompleted)
	var finalErrMsg models.Nullable[string]
	if failedCount > 0 && completedCount > 0 {
		finalStatus = string(models.ExecutionPartialSuccess)
		msg := fmt.Sprintf("%d of %d variants failed during execution", failedCount, totalVariants)
		finalErrMsg = models.NotNull(msg)
	} else if failedCount == totalVariants {
		finalStatus = string(models.ExecutionFailed)
		msg := "All variants failed during execution"
		finalErrMsg = models.NotNull(msg)
	}

	_, updateErr := e.executionsRepo.UpdateExecution(execution.ID, finalStatus, finalErrMsg)
	if updateErr != nil {
		log.Error().Err(updateErr).Int("executionId", execution.ID).Msg("Error updating execution final status")
		return updateErr
	}

	log.Info().
		Int("experimentId", experiment.ID).
		Int("executionId", execution.ID).
		Str("status", finalStatus).
		Int("completedVariants", completedCount).
		Int("failedVariants", failedCount).
		Msg("Finished execution for experiment")

	return nil
}

func (e *ExecutionService) getMaxAttempts(variant *models.Variant) int {
	if variant != nil && variant.MaxAttempts >= 1 {
		return variant.MaxAttempts
	}
	return DefaultMaxAttempts
}

func (e *ExecutionService) mergeVariantParams(experiment *models.Experiment, variant *models.Variant) map[string]any {
	params := make(map[string]any)
	if experiment.Params.Get() != nil {
		for k, v := range *experiment.Params.Get() {
			params[k] = v
		}
	}
	if variant.Params.Get() != nil {
		for k, v := range *variant.Params.Get() {
			params[k] = v
		}
	}
	return params
}

func (e *ExecutionService) evaluateVariant(executionVariantID int, experiment *models.Experiment, variant *models.Variant, params map[string]any) error {
	experimentID := experiment.ID
	variantID := variant.ID
	asset := variant.Asset

	strategy, err := e.strategyFactory.CreateStrategy(variant.Class)
	if err != nil {
		return fmt.Errorf("error creating strategy %s: %w", variant.Class, err)
	}

	e.executionMutex.Lock(experimentID, variantID, asset)
	signal, err := strategy.GenerateSignal(asset, params)
	e.executionMutex.Unlock(experimentID, variantID, asset)

	if err != nil {
		return fmt.Errorf("error evaluating asset %s for variant %d: %w", asset, variantID, err)
	}

	if signal != nil {
		log.Info().
			Str("asset", asset).
			Int("experimentId", experimentID).
			Int("variantId", variantID).
			Str("strategy", variant.Class).
			Str("signal", string(signal.Type)).
			Float64("price", signal.Price).
			Msg("Strategy generated signal")

		_, err = e.signalsRepo.Save(executionVariantID, asset, signal)
		if err != nil {
			return fmt.Errorf("error saving signal for asset %s: %w", asset, err)
		}

		if err := e.signalDispatcher.Dispatch(signal); err != nil {
			log.Warn().Err(err).Str("asset", asset).Msg("Error dispatching signal")
		}
	}

	return nil
}
