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

	"github.com/go-logr/zapr"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nanobus/adapter-go/codec/json"
	"github.com/nanobus/adapter-go/stateful"

	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

func main() {
	ctx := context.Background()
	// Initialize logger
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	log := zapr.NewLogger(zapLog)

	codec := json.New()
	cache, err := stateful.NewLRUCache(200)
	if err != nil {
		panic(err)
	}
	app := customers.NewApp(ctx, codec, cache)
	outbound := app.NewOutbound()
	service := customers.NewService(log, outbound)
	customerActor := customers.NewCustomerActor()

	app.RegisterInbound(service)
	app.RegisterCustomerActor(customerActor)

	if err := app.Start(); err != nil {
		log.Error(err, "Exit with error")
	}
}
