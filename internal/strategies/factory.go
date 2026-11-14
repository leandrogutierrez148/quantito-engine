package strategies

import (
	"fmt"
	"quantito-engine/internal/infraestructures"
	"quantito-engine/internal/interfaces"
)

type StrategyFactory struct {
	apiClient interfaces.IAPIClient
}

func NewStrategyFactory(client ...interfaces.IAPIClient) interfaces.IStrategyFactory {
	var c interfaces.IAPIClient
	if len(client) > 0 && client[0] != nil {
		c = client[0]
	} else {
		c = infraestructures.NewAPIClient("", nil)
	}
	return &StrategyFactory{
		apiClient: c,
	}
}

func (f *StrategyFactory) CreateStrategy(strategyClass string) (interfaces.IStrategy, error) {
	switch strategyClass {
	case "fifty_day_avg":
		return NewFiftyDayAvgStrategy(f.apiClient), nil
	case "price_momentum":
		return NewPriceMomentumStrategy(), nil
	case "clockwork_hgb":
		return NewClockworkHGBStrategy(f.apiClient), nil
	default:
		return nil, fmt.Errorf("strategy class '%s' is not implemented", strategyClass)
	}
}
