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

package mesh

import (
	"context"
	"encoding/binary"
	"sync/atomic"

	"github.com/WasmRS/wasmrs-go/operations"
	"github.com/WasmRS/wasmrs-go/payload"
	"github.com/WasmRS/wasmrs-go/rx/flux"
	"github.com/WasmRS/wasmrs-go/rx/mono"
	"github.com/nanobus/nanobus/compute"
	"go.uber.org/multierr"
)

type (
	Mesh struct {
		instances   map[string]compute.Invoker
		exports     map[string]map[string]*atomic.Pointer[destination]
		unsatisfied []*pending
		done        chan struct{}
	}

	destination struct {
		instance compute.Invoker
		index    uint32
	}

	pending struct {
		instance compute.Invoker
		oper     operations.Operation
	}
)

func New() *Mesh {
	return &Mesh{
		instances:   make(map[string]compute.Invoker),
		exports:     map[string]map[string]*atomic.Pointer[destination]{},
		unsatisfied: make([]*pending, 0, 10),
		done:        make(chan struct{}),
	}
}

func (m *Mesh) RequestResponse(ctx context.Context, namespace, operation string, p payload.Payload) mono.Mono[payload.Payload] {
	ns, ok := m.exports[namespace]
	if !ok {
		return nil
	}

	ptr, ok := ns[operation]
	if !ok {
		return nil
	}
	dest := ptr.Load()

	return dest.RequestResponse(ctx, p)
}

func (m *Mesh) FireAndForget(ctx context.Context, namespace, operation string, p payload.Payload) {
	ns, ok := m.exports[namespace]
	if !ok {
		return
	}

	ptr, ok := ns[operation]
	if !ok {
		return
	}
	dest := ptr.Load()

	dest.FireAndForget(ctx, p)
}

func (m *Mesh) RequestStream(ctx context.Context, namespace, operation string, p payload.Payload) flux.Flux[payload.Payload] {
	ns, ok := m.exports[namespace]
	if !ok {
		return nil
	}

	ptr, ok := ns[operation]
	if !ok {
		return nil
	}
	dest := ptr.Load()

	return dest.RequestStream(ctx, p)
}

func (m *Mesh) RequestChannel(ctx context.Context, namespace, operation string, p payload.Payload, in flux.Flux[payload.Payload]) flux.Flux[payload.Payload] {
	ns, ok := m.exports[namespace]
	if !ok {
		return nil
	}

	ptr, ok := ns[operation]
	if !ok {
		return nil
	}
	dest := ptr.Load()

	return dest.RequestChannel(ctx, p, in)
}

func (m *Mesh) Close() error {
	var merr error
	for _, inst := range m.instances {
		if err := inst.Close(); err != nil {
			merr = multierr.Append(merr, err)
		}
	}
	close(m.done)
	return merr
}

func (m *Mesh) WaitUntilShutdown() error {
	<-m.done
	return nil
}

func (m *Mesh) Link(inst compute.Invoker) {
	opers := inst.Operations()

	numExported := 0
	for _, op := range opers {
		switch op.Direction {
		case operations.Export:
			ns, ok := m.exports[op.Namespace]
			if !ok {
				ns = make(map[string]*atomic.Pointer[destination])
				m.exports[op.Namespace] = ns
			}
			ptr, ok := ns[op.Operation]
			if !ok {
				ptr = &atomic.Pointer[destination]{}
				ns[op.Operation] = ptr
			}

			ptr.Store(&destination{
				instance: inst,
				index:    op.Index,
			})
			numExported++

		case operations.Import:
			if ok := m.linkOperation(inst, op); !ok {
				m.unsatisfied = append(m.unsatisfied, &pending{
					instance: inst,
					oper:     op,
				})
			}
		}
	}

	if numExported > 0 && len(m.unsatisfied) > 0 {
		filtered := m.unsatisfied[:0]
		for _, u := range m.unsatisfied {
			if ok := m.linkOperation(u.instance, u.oper); !ok {
				filtered = append(filtered, u)
			}
		}
		m.unsatisfied = filtered
	}
}

func (m *Mesh) linkOperation(inst compute.Invoker, op operations.Operation) bool {
	ns, ok := m.exports[op.Namespace]
	if !ok {
		return false
	}

	ptr, ok := ns[op.Operation]
	if !ok {
		return false
	}
	dest := ptr.Load()

	switch op.Type {
	case operations.RequestResponse:
		inst.SetRequestResponseHandler(op.Index, dest.RequestResponse)
	case operations.FireAndForget:
		inst.SetFireAndForgetHandler(op.Index, dest.FireAndForget)
	case operations.RequestStream:
		inst.SetRequestStreamHandler(op.Index, dest.RequestStream)
	case operations.RequestChannel:
		inst.SetRequestChannelHandler(op.Index, dest.RequestChannel)
	}

	return true
}

func (d *destination) RequestResponse(ctx context.Context, p payload.Payload) mono.Mono[payload.Payload] {
	md := p.Metadata()
	if md != nil {
		binary.BigEndian.PutUint32(md, d.index)
	}
	return d.instance.RequestResponse(ctx, p)
}

func (d *destination) FireAndForget(ctx context.Context, p payload.Payload) {
	md := p.Metadata()
	if md != nil {
		binary.BigEndian.PutUint32(md, d.index)
	}
	d.instance.FireAndForget(ctx, p)
}

func (d *destination) RequestStream(ctx context.Context, p payload.Payload) flux.Flux[payload.Payload] {
	md := p.Metadata()
	if md != nil {
		binary.BigEndian.PutUint32(md, d.index)
	}
	return d.instance.RequestStream(ctx, p)
}

func (d *destination) RequestChannel(ctx context.Context, p payload.Payload, in flux.Flux[payload.Payload]) flux.Flux[payload.Payload] {
	md := p.Metadata()
	if md != nil {
		binary.BigEndian.PutUint32(md, d.index)
	}
	return d.instance.RequestChannel(ctx, p, in)
}
