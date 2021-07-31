package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nanobus/go-functions"
	jsoncodec "github.com/nanobus/go-functions/codecs/json"
	"github.com/nanobus/go-functions/transports/mux"
)

func main() {
	ctx := context.Background()
	codec := jsoncodec.New()
	m := mux.New("http://localhost:8081", codec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, codec)

	var response interface{}
	err := invoker.InvokeWithReturn(ctx, "customers.v1.Inbound", "createCustomer", map[string]interface{}{
		"id":        1234,
		"firstName": "John",
		"lastName":  "Doe",
		"email":     "john.doe@gmail.com",
	}, &response)
	if err != nil {
		log.Fatal(err)
	}

	if jsonBytes, err := json.MarshalIndent(response, "", "  "); err == nil {
		log.Printf("RESPONSE: %s", string(jsonBytes))
	}
}
