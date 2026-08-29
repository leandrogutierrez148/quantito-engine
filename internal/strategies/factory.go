package strategies

import (
	"fmt"
	"quantito-engine/internal/interfaces"
)

type StrategyFactory struct{}

func NewStrategyFactory() interfaces.IStrategyFactory {
	return &StrategyFactory{}
}

func (f *StrategyFactory) CreateStrategy(strategyClass string) (interfaces.IStrategy, error) {
	switch strategyClass {
	case "fifty_day_avg":
		return NewFiftyDayAvgStrategy(), nil
	case "price_momentum":
		return NewPriceMomentumStrategy(), nil
	default:
		return nil, fmt.Errorf("strategy class '%s' is not implemented", strategyClass)
	}
}
