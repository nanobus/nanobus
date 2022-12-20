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

	"github.com/go-logr/logr"
	"go.opentelemetry.io/otel/trace"

	"github.com/nanobus/nanobus/pkg/config"
	"github.com/nanobus/nanobus/pkg/resolve"
	"github.com/nanobus/nanobus/pkg/runtime"
	"github.com/nanobus/nanobus/pkg/transport"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	log         logr.Logger
	tracer      trace.Tracer
	schedule    string
	daemon      *gocron.Scheduler
	lastruntime time.Time
	numruns     int
	action      []runtime.Component
}

func TimeSchedulerV1Loader(ctx context.Context, with interface{}, resolver resolve.ResolveAs) (transport.Transport, error) {
	var log logr.Logger
	var tracer trace.Tracer
	if err := resolve.Resolve(resolver,
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

	return NewScheduler(log, tracer, c)
}

func NewScheduler(log logr.Logger, tracer trace.Tracer, config TimeSchedulerV1Config) (*Scheduler, error) {
	return &Scheduler{
		log:         log,
		tracer:      tracer,
		daemon:      nil,
		schedule:    config.Schedule,
		lastruntime: time.Time{},
		numruns:     0,
		action:      config.Action,
	}, nil
}

func (t *Scheduler) Listen() error {
	s := gocron.NewScheduler(time.UTC)
	s.Cron("*/1 * * * *").Do(true)
	s.StartAsync()

	t.daemon = s
	t.log.Info("Schedule Deamon Started", "schedule", t.schedule)

	return nil
}

func (t *Scheduler) Close() (err error) {
	if t.daemon != nil {
		t.daemon.Stop()
		t.daemon = nil
	}

	return nil
}
