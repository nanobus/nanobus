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

const defaultInvokeURL = "http://127.0.0.1:9000"

type MuxConfig struct {
	BaseURL string `mapstructure:"baseUrl"`
}

// Mux is the NamedLoader for the mux compute.
func Mux() (string, compute.Loader) {
	return "mux", MuxLoader
}

func MuxLoader(with interface{}, resolver resolve.ResolveAs) (*functions.Invoker, error) {
	baseURL := os.Getenv("APP_URL")
	if baseURL == "" {
		baseURL = defaultInvokeURL
	}
	c := MuxConfig{
		BaseURL: baseURL,
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	msgpackcodec := msgpack_codec.New()
	m := transport_mux.New(c.BaseURL, msgpackcodec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, msgpackcodec)

	return invoker, nil
}
