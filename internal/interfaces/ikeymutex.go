package interfaces

type IKeyMutex interface {
	LockKey(key string)
	UnlockKey(key string) error
}
