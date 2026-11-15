package sinks_test

import (
	"quantito-engine/internal/interfaces/mocks"
	"quantito-engine/internal/models"
	"quantito-engine/internal/sinks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSignalDispatcher(t *testing.T) {
	mockSink := mocks.NewMockISignalSink(t)
	console := sinks.NewConsoleSink()
	dispatcher := sinks.NewSignalDispatcher(console, mockSink)

	signal := models.NewSignal("AAPL", models.SignalBuy, 250.75, "fifty_day_avg")
	mockSink.EXPECT().Send(mock.MatchedBy(func(s *models.Signal) bool {
		return s.Asset == "AAPL" && s.Type == models.SignalBuy && s.Price == 250.75
	})).Return(nil).Once()

	err := dispatcher.Dispatch(signal)

	assert.NoError(t, err)
}
