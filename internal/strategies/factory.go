package strategies

import (
	"fmt"
	"quantito-engine/internal/interfaces"
)

type StrategyFactory struct {
	quantitoClient interfaces.IAPIClient
}

func NewStrategyFactory(client ...interfaces.IAPIClient) interfaces.IStrategyFactory {
	var c interfaces.IAPIClient
	if len(client) > 0 {
		c = client[0]
	}
	return &StrategyFactory{
		quantitoClient: c,
	}
}

func (f *StrategyFactory) CreateStrategy(strategyClass string) (interfaces.IStrategy, error) {
	switch strategyClass {
	case "fifty_day_avg":
		return NewFiftyDayAvgStrategy(f.quantitoClient), nil
	case "price_momentum":
		return NewPriceMomentumStrategy(), nil
	default:
		return nil, fmt.Errorf("strategy class '%s' is not implemented", strategyClass)
	}
}
