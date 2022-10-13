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
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/oklog/run"
	"github.com/vmihailenco/msgpack/v5"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otel_resource "go.opentelemetry.io/otel/sdk/resource"
	sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"

	// COMPONENT MODEL / PLUGGABLE COMPONENTS
	"github.com/WasmRS/wasmrs-go/payload"
	proto "github.com/dapr/dapr/pkg/proto/components/v1"

	// NANOBUS CORE
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/function"
	"github.com/nanobus/nanobus/mesh"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/resource"
	"github.com/nanobus/nanobus/runtime"
	"github.com/nanobus/nanobus/security/claims"
	"github.com/nanobus/nanobus/version"

	// CHANNELS
	json_codec "github.com/nanobus/nanobus/channel/codecs/json"
	msgpack_codec "github.com/nanobus/nanobus/channel/codecs/msgpack"

	// SPECIFICATION LANGUAGES
	"github.com/nanobus/nanobus/spec"
	spec_apex "github.com/nanobus/nanobus/spec/apex"

	// COMPONENTS
	"github.com/nanobus/nanobus/compute"
	compute_wasmrs "github.com/nanobus/nanobus/compute/wasmrs"

	// ACTIONS
	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/actions/core"
	"github.com/nanobus/nanobus/actions/dapr"
	"github.com/nanobus/nanobus/actions/gorm"
	"github.com/nanobus/nanobus/actions/postgres"

	// CODECS
	"github.com/nanobus/nanobus/codec"
	cloudevents_avro "github.com/nanobus/nanobus/codec/cloudevents/avro"
	cloudevents_json "github.com/nanobus/nanobus/codec/cloudevents/json"
	"github.com/nanobus/nanobus/codec/confluentavro"
	codec_json "github.com/nanobus/nanobus/codec/json"
	codec_msgpack "github.com/nanobus/nanobus/codec/msgpack"
	codec_text "github.com/nanobus/nanobus/codec/text"

	// TELEMETRY / TRACING
	otel_tracing "github.com/nanobus/nanobus/telemetry/tracing"
	tracing_jaeger "github.com/nanobus/nanobus/telemetry/tracing/jaeger"
	tracing_otlp "github.com/nanobus/nanobus/telemetry/tracing/otlp"
	tracing_stdout "github.com/nanobus/nanobus/telemetry/tracing/stdout"

	// TRANSPORTS
	"github.com/nanobus/nanobus/transport"
	"github.com/nanobus/nanobus/transport/filter"
	"github.com/nanobus/nanobus/transport/filter/jwt"
	"github.com/nanobus/nanobus/transport/httprpc"
	"github.com/nanobus/nanobus/transport/nats"
	"github.com/nanobus/nanobus/transport/rest"
)

type Runtime struct {
	log        logr.Logger
	config     *runtime.Configuration
	namespaces spec.Namespaces
	processor  *runtime.Processor
	resolver   resolve.DependencyResolver
	resolveAs  resolve.ResolveAs
	env        runtime.Environment
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zapLog := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig),
		zapcore.AddSync(colorable.NewColorableStdout()),
		zapcore.DebugLevel,
	))
	//zapLog, err := zapConfig.Build()
	//zapLog, err := zap.NewProduction()
	// if err != nil {
	// 	panic(err)
	// }
	// zapLog := zap.NewExample()
	log := zapr.NewLogger(zapLog)

	// NanoBus flags

	var busListenAddr string
	flag.StringVar(
		&busListenAddr,
		"bus-listen-addr",
		LookupEnvOrString("BUS_LISTEN_ADDR", "localhost:32320"),
		"bus listen address",
	)
	var busFile string
	flag.StringVar(
		&busFile,
		"bus",
		LookupEnvOrString("CONFIG", "bus.yaml"),
		"The bus configuration file",
	)
	flag.Parse()
	args := flag.Args()

	// Load the configuration
	config, err := loadConfiguration(busFile)
	if err != nil {
		log.Error(err, "could not load configuration", "file", busFile)
		os.Exit(1)
	}

	// Transport registration
	transportRegistry := transport.Registry{}
	transportRegistry.Register(
		httprpc.Load,
		rest.Load,
		nats.Load,
	)

	// Spec registration
	specRegistry := spec.Registry{}
	specRegistry.Register(
		spec_apex.Apex,
	)

	// Filter registration
	filterRegistry := filter.Registry{}
	filterRegistry.Register(
		jwt.JWT,
	)

	// Compute registration
	computeRegistry := compute.Registry{}
	computeRegistry.Register(
		compute_wasmrs.WasmRS,
		// compute_mux.Mux,
		// compute_wapc.WaPC,
		// compute_rsocket.RSocket,
	)

	// Codec registration
	codecRegistry := codec.Registry{}
	codecRegistry.Register(
		codec_json.JSON,
		codec_msgpack.MsgPack,
		confluentavro.ConfluentAvro,
		cloudevents_avro.CloudEventsAvro,
		cloudevents_json.CloudEventsJSON,
		codec_text.Plain,
		codec_text.HTML,
	)

	resourceRegistry := resource.Registry{}
	resourceRegistry.Register(
		postgres.Connection,
		gorm.Connection,
		dapr.PubSub,
		dapr.StateStore,
		dapr.OutputBinding,
	)

	tracingRegistry := otel_tracing.Registry{}
	tracingRegistry.Register(
		tracing_jaeger.Jaeger,
		tracing_otlp.OTLP,
		tracing_stdout.Stdout,
	)

	// Action registration
	actionRegistry := actions.Registry{}
	actionRegistry.Register(core.All...)
	actionRegistry.Register(postgres.All...)
	actionRegistry.Register(gorm.All...)
	actionRegistry.Register(dapr.All...)

	// Codecs
	jsoncodec := json_codec.New()
	msgpackcodec := msgpack_codec.New()

	// Dependencies
	// var invoker *channel.Invoker
	// var busInvoker compute.BusInvoker
	httpClient := getHTTPClient()
	env := getEnvironment()
	namespaces := make(spec.Namespaces)
	dependencies := map[string]interface{}{
		"system:logger": log,
		//"system:tracer":   tracer,
		"client:http":     httpClient,
		"codec:json":      jsoncodec,
		"codec:msgpack":   msgpackcodec,
		"spec:namespaces": namespaces,
		"os:env":          env,
	}
	resolver := func(name string) (interface{}, bool) {
		dep, ok := dependencies[name]
		return dep, ok
	}
	resolveAs := resolve.ToResolveAs(resolver)

	var spanExporter sdk_trace.SpanExporter
	if config.Tracing != nil {
		loadable, ok := tracingRegistry[config.Tracing.Uses]
		if !ok {
			log.Error(nil, "Could not find codec", "type", config.Tracing.Uses)
			os.Exit(1)
		}
		var err error
		spanExporter, err = loadable(ctx, config.Tracing.With, resolveAs)
		if err != nil {
			log.Error(err, "Error loading codec", "type", config.Tracing.Uses)
			os.Exit(1)
		}
	}

	var tp trace.TracerProvider

	if spanExporter == nil {
		tp = trace.NewNoopTracerProvider()
	} else {
		ntp := sdk_trace.NewTracerProvider(
			sdk_trace.WithBatcher(spanExporter),
			sdk_trace.WithResource(newOtelResource(config.Application)),
		)
		defer func() {
			if err := ntp.Shutdown(ctx); err != nil {
				log.Error(err, "error shutting down trace provider")
				os.Exit(1)
			}
		}()
		tp = ntp
	}

	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("NanoBus")
	dependencies["system:tracer"] = tracer

	if len(config.Specs) == 0 {
		config.Specs = append(config.Specs, runtime.Component{
			Uses: "apex",
			With: map[string]interface{}{
				"filename": "spec.apexlang",
			},
		})
	}
	for _, spec := range config.Specs {
		loader, ok := specRegistry[spec.Uses]
		if !ok {
			log.Error(nil, "could not find spec", "type", spec.Uses)
			os.Exit(1)
		}
		nss, err := loader(ctx, spec.With, resolveAs)
		if err != nil {
			log.Error(err, "error loading spec", "type", spec.Uses)
			os.Exit(1)
		}
		for _, ns := range nss {
			namespaces[ns.Name] = ns
		}
	}

	if config.Codecs == nil {
		config.Codecs = map[string]runtime.Component{}
	}
	for name, loadable := range codecRegistry {
		if loadable.Auto {
			if _, exists := config.Codecs[name]; !exists {
				config.Codecs[name] = runtime.Component{
					Uses: name,
				}
			}
		}
	}

	codecs := make(codec.Codecs)
	codecsByContentType := make(codec.Codecs)
	for name, component := range config.Codecs {
		loadable, ok := codecRegistry[component.Uses]
		if !ok {
			log.Error(nil, "Could not find codec", "type", component.Uses)
			os.Exit(1)
		}
		c, err := loadable.Loader(component.With, resolveAs)
		if err != nil {
			log.Error(err, "Error loading codec", "type", component.Uses)
			os.Exit(1)
		}
		codecs[name] = c
		codecsByContentType[c.ContentType()] = c
	}
	dependencies["codec:lookup"] = codecs
	dependencies["codec:byContentType"] = codecsByContentType

	resources := resource.Resources{}
	for name, component := range config.Resources {
		log.Info("Initializing resource", "name", name)

		loader, ok := resourceRegistry[component.Uses]
		if !ok {
			log.Error(nil, "Could not find resource", "type", component.Uses)
			os.Exit(1)
		}
		c, err := loader(ctx, component.With, resolveAs)
		if err != nil {
			log.Error(err, "Error loading resource", "type", component.Uses)
			os.Exit(1)
		}

		resources[name] = c
	}
	dependencies["resource:lookup"] = resources

	// Create processor
	processor, err := runtime.NewProcessor(ctx, log, tracer, config, actionRegistry, resolver)
	if err != nil {
		log.Error(err, "Could not create NanoBus runtime")
		os.Exit(1)
	}

	rt := Runtime{
		log:        log,
		config:     config,
		namespaces: namespaces,
		processor:  processor,
		resolver:   resolver,
		resolveAs:  resolveAs,
		env:        env,
	}
	// busInvoker = rt.BusInvoker
	// dependencies["bus:invoker"] = busInvoker
	dependencies["state:invoker"] = func(ctx context.Context, namespace, id, key string) ([]byte, error) {
		// TODO: Retrieve state
		return []byte{}, nil
	}

	if err = processor.Initialize(); err != nil {
		log.Error(err, "Could not initialize processor")
		os.Exit(1)
	}

	m := mesh.New(tracer)

	m.Link(runtime.NewInvoker(log, processor.GetProviders(), msgpackcodec))

	for _, comp := range config.Compute {
		computeLoader, ok := computeRegistry[comp.Uses]
		if !ok {
			log.Error(err, "could not find compute", "type", comp.Uses)
			os.Exit(1)
		}
		invoker, err := computeLoader(ctx, comp.With, resolveAs)
		if err != nil {
			log.Error(err, "could not load compute", "type", comp.Uses)
			os.Exit(1)
		}
		m.Link(invoker)
	}
	dependencies["compute:mesh"] = m

	for _, subscription := range config.Subscriptions {
		pubsub, err := resource.Get[proto.PubSubClient](resources, subscription.Resource)
		if err != nil {
			log.Error(err, "Could not load resource", "name", subscription.Resource)
			os.Exit(1)
		}

		c, ok := codecs[subscription.Codec]
		if !ok {
			log.Error(nil, "Could not find codec", "name", subscription.Resource)
			os.Exit(1)
		}

		pull, err := pubsub.PullMessages(ctx)
		if err != nil {
			log.Error(nil, "Could not pull messages", "name", subscription.Resource)
			os.Exit(1)
		}

		go func(pull proto.PubSub_PullMessagesClient, c codec.Codec, sub runtime.Subscription) {
			if err := pull.Send(&proto.PullMessagesRequest{
				Topic: &proto.Topic{
					Name:     sub.Topic,
					Metadata: sub.Metadata,
				},
			}); err != nil {
				log.Error(err, "Error subscribing")
				return
			}

			log.Info("Subscribed to pubsub", "resource", sub.Resource, "topic", sub.Topic)

			for {
				recv, err := pull.Recv()
				if err == io.EOF || err == context.Canceled {
					return
				}
				if err != nil {
					log.Error(err, "Error receiving messages")
					return
				}

				input, messageType, err := c.Decode(recv.Data, sub.CodecArgs...)
				if err != nil {
					log.Error(err, "could not decode message")
					pull.Send(&proto.PullMessagesRequest{
						AckMessageId: recv.Id,
						AckError: &proto.AckMessageError{
							Message: err.Error(),
						},
					})
					continue
				}

				data := actions.Data{
					"input": input,
				}

				pipelineName := sub.Function
				if pipelineName == "" {
					pipelineName = messageType
				}

				if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
					logInbound(log, "events::type="+pipelineName, string(jsonBytes))
				}

				_, err = processor.Event(ctx, pipelineName, data)
				if err != nil {
					log.Error(err, "could not process message")
					pull.Send(&proto.PullMessagesRequest{
						AckMessageId: recv.Id,
						AckError: &proto.AckMessageError{
							Message: err.Error(),
						},
					})
					continue
				}

				pull.Send(&proto.PullMessagesRequest{
					AckMessageId: recv.Id,
				})
			}
		}(pull, c, subscription)
	}

	// Big 'ol TODO
	// invoker = computeInstance.Invoker
	// dependencies["client:invoker"] = invoker

	filters := []filter.Filter{}
	if configFilters, ok := config.Filters["http"]; ok {
		for _, f := range configFilters {
			filterLoader, ok := filterRegistry[f.Uses]
			if !ok {
				log.Error(nil, "could not find filter", "type", f.Uses)
				os.Exit(1)
			}

			filter, err := filterLoader(ctx, f.With, resolveAs)
			if err != nil {
				log.Error(err, "could not load filter", "type", f.Uses)
				os.Exit(1)
			}

			filters = append(filters, filter)
		}
	}
	dependencies["filter:lookup"] = filters

	translateError := func(err error) *errorz.Error {
		if errz, ok := err.(*errorz.Error); ok {
			return errz
		}
		var te errorz.TemplateError

		if terrz, ok := err.(*errorz.TemplateError); ok && terrz != nil {
			te = *terrz
		} else {
			te = errorz.ParseTemplateError(err.Error())
		}

		tmpl, ok := config.Errors[te.Template]
		if !ok {
			// Default error if template matches a code name.
			if code, ok := errorz.CodeLookup[te.Template]; ok {
				return errorz.New(code)
			}

			return errorz.New(errorz.Internal, err.Error())
		}

		message := err.Error()
		if tmpl.Message != nil {
			message, _ = tmpl.Message.Eval(te.Metadata)
		}

		e := errorz.New(tmpl.Code, message)
		e.Type = te.Template
		if tmpl.Type != "" {
			e.Type = tmpl.Type
		}
		if tmpl.Status != 0 {
			e.Status = tmpl.Status
		}
		if tmpl.Title != nil {
			title, _ := tmpl.Title.Eval(te.Metadata)
			e.Title = title
		}
		if tmpl.Help != nil {
			instance, _ := tmpl.Help.Eval(te.Metadata)
			e.Help = instance
		}
		e.Metadata = te.Metadata

		return e
	}
	dependencies["errors:resolver"] = errorz.Resolver(translateError)

	// healthHandler := func(w http.ResponseWriter, r *http.Request) {
	// 	//fmt.Println("HEALTH HANDLER CALLED")
	// 	defer r.Body.Close()

	// 	w.Header().Set("Content-Type", "text/plain")
	// 	w.Write([]byte("OK"))
	// }

	transportInvoker := func(ctx context.Context, namespace, service, id, fn string, input interface{}) (interface{}, error) {
		if err := coalesceInput(namespaces, namespace, service, fn, input); err != nil {
			return nil, err
		}

		claimsMap := claims.FromContext(ctx)

		data := actions.Data{
			"claims": claimsMap,
			"input":  input,
		}

		ns := namespace + "." + service

		if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
			logInbound(rt.log, ns+"/"+fn, string(jsonBytes))
		}

		data["env"] = env

		ctx = function.ToContext(ctx, function.Function{
			Namespace: ns,
			Operation: fn,
		})

		response, ok, err := rt.processor.Service(ctx, namespace, service, fn, data)
		if err != nil {
			return nil, translateError(err)
		}

		// No pipeline exits for the operation so invoke directly.
		if !ok {
			payloadData, err := msgpack.Marshal(input)
			if err != nil {
				return nil, translateError(err)
			}

			metadata := make([]byte, 8)
			p := payload.New(payloadData, metadata)

			result, err := m.RequestResponse(ctx, ns, fn, p).Block()
			if err != nil {
				return nil, translateError(err)
			}

			var resultDecoded interface{}
			if err := msgpack.Unmarshal(result.Data(), &resultDecoded); err != nil {
				return nil, translateError(err)
			}

			response = resultDecoded
		}

		return response, err
	}
	dependencies["transport:invoker"] = transport.Invoker(transportInvoker)

	// go computeInstance.Start()

	var g run.Group
	if len(args) > 0 {
		log.Info("Executing process", "cmd", strings.Join(args, " "))
		command := args[0]
		args = args[1:]
		cmd := exec.CommandContext(ctx, command, args...)
		busPort := busListenAddr
		if i := strings.Index(busPort, ":"); i != -1 {
			busPort = busPort[i+1:]
		}
		g.Add(func() error {
			appEnv := []string{
				// fmt.Sprintf("PORT=%d", appPort),
				fmt.Sprintf("BUS_URL=http://127.0.0.1:%s", busPort),
			}
			// if computeInstance.Environ != nil {
			// 	appEnv = append(appEnv, computeInstance.Environ()...)
			// }
			env := []string{}
			env = append(env, os.Environ()...)
			env = append(env, appEnv...)
			cmd.Env = env
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}, func(error) {
			// TODO: send term sig instead
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		})
	}
	// {
	// 	// Expose the bus
	// 	r := mux.NewRouter()
	// 	r.HandleFunc("/healthz", healthHandler).Methods("GET")

	// 	log.Info("Bus listening", "address", busListenAddr)
	// 	ln, err := net.Listen("tcp", busListenAddr)
	// 	if err != nil {
	// 		log.Error(err, "could not create bus")
	// 		os.Exit(1)
	// 	}
	// 	g.Add(func() error {
	// 		return http.Serve(ln, r)
	// 	}, func(error) {
	// 		ln.Close()
	// 	})
	// }
	{
		g.Add(func() error {
			return m.WaitUntilShutdown()
		}, func(error) {
			m.Close()
		})
	}

	for name, comp := range config.Transports {
		name := name // Make copy
		loader, ok := transportRegistry[comp.Uses]
		if !ok {
			log.Error(nil, "unknown transport", "type", comp.Uses)
			os.Exit(1)
		}
		log.Info("Initializing transport", "name", name)
		t, err := loader(ctx, comp.With, resolveAs)
		if err != nil {
			log.Error(err, "could not load transport", "type", comp.Uses)
			os.Exit(1)
		}

		g.Add(func() error {
			return t.Listen()
		}, func(error) {
			t.Close()
		})
	}

	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}

	err = g.Run()
	if err != nil {
		if _, isSignal := err.(run.SignalError); !isSignal {
			log.Error(err, "unexpected error")
		}
	}
}

// func (rt *Runtime) BusInvoker(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error) {
// 	if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
// 		logOutbound(rt.log, namespace+"."+service+"/"+function, string(jsonBytes))
// 	}

// 	data := actions.Data{
// 		"input": input,
// 		"env":   rt.env,
// 	}

// 	output, err := rt.processor.Provider(ctx, namespace, service, function, data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if output, err = coalesceOutput(rt.namespaces, namespace, service, function, output); err != nil {
// 		return nil, err
// 	}

// 	return output, nil
// }

func loadConfiguration(filename string) (*runtime.Configuration, error) {
	// TODO: Load from file or URI
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c, err := runtime.LoadYAML(f)
	if err != nil {
		return nil, err
	}

	for _, imp := range c.Import {
		imported, err := loadConfiguration(imp)
		if err != nil {
			return nil, err
		}
		runtime.Combine(c, imported)
	}

	return c, nil
}

func getHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}

func getEnvironment() runtime.Environment {
	return environmentToMap(os.Environ(), func(item string) (key, val string) {
		splits := strings.SplitN(item, "=", 1)
		key = splits[0]
		if len(splits) > 1 {
			val = splits[1]
		}

		return
	})
}

func environmentToMap(environment []string, getkeyval func(item string) (key, val string)) map[string]string {
	items := make(map[string]string)
	for _, item := range environment {
		key, val := getkeyval(item)
		items[key] = val
	}

	return items
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return defaultVal
}

func LookupEnvOrInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err != nil {
			return intVal
		}
	}

	return defaultVal
}

func coalesceInput(namespaces spec.Namespaces, namespace, service, function string, input interface{}) error {
	if oper, ok := namespaces.Operation(namespace, service, function); ok {
		if oper.Parameters != nil {
			inputMap, ok := coalesce.ToMapSI(input, true)
			if !ok {
				return fmt.Errorf("%w: input is not a map", transport.ErrBadInput)
			}
			input = inputMap
			if err := oper.Parameters.Coalesce(inputMap, true); err != nil {
				var errz *errorz.Error
				if errors.As(err, &errz) {
					return errz
				}
				return fmt.Errorf("%w: %v", transport.ErrBadInput, err)
			}
		}
	} else {
		coalesce.Integers(input)
	}
	return nil
}

func coalesceOutput(namespaces spec.Namespaces, namespace, service, function string, output interface{}) (interface{}, error) {
	var err error
	if oper, ok := namespaces.Operation(namespace, service, function); ok {
		if oper.Returns != nil && output != nil {
			output, _, err = oper.Returns.Coalesce(output, false)
			if err != nil {
				return nil, err
			}
		} else {
			coalesce.Integers(output)
		}
	} else {
		coalesce.Integers(output)
	}
	return output, err
}

func logInbound(log logr.Logger, target string, data string) {
	l := log //.V(10)
	if l.Enabled() {
		l.Info("==> " + target + " " + data)
	}
}

func logOutbound(log logr.Logger, target string, data string) {
	l := log //.V(10)
	if l.Enabled() {
		l.Info("<== " + target + " " + data)
	}
} // )

// newOtelResource returns a resource describing this application.
func newOtelResource(app *runtime.Application) *otel_resource.Resource {
	serviceKey := "nanobus"
	version := version.Version
	environment := "non-prod"

	if app != nil {
		if app.ID != "" {
			serviceKey = app.ID
		}
		if app.Version != "" {
			version = app.Version
		}
		if app.Environment != "" {
			environment = app.Environment
		}
	}

	r, _ := otel_resource.Merge(
		otel_resource.Default(),
		otel_resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceKey),
			semconv.ServiceVersionKey.String(version),
			attribute.String("environment", environment),
		),
	)
	return r
}
