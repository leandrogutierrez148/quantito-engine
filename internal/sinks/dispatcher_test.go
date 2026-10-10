package sinks_test

import (
	"quantito-engine/internal/models"
	"quantito-engine/internal/sinks"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockSink struct {
	received []*models.Signal
}

func (m *mockSink) Name() string {
	return "MockSink"
}

func (m *mockSink) Send(signal *models.Signal) error {
	m.received = append(m.received, signal)
	return nil
}

func TestSignalDispatcher(t *testing.T) {
	mock := &mockSink{}
	console := sinks.NewConsoleSink()
	dispatcher := sinks.NewSignalDispatcher(console, mock)

	signal := models.NewSignal("AAPL", models.SignalBuy, 250.75, "fifty_day_avg")
	err := dispatcher.Dispatch(signal)

	assert.NoError(t, err)
	assert.Len(t, mock.received, 1)
	assert.Equal(t, models.SignalBuy, mock.received[0].Type)
	assert.Equal(t, 250.75, mock.received[0].Price)
	assert.Equal(t, "AAPL", mock.received[0].Asset)
}
