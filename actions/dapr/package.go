package dapr

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/nanobus/nanobus/actions"
)

var daprBaseURI string

func init() {
	daprPort := os.Getenv("DAPR_HTTP_PORT")
	if daprPort == "" {
		daprPort = "3500"
	}
	daprBaseURI = fmt.Sprintf("http://localhost:%s", daprPort)
}

var All = []actions.NamedLoader{
	InvokeActor,
	InvokeBinding,
	SetState,
	GetState,
	PublishMessage,
	SQLExec,
}

// Common dependencies

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func GET(ctx context.Context, httpClient HTTPClient, url string, decode func([]byte) error) error {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		url,
		nil)
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("expected 2XX status code; received %d", resp.StatusCode)
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBytes) > 0 {
		if err = decode(responseBytes); err != nil {
			return err
		}
	}

	return nil
}

func POST(ctx context.Context, httpClient HTTPClient, url string, encode func() ([]byte, error), decode func([]byte) error) error {
	requestBytes, err := encode()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		url,
		bytes.NewReader(requestBytes))
	if err != nil {
		return err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("expected 2XX status code; received %d", resp.StatusCode)
	}

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(responseBytes) > 0 {
		if err = decode(responseBytes); err != nil {
			return err
		}
	}

	return nil
}
