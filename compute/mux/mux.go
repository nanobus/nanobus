package mux

import (
	"os"

	"github.com/nanobus/go-functions"
	msgpack_codec "github.com/nanobus/go-functions/codecs/msgpack"
	transport_mux "github.com/nanobus/go-functions/transports/mux"

	"github.com/nanobus/nanobus/compute"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
)

var defaultInvokeURL string

func init() {
	defaultInvokeURL = os.Getenv("INVOKE_BASE_URL")
	if defaultInvokeURL == "" {
		defaultInvokeURL = "http://localhost:8000"
	}
}

type MuxConfig struct {
	BaseURL string `mapstructure:"baseUrl"`
}

// Mux is the NamedLoader for the mux compute.
func Mux() (string, compute.Loader) {
	return "mux", MuxLoader
}

func MuxLoader(with interface{}, resolver resolve.ResolveAs) (*functions.Invoker, error) {
	c := MuxConfig{
		BaseURL: defaultInvokeURL,
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	msgpackcodec := msgpack_codec.New()
	m := transport_mux.New(c.BaseURL, msgpackcodec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, msgpackcodec)

	return invoker, nil
}
