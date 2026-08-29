package interfaces

import "quantito-engine/internal/models"

type ISignalSink interface {
	Name() string
	Send(signal *models.Signal) error
}

type ISignalDispatcher interface {
	RegisterSink(sink ISignalSink)
	Dispatch(signal *models.Signal) error
}
