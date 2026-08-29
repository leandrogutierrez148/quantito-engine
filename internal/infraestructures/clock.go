package infraestructures

import (
	"quantito-engine/internal/interfaces"
	"time"
)

type clock struct{}

func NewClock() interfaces.IClock {
	return &clock{}
}

func (c *clock) Now() time.Time {
	return time.Now()
}
