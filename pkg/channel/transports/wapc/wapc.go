package wapc

import (
	"context"
	"runtime"

	wapc "github.com/wapc/wapc-go"

	functions "github.com/nanobus/nanobus/pkg/channel"
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
