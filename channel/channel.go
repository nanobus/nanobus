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
