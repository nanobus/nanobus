/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package channel

import (
	"context"
	"fmt"
	"io"
	"reflect"

	"github.com/nanobus/nanobus/channel/metadata"
)

type (
	Receiver struct {
		Namespace string
		Operation string
		EntityID  string
	}

	///////////////////////////////////////////////////////
	// Function handling

	// Register registers a function handler with the transport layer.
	Register func(namespace, operation string, handler Handler)
	// Handler is a function that handles the invocation of a named function.
	Handler func(ctx context.Context, payload []byte) ([]byte, error)

	RegisterStateful func(namespace, operation string, method StatefulHandler)
	StatefulHandler  func(ctx context.Context, id string, payload []byte) ([]byte, error)

	///////////////////////////////////////////////////////
	// Function invoking

	// Invoke calls a function over the transport layer.
	Invoke func(ctx context.Context, receiver Receiver, payload []byte) ([]byte, error)

	InvokeStream func(ctx context.Context, receiver Receiver) (Streamer, error)

	Streamer interface {
		SendMetadata(md metadata.MD, end ...bool) error
		SendData(data []byte, end ...bool) error
		Metadata() metadata.MD
		RecvData() ([]byte, error)
		Close() error
	}
)

// Codec is an interface that handles encoding and decoding payloads send to and
// received from functions.
type Codec interface {
	ContentType() string
	Encode(v interface{}) ([]byte, error)
	Decode(data []byte, v interface{}) error
}

type Invoker struct {
	invoke       Invoke
	invokeStream InvokeStream
	codec        Codec
}

func NewInvoker(invoke Invoke, invokeStream InvokeStream, codec Codec) *Invoker {
	return &Invoker{
		invoke:       invoke,
		invokeStream: invokeStream,
		codec:        codec,
	}
}

func (i *Invoker) Invoke(ctx context.Context, receiver Receiver, input interface{}) error {
	reqBytes, err := i.codec.Encode(input)
	if err != nil {
		return err
	}

	_, err = i.invoke(ctx, receiver, reqBytes)
	return err
}

func (i *Invoker) InvokeWithReturn(ctx context.Context, receiver Receiver, input, output interface{}) error {
	reqBytes, err := i.codec.Encode(input)
	if err != nil {
		return err
	}

	respBytes, err := i.invoke(ctx, receiver, reqBytes)
	if err != nil {
		return err
	}

	if len(respBytes) == 0 {
		v := reflect.ValueOf(output)
		t := v.Type()
		if t.Kind() == reflect.Ptr && !v.IsNil() {
			p := v.Elem()
			p.Set(reflect.Zero(p.Type()))
		}

		return nil
	}

	return i.codec.Decode(respBytes, output)
}

func (i *Invoker) InvokeStream(ctx context.Context, receiver Receiver) (*Stream, error) {
	streamer, err := i.invokeStream(ctx, receiver)
	if err != nil {
		return nil, err
	}
	return &Stream{
		streamer: streamer,
		codec:    i.codec,
	}, nil
}

type Stream struct {
	streamer Streamer
	codec    Codec
}

func (s *Stream) Metadata() metadata.MD {
	return s.streamer.Metadata()
}

func (s *Stream) RecvData(dst interface{}) error {
	msg, err := s.streamer.RecvData()
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}
		return fmt.Errorf("could not receive data: %w", err)
	}
	if len(msg) == 0 {
		return io.EOF
	}

	return s.codec.Decode(msg, dst)
}

func (s *Stream) SendHeaders(md metadata.MD) error {
	return s.streamer.SendMetadata(md, false)
}

func (s *Stream) SendData(v interface{}, end ...bool) error {
	var endVal bool
	if len(end) > 0 {
		endVal = end[0]
	}
	var vBytes []byte
	var err error
	if !isNil(v) {
		vBytes, err = s.codec.Encode(v)
		if err != nil {
			return fmt.Errorf("could not marshal value to send: %w", err)
		}
	} else {
		vBytes = []byte{}
	}
	return s.streamer.SendData(vBytes, endVal)
}

func (s *Stream) Close() error {
	return s.streamer.Close()
}

func isNil(val interface{}) bool {
	return val == nil ||
		(reflect.ValueOf(val).Kind() == reflect.Ptr &&
			reflect.ValueOf(val).IsNil())
}
