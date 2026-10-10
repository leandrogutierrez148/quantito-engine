package infraestructures

import (
	"quantito-engine/internal/interfaces/mocks"
	"testing"
)

func TestExecutionMutexLock(t *testing.T) {
	mutexMock := mocks.NewMockIKeyMutex(t)
	mutex := executionMutex{
		mutex: mutexMock,
	}

	expectedKey := "1:10:AAPL"
	mutexMock.EXPECT().LockKey(expectedKey).Once()
	mutexMock.EXPECT().UnlockKey(expectedKey).Return(nil).Once()

	mutex.Lock(1, 10, "AAPL")
	mutex.Unlock(1, 10, "AAPL")
}
