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
	"encoding/json"
	"flag"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 4 {
		log.Fatal("usage: nats-req <namespace> <service> <operation> <data>")
	}

	ns := args[0]
	service := args[1]
	operation := args[2]
	data := args[3]

	// Connect to a server
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	subject := ns + "." + service + "." + operation
	log.Println(subject)
	request := nats.NewMsg(subject)
	request.Header.Set("Namespace", ns)
	request.Header.Set("Service", service)
	request.Header.Set("Function", operation)

	var payload interface{}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		log.Fatal(err)
	}

	request.Data, err = msgpack.Marshal(&payload)
	if err != nil {
		log.Fatal(err)
	}
	request.Header.Set("Content-Type", "application/msgpack")
	reply, err := nc.RequestMsg(request, 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range reply.Header {
		log.Println(k, "=", v)
	}

	payload = nil
	jsonBytes := reply.Data
	if err = msgpack.Unmarshal(reply.Data, &payload); err == nil {
		jsonBytes, err = json.Marshal(&payload)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println(string(jsonBytes))
}
