package infraestructures

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClockNow(t *testing.T) {
	c := NewClock()

	assert.NotNil(t, c.Now())
}
