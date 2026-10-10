package interfaces

type IExecutionMutex interface {
	Lock(experimentID, variantID int, asset string)
	Unlock(experimentID, variantID int, asset string)
}
