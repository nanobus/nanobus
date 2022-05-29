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

package main

import (
	"context"
	"encoding/json"
	"log"

	functions "github.com/nanobus/nanobus/channel"
	jsoncodec "github.com/nanobus/nanobus/channel/codecs/json"
	"github.com/nanobus/nanobus/channel/transports/mux"
)

func main() {
	ctx := context.Background()
	codec := jsoncodec.New()
	m := mux.New("http://localhost:8081", codec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, nil, codec)

	var response interface{}
	err := invoker.InvokeWithReturn(ctx, functions.Receiver{
		Namespace: "customers.v1.Inbound",
		Operation: "createCustomer",
	}, map[string]interface{}{
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
