/*
 * Copyright 2022 The NanoBus Authors.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package time

import (
	"context"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"

	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/handler"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/transport"
)

type Scheduler struct {
	id          string
	ctx         context.Context
	log         logr.Logger
	tracer      trace.Tracer
	schedule    string
	daemon      *gocron.Scheduler
	lastruntime time.Time
	numruns     int
	invoker     transport.Invoker
	handler     handler.Handler
}

func TimeSchedulerV1Loader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (transport.Transport, error) {
	var log logr.Logger
	var tracer trace.Tracer
	var transportInvoker transport.Invoker
	if err := resolve.Resolve(resolver,
		"transport:invoker", &transportInvoker,
		"system:logger", &log,
		"system:tracer", &tracer,
	); err != nil {
		return nil, err
	}

	// Defaults
	c := TimeSchedulerV1Config{}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	return NewScheduler(ctx, log, tracer, transportInvoker, c)
}

func NewScheduler(ctx context.Context, log logr.Logger, tracer trace.Tracer, transportInvoker transport.Invoker, config TimeSchedulerV1Config) (*Scheduler, error) {
	return &Scheduler{
		id:          uuid.New().String(),
		ctx:         ctx,
		log:         log,
		tracer:      tracer,
		daemon:      nil,
		schedule:    config.Schedule,
		lastruntime: time.Time{},
		numruns:     0,
		invoker:     transportInvoker,
		handler:     config.Handler,
	}, nil
}

func (t *Scheduler) Listen() error {
	input := map[string]interface{}{}
	s := gocron.NewScheduler(time.UTC)
	s.Cron(t.schedule).Do(func() {
		_, err := t.invoker(t.ctx, t.handler.Interface, t.id, t.handler.Operation, input, transport.BypassAuthorization)
		if err != nil {
			t.log.Error(err, "Error in %q", t.handler)
		}
	})

	t.daemon = s

	t.log.Info("Schedule Deamon Started", "schedule", t.schedule)
	s.StartBlocking()

	return nil
}

func (t *Scheduler) Close() (err error) {
	if t.daemon != nil {
		t.daemon.Stop()
		t.daemon = nil
	}

	return nil
}
