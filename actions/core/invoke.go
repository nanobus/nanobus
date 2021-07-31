package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/expr"
	"github.com/nanobus/nanobus/resolve"
)

var defaultInvokeURL string

func init() {
	defaultInvokeURL = os.Getenv("INVOKE_BASE_URL")
	if defaultInvokeURL == "" {
		defaultInvokeURL = "http://localhost:8000"
	}
}

type InvokeConfig struct {
	BaseURI string `mapstructure:"baseUri"`
	// Function is name of function to invoke.
	Function string `mapstructure:"function"`
	// Input optionally transforms the input sent to the function.
	Input *expr.DataExpr `mapstructure:"input"`
}

// Invoke is the NamedLoader for the invoke action.
func Invoke() (string, actions.Loader) {
	return "invoke", InvokeLoader
}

func InvokeLoader(with interface{}, resolver resolve.ResolveAs) (actions.Action, error) {
	c := InvokeConfig{
		BaseURI: defaultInvokeURL,
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var httpClient HTTPClient
	if err := resolve.Resolve(resolver,
		"client:http", &httpClient); err != nil {
		return nil, err
	}

	return InvokeAction(httpClient, &c), nil
}

func InvokeAction(
	httpClient HTTPClient,
	config *InvokeConfig) actions.Action {
	return func(ctx context.Context, data actions.Data) (interface{}, error) {
		input := data["input"]
		if config.Input != nil {
			var err error
			input, err = config.Input.Eval(data)
			if err != nil {
				return nil, err
			}
		}

		if inputBytes, ok := input.([]byte); ok {
			if err := json.Unmarshal(inputBytes, &input); err != nil {
				return nil, err
			}
		}
		if inputString, ok := input.(string); ok {
			if err := json.Unmarshal([]byte(inputString), &input); err != nil {
				return nil, err
			}
		}

		requestBytes, err := msgpack.Marshal(&input)
		if err != nil {
			return nil, err
		}

		u, err := url.Parse(config.BaseURI)
		if err != nil {
			return nil, err
		}
		u.Path = path.Join(u.Path, config.Function)

		req, err := http.NewRequestWithContext(
			ctx,
			"POST",
			u.String(),
			bytes.NewReader(requestBytes))
		if err != nil {
			return nil, err
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode/100 != 2 {
			return nil, fmt.Errorf("expected 2XX status code; received %d", resp.StatusCode)
		}

		responseBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		if len(responseBytes) > 0 {
			var response interface{}
			if err = msgpack.Unmarshal(responseBytes, &response); err != nil {
				return nil, err
			}

			return response, nil
		}

		return nil, nil
	}
}
