package models_test

import (
	"quantito-engine/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExperiment_IsActive(t *testing.T) {
	expActive := models.Experiment{Status: "ACTIVE"}
	assert.True(t, expActive.IsActive())

	expA := models.Experiment{Status: "A"}
	assert.True(t, expA.IsActive())

	expPaused := models.Experiment{Status: "PAUSED"}
	assert.False(t, expPaused.IsActive())
}

func TestExperiment_GetParam(t *testing.T) {
	exp := models.Experiment{
		Params: models.NewJsonColumn(map[string]any{
			"execution_mode": "PAPER",
			"max_allocation": 1000.0,
		}),
	}

	mode, ok := exp.GetParam("execution_mode")
	assert.True(t, ok)
	assert.Equal(t, "PAPER", mode)

	_, ok = exp.GetParam("non_existing")
	assert.False(t, ok)
}

func TestExperiment_GetAssets(t *testing.T) {
	exp := models.Experiment{
		Assets: models.NewJsonColumn([]string{"AAPL", "NVDA"}),
	}

	assets := exp.GetAssets()
	assert.Equal(t, []string{"AAPL", "NVDA"}, assets)

	var nilExp *models.Experiment
	assert.Nil(t, nilExp.GetAssets())
}
