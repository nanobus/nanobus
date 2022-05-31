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

package wapc

import (
	"context"
	"runtime"

	wapc "github.com/wapc/wapc-go"

	functions "github.com/nanobus/nanobus/channel"
)

const DefaultPoolSize = 0

type WaPC struct {
	module wapc.Module
	pool   *wapc.Pool
}

// Ensure `Invoke` conforms to `functions.Invoke`
var _ = (functions.Invoke)(((*WaPC)(nil)).Invoke)

// Registering handlers is handled by waPC itself.

func New(module wapc.Module, poolSize uint64) (*WaPC, error) {
	if poolSize == DefaultPoolSize {
		poolSize = uint64(runtime.NumCPU() * 2)
	}
	pool, err := wapc.NewPool(context.Background(), module, poolSize)
	if err != nil {
		return nil, err
	}

	return &WaPC{
		module: module,
		pool:   pool,
	}, nil
}

func (w *WaPC) Invoke(ctx context.Context, receiver functions.Receiver, payload []byte) ([]byte, error) {
	instance, err := w.pool.Get(0)
	if err != nil {
		return nil, err
	}
	defer w.pool.Return(instance)

	return instance.Invoke(ctx, receiver.Namespace+"/"+receiver.Operation, payload)
}
