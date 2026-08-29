package infraestructures

import (
	"fmt"
	"quantito-engine/internal/interfaces"

	"github.com/rs/zerolog/log"
	"k8s.io/utils/keymutex"
)

const mutexSize = 64

type executionMutex struct {
	mutex interfaces.IKeyMutex
}

// NewExecutionMutex creates a mutex that locks by the tuple (experimentID, variantID, assetID) to prevent duplicate executions.
func NewExecutionMutex() interfaces.IExecutionMutex {
	return &executionMutex{mutex: keymutex.NewHashed(mutexSize)}
}

func (m *executionMutex) Lock(experimentID, variantID, assetID int) {
	key := m.formatKey(experimentID, variantID, assetID)
	m.mutex.LockKey(key)
}

func (m *executionMutex) Unlock(experimentID, variantID, assetID int) {
	key := m.formatKey(experimentID, variantID, assetID)
	err := m.mutex.UnlockKey(key)
	if err != nil {
		log.Error().Err(err).Str("key", key).Msg("Error unlocking execution mutex")
	}
}

func (m *executionMutex) formatKey(experimentID, variantID, assetID int) string {
	return fmt.Sprintf("%d:%d:%d", experimentID, variantID, assetID)
}
