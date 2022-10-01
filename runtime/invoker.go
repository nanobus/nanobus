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

package runtime

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"reflect"

	"github.com/WasmRS/wasmrs-go/operations"
	"github.com/WasmRS/wasmrs-go/payload"
	"github.com/WasmRS/wasmrs-go/rx/flux"
	"github.com/WasmRS/wasmrs-go/rx/mono"
	"github.com/go-logr/logr"
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/channel"
	"github.com/nanobus/nanobus/compute"
)

type Invoker struct {
	log logr.Logger
	compute.Invoker
	codec     channel.Codec
	ns        Namespaces
	ops       operations.Table
	runnables []Runnable
	targets   []Target
}

type Target struct {
	Namespace string
	Operation string
}

func NewInvoker(log logr.Logger, ns Namespaces, codec channel.Codec) *Invoker {
	ops := make(operations.Table, 0, 10)
	runnables := make([]Runnable, 0, 10)
	targets := make([]Target, 0, 10)

	index := uint32(0)
	for namespace, functions := range ns {
		for operation, r := range functions {
			targets = append(targets, Target{
				Namespace: namespace,
				Operation: operation,
			})
			ops = append(ops, operations.Operation{
				Index:     index,
				Type:      operations.RequestResponse,
				Direction: operations.Export,
				Namespace: namespace,
				Operation: operation,
			})
			ops = append(ops, operations.Operation{
				Index:     index,
				Type:      operations.FireAndForget,
				Direction: operations.Export,
				Namespace: namespace,
				Operation: operation,
			})
			ops = append(ops, operations.Operation{
				Index:     index,
				Type:      operations.RequestStream,
				Direction: operations.Export,
				Namespace: namespace,
				Operation: operation,
			})
			ops = append(ops, operations.Operation{
				Index:     index,
				Type:      operations.RequestChannel,
				Direction: operations.Export,
				Namespace: namespace,
				Operation: operation,
			})
			runnables = append(runnables, r)
			index++
		}
	}

	return &Invoker{
		log:       log,
		codec:     codec,
		ns:        ns,
		ops:       ops,
		runnables: runnables,
		targets:   targets,
	}
}

func (i *Invoker) Close() error { return nil }

func (i *Invoker) Operations() operations.Table {
	return i.ops
}

func (i *Invoker) FireAndForget(ctx context.Context, p payload.Payload) {
	r, data := i.lookup(p)
	go r(ctx, data)
}

func (i *Invoker) RequestResponse(ctx context.Context, p payload.Payload) mono.Mono[payload.Payload] {
	r, data := i.lookup(p)
	return mono.Create(func(sink mono.Sink[payload.Payload]) {
		go func() {
			result, err := r(ctx, data)
			if err != nil {
				sink.Error(err)
				return
			}

			if isNil(result) {
				sink.Success(payload.New(nil))
				return
			}

			data, err := i.codec.Encode(result)
			if err != nil {
				sink.Error(err)
				return
			}

			sink.Success(payload.New(data))
		}()
	})
}

func (i *Invoker) RequestStream(ctx context.Context, p payload.Payload) flux.Flux[payload.Payload] {
	r, data := i.lookup(p)
	return flux.Create(func(sink flux.Sink[payload.Payload]) {
		go func() {
			// TODO: set sink in context
			result, err := r(ctx, data)
			if err != nil {
				sink.Error(err)
				return
			}

			if isNil(result) {
				sink.Next(payload.New(nil))
				return
			}

			sink.Complete()
		}()
	})
}

func (*Invoker) RequestChannel(context.Context, payload.Payload, flux.Flux[payload.Payload]) flux.Flux[payload.Payload] {
	return nil
}

func (i *Invoker) lookup(p payload.Payload) (Runnable, actions.Data) {
	md := p.Metadata()
	index := binary.BigEndian.Uint32(md)
	r := i.runnables[index]
	t := i.targets[index]
	var input interface{}
	i.codec.Decode(p.Data(), &input)
	data := actions.Data{
		"input": input,
	}

	if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
		logOutbound(i.log, t.Namespace+"/"+t.Operation, string(jsonBytes))
	}

	return r, data
}

func isNil(val interface{}) bool {
	return val == nil ||
		(reflect.ValueOf(val).Kind() == reflect.Ptr &&
			reflect.ValueOf(val).IsNil())
}

func logOutbound(log logr.Logger, target string, data string) {
	l := log //.V(10)
	if l.Enabled() {
		l.Info("<== " + target + " " + data)
	}
}
