package sinks

import (
	"quantito-engine/internal/interfaces"
	"quantito-engine/internal/models"

	"github.com/rs/zerolog/log"
)

type SignalDispatcher struct {
	sinks []interfaces.ISignalSink
}

func NewSignalDispatcher(sinks ...interfaces.ISignalSink) interfaces.ISignalDispatcher {
	return &SignalDispatcher{
		sinks: sinks,
	}
}

func (d *SignalDispatcher) RegisterSink(sink interfaces.ISignalSink) {
	d.sinks = append(d.sinks, sink)
}

func (d *SignalDispatcher) Dispatch(signal *models.Signal) error {
	for _, sink := range d.sinks {
		err := sink.Send(signal)
		if err != nil {
			log.Error().
				Err(err).
				Str("sink", sink.Name()).
				Int("assetId", signal.AssetID).
				Str("signal", string(signal.Type)).
				Msg("Error dispatching signal to sink")
		}
	}
	return nil
}
