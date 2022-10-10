package stream

import (
	"github.com/WasmRS/wasmrs-go/payload"
	"github.com/WasmRS/wasmrs-go/rx/flux"
	"github.com/nanobus/nanobus/channel/metadata"
	"github.com/vmihailenco/msgpack/v5"
)

type sink struct {
	f flux.Sink[payload.Payload]
}

func SinkFromFlux(f flux.Sink[payload.Payload]) Sink {
	return &sink{
		f: f,
	}
}

func (s *sink) Next(data any, md metadata.MD) error {
	dataBytes, err := msgpack.Marshal(data)
	if err != nil {
		return err
	}
	s.f.Next(payload.New(dataBytes))
	return nil
}

func (s *sink) Complete() {
	s.f.Complete()
}

func (s *sink) Error(err error) {
	s.f.Error(err)
}
