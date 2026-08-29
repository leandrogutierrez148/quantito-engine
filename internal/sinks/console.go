package sinks

import (
	"fmt"
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"
)

type ConsoleSink struct{}

func NewConsoleSink() interfaces.ISignalSink {
	return &ConsoleSink{}
}

func (s *ConsoleSink) Name() string {
	return "ConsoleSink"
}

func (s *ConsoleSink) Send(signal *models.Signal) error {
	fmt.Printf("[SIGNAL DISPATCHED] Asset: %d | Strategy: %s | Signal: %s | Price: %.4f | Timestamp: %s\n",
		signal.AssetID, signal.StrategyName, signal.Type, signal.Price, signal.Timestamp.Format("2006-01-02 15:04:05"))
	return nil
}
