package services

import (
	"fmt"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"

	"github.com/rs/zerolog/log"
)

// ExecutionService orchestrates experiment executions across their assigned assets and competing strategy variants.
type ExecutionService struct {
	experimentsRepo  interfaces.IExperimentsRepository
	signalsRepo      interfaces.ISignalsRepository
	strategyFactory  interfaces.IStrategyFactory
	signalDispatcher interfaces.ISignalDispatcher
	executionMutex   interfaces.IExecutionMutex
}

// NewExecutionService creates a new ExecutionService with all required dependencies.
func NewExecutionService(
	experimentsRepo interfaces.IExperimentsRepository,
	signalsRepo interfaces.ISignalsRepository,
	strategyFactory interfaces.IStrategyFactory,
	signalDispatcher interfaces.ISignalDispatcher,
	executionMutex interfaces.IExecutionMutex,
) interfaces.IExecutionService {
	return &ExecutionService{
		experimentsRepo:  experimentsRepo,
		signalsRepo:      signalsRepo,
		strategyFactory:  strategyFactory,
		signalDispatcher: signalDispatcher,
		executionMutex:   executionMutex,
	}
}

// Execute loads an experiment by its ID, validates its active state, and orchestrates its execution.
func (e *ExecutionService) Execute(id int) error {
	experiment, err := e.experimentsRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("error finding experiment %d: %w", id, err)
	}

	if experiment == nil || !experiment.IsActive() {
		log.Warn().Int("id", id).Msg("Experiment is inactive or not found")
		return nil
	}

	return e.executeExperiment(experiment)
}

// executeExperiment executes all variants of an experiment concurrently across its assigned assets.
func (e *ExecutionService) executeExperiment(experiment *models.Experiment) error {
	log.Info().
		Int("experimentId", experiment.ID).
		Str("experiment", experiment.Experiment).
		Msg("Starting execution for experiment")

	if len(experiment.Assets) == 0 {
		log.Warn().Int("experimentId", experiment.ID).Msg("No active assets assigned to experiment")
		return nil
	}

	if len(experiment.Variants) == 0 {
		log.Warn().Int("experimentId", experiment.ID).Msg("No variants configured for experiment")
		return nil
	}

	for _, ea := range experiment.Assets {
		err := e.evaluateAssetAcrossVariants(ea.AssetID, experiment.ID, experiment.Variants)
		if err != nil {
			log.Error().
				Err(err).
				Int("experimentId", experiment.ID).
				Int("assetId", ea.AssetID).
				Msg("Error evaluating asset across variants")
			return err
		}
	}

	log.Info().
		Int("experimentId", experiment.ID).
		Int("assetsCount", len(experiment.Assets)).
		Int("variantsCount", len(experiment.Variants)).
		Msg("Finished execution for experiment")

	return nil
}

// evaluateAssetAcrossVariants evaluates each competing strategy variant for an asset.
func (e *ExecutionService) evaluateAssetAcrossVariants(assetID, experimentID int, variants []*models.Variant) error {
	for _, variant := range variants {
		err := e.executeVariant(assetID, experimentID, variant)
		if err != nil {
			return err
		}
	}

	return nil
}

// executeVariant evaluates a single strategy algorithm under a mutex lock on (experiment, variant, asset), saves the generated signal, and dispatches it.
func (e *ExecutionService) executeVariant(assetID, experimentID int, variant *models.Variant) error {
	variantID := variant.ID

	e.executionMutex.Lock(experimentID, variantID, assetID)
	defer e.executionMutex.Unlock(experimentID, variantID, assetID)

	strategy, err := e.strategyFactory.CreateStrategy(variant.Class)
	if err != nil {
		return err
	}

	params := make(map[string]any)
	if variant.Params.Get() != nil {
		params = *variant.Params.Get()
	}

	signal, err := strategy.GenerateSignal(assetID, params)
	if err != nil {
		return err
	}

	log.Info().
		Int("assetId", assetID).
		Int("experimentId", experimentID).
		Int("variantId", variantID).
		Str("strategy", variant.Class).
		Str("signal", string(signal.Type)).
		Float64("price", signal.Price).
		Msg("Strategy generated signal")

	_, err = e.signalsRepo.Save(assetID, variantID, signal)
	if err != nil {
		return err
	}

	return e.signalDispatcher.Dispatch(signal)
}
