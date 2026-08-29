package interfaces

type IExecutionMutex interface {
	Lock(experimentID, variantID, assetID int)
	Unlock(experimentID, variantID, assetID int)
}
