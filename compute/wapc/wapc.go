package wapc

import (
	"context"
	"fmt"
	"os"
	go_runtime "runtime"
	"strings"

	"github.com/nanobus/go-functions"
	wapc_mux "github.com/nanobus/go-functions/transports/wapc"
	wapc "github.com/wapc/wapc-go"

	"github.com/nanobus/nanobus/compute"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

type WaPCConfig struct {
	// Filename is the file name of the waPC/WebAssembly module to load.
	Filename string `mapstructure:"filename"` // TODO: Load from external location
	// PoolSize controls the number of waPC instance of the module to create and pool.
	// It also represents the maximum number of concurrent requests the module can process.
	PoolSize uint64 `mapstructure:"poolSize"`
}

// Mux is the NamedLoader for the waPC compute.
func WaPC() (string, compute.Loader) {
	return "wapc", WaPCLoader
}

func WaPCLoader(with interface{}, resolver resolve.ResolveAs) (*functions.Invoker, error) {
	var busInvoker compute.BusInvoker
	var msgpackcodec functions.Codec
	if err := resolve.Resolve(resolver,
		"bus:invoker", &busInvoker,
		"codec:msgpack", &msgpackcodec); err != nil {
		return nil, err
	}

	c := WaPCConfig{
		PoolSize: uint64(go_runtime.NumCPU() * 5),
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	wasmBytes, err := os.ReadFile(c.Filename)
	if err != nil {
		return nil, err
	}

	module, err := wapc.New(wasmBytes, func(ctx context.Context, binding, namespace, operation string, payload []byte) ([]byte, error) {
		lastDot := strings.LastIndexByte(namespace, '.')
		if lastDot < 0 {
			return nil, fmt.Errorf("invalid namespace %q", namespace)
		}
		service := namespace[lastDot+1:]
		namespace = namespace[:lastDot]

		var input interface{}
		if err := msgpackcodec.Decode(payload, &input); err != nil {
			return nil, err
		}

		result, err := busInvoker(ctx, namespace, service, operation, input)
		if err != nil {
			return nil, err
		}

		return msgpackcodec.Encode(result)
	})
	if err != nil {
		return nil, err
	}

	module.SetLogger(wapc.Println)
	module.SetWriter(wapc.Print)

	m, err := wapc_mux.New(module, uint64(c.PoolSize))
	if err != nil {
		return nil, err
	}
	invoker := functions.NewInvoker(m.Invoke, msgpackcodec)

	return invoker, nil
}
