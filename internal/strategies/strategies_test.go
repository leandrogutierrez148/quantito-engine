package strategies_test

import (
	"quantito-engine/internal/models"
	"quantito-engine/internal/strategies"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrategyFactory(t *testing.T) {
	factory := strategies.NewStrategyFactory()

	strat, err := factory.CreateStrategy("fifty_day_avg")
	assert.NoError(t, err)
	assert.NotNil(t, strat)

	strat, err = factory.CreateStrategy("price_momentum")
	assert.NoError(t, err)
	assert.NotNil(t, strat)

	_, err = factory.CreateStrategy("unknown")
	assert.Error(t, err)
}

func TestFiftyDayAvgStrategy(t *testing.T) {
	strat := strategies.NewFiftyDayAvgStrategy()

	// Price > Avg -> Buy
	signal, err := strat.GenerateSignal(101, map[string]any{
		"current_price": 110.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalBuy, signal.Type)
	assert.Equal(t, 110.0, signal.Price)

	// Price < Avg -> Sell
	signal, err = strat.GenerateSignal(101, map[string]any{
		"current_price": 90.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalSell, signal.Type)

	// Price close to Avg -> Hold
	signal, err = strat.GenerateSignal(101, map[string]any{
		"current_price": 101.0,
		"fifty_day_avg": 100.0,
		"threshold":     0.05,
	})
	assert.NoError(t, err)
	assert.Equal(t, models.SignalHold, signal.Type)
}
