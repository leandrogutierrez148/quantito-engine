package interfaces

type IStrategyFactory interface {
	CreateStrategy(strategyClass string) (IStrategy, error)
}
