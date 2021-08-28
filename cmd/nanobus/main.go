package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/dapr/components-contrib/bindings"
	"github.com/dapr/components-contrib/pubsub"
	dapr_rt "github.com/dapr/dapr/pkg/runtime"
	"github.com/dapr/dapr/pkg/runtime/embedded"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/nanobus/go-functions"
	json_codec "github.com/nanobus/go-functions/codecs/json"
	msgpack_codec "github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/oklog/run"
	"github.com/vmihailenco/msgpack/v5"

	"github.com/nanobus/nanobus/actions"
	"github.com/nanobus/nanobus/actions/core"
	"github.com/nanobus/nanobus/actions/dapr"
	"github.com/nanobus/nanobus/coalesce"
	"github.com/nanobus/nanobus/codec"
	cloudevents_avro "github.com/nanobus/nanobus/codec/cloudevents/avro"
	"github.com/nanobus/nanobus/codec/confluentavro"
	codec_json "github.com/nanobus/nanobus/codec/json"
	codec_msgpack "github.com/nanobus/nanobus/codec/msgpack"
	"github.com/nanobus/nanobus/compute"
	compute_mux "github.com/nanobus/nanobus/compute/mux"
	compute_wapc "github.com/nanobus/nanobus/compute/wapc"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/runtime"
	dapr_runtime "github.com/nanobus/nanobus/runtime/dapr"
	"github.com/nanobus/nanobus/spec"
	spec_widl "github.com/nanobus/nanobus/spec/widl"
	"github.com/nanobus/nanobus/transport"
	"github.com/nanobus/nanobus/transport/httprpc"
	"github.com/nanobus/nanobus/transport/rest"
)

var ErrInvalidURISyntax = errors.New("invalid invocation URI syntax")

type decoding struct {
	Pubsub  []pubsubDecoding  `mapstructure:"pubsub"`
	Binding []bindingDecoding `mapstructure:"binding"`
}

type pubsubDecoding struct {
	PubsubName string        `mapstructure:"pubsubname"`
	Topic      string        `mapstructure:"topic"`
	Codec      string        `mapstructure:"codec"`
	Args       []interface{} `mapstructure:"args"`
}

type bindingDecoding struct {
	BindingName string        `mapstructure:"bindingname"`
	Codec       string        `mapstructure:"codec"`
	Args        []interface{} `mapstructure:"args"`
}

type pubsubKey struct {
	PubsubName string `mapstructure:"pubsubname"`
	Topic      string `mapstructure:"topic"`
}

type codecConfig struct {
	codec codec.Codec
	args  []interface{}
}

type pubsubDecoders map[pubsubKey]codecConfig
type bindingDecoders map[string]codecConfig

type Runtime struct {
	config     *runtime.Configuration
	namespaces spec.Namespaces
	processor  *runtime.Processor
	resolver   resolve.DependencyResolver
	resolveAs  resolve.ResolveAs
	env        runtime.Environment
}

func main() {
	daprRuntime := dapr_runtime.New()
	daprRuntime.AttachFlags()

	// NanoBus flags

	var httpListenAddr string
	flag.StringVar(
		&httpListenAddr,
		"http-listen-addr",
		LookupEnvOrString("HTTP_LISTEN_ADDR", ":8080"),
		"http listen address",
	)
	var busListenAddr string
	flag.StringVar(
		&busListenAddr,
		"bus-listen-addr",
		LookupEnvOrString("BUS_LISTEN_ADDR", "localhost:32320"),
		"bus listen address",
	)
	var restListenAddr string
	flag.StringVar(
		&restListenAddr,
		"rest-listen-addr",
		LookupEnvOrString("REST_LISTEN_ADDR", ":8090"),
		"rest listen address",
	)
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("usage: nanobus <configuration file>")
		os.Exit(1)
	}

	if err := daprRuntime.Initialize(); err != nil {
		log.Fatal(err)
	}

	filename := args[0]
	// Load the configuration
	config, err := loadConfiguration(filename)
	if err != nil {
		log.Fatal(err)
	}

	daprComponents := dapr.DaprComponents{}
	dapr_runtime.RegisterComponents(daprRuntime)
	daprRuntime.AddOptions(dapr_rt.WithEmbeddedHandlers(&daprComponents))
	if err = daprRuntime.Run(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// Spec registration
	specRegistry := spec.Registry{}
	specRegistry.Register(
		spec_widl.WIDL,
	)

	// Compute registration
	computeRegistry := compute.Registry{}
	computeRegistry.Register(
		compute_mux.Mux,
		compute_wapc.WaPC,
	)

	// Codec registration
	codecRegistry := codec.Registry{}
	codecRegistry.Register(
		codec_json.JSON,
		codec_msgpack.MsgPack,
		confluentavro.ConfluentAvro,
		cloudevents_avro.CloudEventsAvro,
	)

	// Action registration
	actionRegistry := actions.Registry{}
	actionRegistry.Register(core.All...)
	actionRegistry.Register(dapr.All...)

	namespaces := make(spec.Namespaces)
	for _, spec := range config.Specs {
		loader, ok := specRegistry[spec.Type]
		if !ok {
			log.Fatal(fmt.Errorf("could not find spec type %q", spec.Type))
		}
		nss, err := loader(spec.With)
		if err != nil {
			log.Fatal(fmt.Errorf("error loading spec of type %q", spec.Type))
		}
		for _, ns := range nss {
			namespaces[ns.Name] = ns
		}
	}

	// Codecs
	jsoncodec := json_codec.New()
	msgpackcodec := msgpack_codec.New()

	// Dependencies
	var invoker *functions.Invoker
	var busInvoker compute.BusInvoker
	httpClient := getHTTPClient()
	env := getEnvironment()
	dependencies := map[string]interface{}{
		"client:http":     httpClient,
		"codec:json":      jsoncodec,
		"codec:msgpack":   msgpackcodec,
		"spec:namespaces": namespaces,
		"os:env":          env,
		"dapr:components": &daprComponents,
	}
	resolver := func(name string) (interface{}, bool) {
		dep, ok := dependencies[name]
		return dep, ok
	}
	resolveAs := resolve.ToResolveAs(resolver)

	if config.Codecs == nil {
		config.Codecs = map[string]runtime.Codec{}
	}
	if _, exists := config.Codecs["json"]; !exists {
		config.Codecs["json"] = runtime.Codec{
			Type: "json",
		}
	}
	if _, exists := config.Codecs["msgpack"]; !exists {
		config.Codecs["msgpack"] = runtime.Codec{
			Type: "msgpack",
		}
	}
	codecs := make(codec.Codecs)
	for name, codec := range config.Codecs {
		loader, ok := codecRegistry[codec.Type]
		if !ok {
			log.Fatal(fmt.Errorf("could not find codec type %q", codec.Type))
		}
		c, err := loader(codec.With, resolveAs)
		if err != nil {
			log.Fatal(fmt.Errorf("error loading codec of type %q", codec.Type))
		}
		codecs[name] = c
	}
	dependencies["codec:lookup"] = codecs

	pubsubDecs := pubsubDecoders{}
	bindingDecs := bindingDecoders{}
	if config.Decoding != nil {
		var dec decoding
		err = mapstructure.Decode(config.Decoding, &dec)
		if err != nil {
			log.Fatal(fmt.Errorf("error reading decoding config %w", err))
		}

		for _, psd := range dec.Pubsub {
			codec, ok := codecs[psd.Codec]
			if !ok {
				log.Fatal(fmt.Errorf("codec %q is not configured", psd.Codec))
			}
			if psd.Args == nil {
				psd.Args = []interface{}{}
			}
			pubsubDecs[pubsubKey{
				PubsubName: psd.PubsubName,
				Topic:      psd.Topic,
			}] = codecConfig{
				codec: codec,
				args:  psd.Args,
			}
		}

		for _, bd := range dec.Binding {
			codec, ok := codecs[bd.Codec]
			if !ok {
				log.Fatal(fmt.Errorf("codec %q is not configured", bd.Codec))
			}
			if bd.Args == nil {
				bd.Args = []interface{}{}
			}
			bindingDecs[bd.BindingName] = codecConfig{
				codec: codec,
				args:  bd.Args,
			}
		}
	}

	// Create processor
	processor, err := runtime.New(config, actionRegistry, resolver)
	if err != nil {
		log.Fatal(err)
	}

	rt := Runtime{
		config:     config,
		namespaces: namespaces,
		processor:  processor,
		resolver:   resolver,
		resolveAs:  resolveAs,
		env:        env,
	}
	busInvoker = rt.BusInvoker
	dependencies["bus:invoker"] = busInvoker

	// Internal invoker
	if config.Compute.Type == "" {
		config.Compute.Type = "mux"
	}
	computeLoader, ok := computeRegistry[config.Compute.Type]
	if !ok {
		log.Fatal(fmt.Errorf("could not find compute type %q", config.Compute.Type))
	}
	invoker, err = computeLoader(config.Compute.With, resolveAs)
	if err != nil {
		log.Fatal(err)
	}
	dependencies["client:invoker"] = invoker

	if err = processor.Initialize(); err != nil {
		log.Fatal(err)
	}

	daprComponents.InvokeHandler(func(ctx context.Context, method, contentType string, payload []byte, metadata map[string][]string) ([]byte, string, error) {
		var input interface{}
		// TODO: Decoder
		if err := json.Unmarshal(payload, &input); err != nil {
			return nil, "", err
		}

		target := method
		idx := strings.LastIndex(method, "/")
		if idx < 0 {
			return nil, "", fmt.Errorf("invalid method %q", method)
		}
		function := method[idx+1:]
		method = method[:idx]

		lastDot := strings.LastIndexByte(method, '.')
		if lastDot < 0 {
			return nil, "", fmt.Errorf("invalid method %q", method)
		}
		service := method[lastDot+1:]
		namespace := method[:lastDot]

		data := actions.Data{
			"input":    input,
			"metadata": metadata,
			"env":      env,
		}

		if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
			log.Println("-->", target, string(jsonBytes)+"\n")
		}

		output, _, err := rt.processor.Service(ctx, namespace, service, function, data)
		if err != nil {
			return nil, "", err
		}

		var respData []byte
		if output != nil {
			if respData, err = json.Marshal(&output); err != nil {
				return nil, "", err
			}
		}

		return respData, "", err
	})

	inputBindingHandler := func(function string, codec codec.Codec, args []interface{}) func(*bindings.ReadResponse) ([]byte, error) {
		return func(msg *bindings.ReadResponse) ([]byte, error) {
			target := function

			decoded, typeName, err := codec.Decode(msg.Data, args...)
			if err != nil {
				log.Println("error decoding event payload", err)
				return nil, err
			}
			input := map[string]interface{}{
				"data": decoded,
				"type": typeName,
			}

			data := actions.Data{
				"input":    input,
				"metadata": msg.Metadata,
				"env":      env,
			}

			if target == "" {
				target = typeName
			}

			// if jsonBytes, err := json.MarshalIndent(data, "", "  "); err == nil {
			// 	log.Println("-->", target, string(jsonBytes)+"\n")
			// }

			result, err := rt.processor.Inbound(ctx, target, data)
			if err != nil {
				log.Println("error processing event", err)
				return nil, err
			}

			var resultBytes []byte
			if result != nil {
				resultBytes, err = codec.Encode(result)
			}

			return resultBytes, err
		}
	}

	type InputBinding struct {
		Binding   string        `mapstructure:"binding"`
		Codec     string        `mapstructure:"codec"`
		CodecArgs []interface{} `mapstructure:"codecArgs"`
		Function  string        `mapstructure:"function"`
	}

	var inputBindings []InputBinding
	if rt.config.InputBindings != nil {
		if err := mapstructure.Decode(rt.config.InputBindings, &inputBindings); err != nil {
			log.Fatal(err)
		}
	}

	// Direct input bindings
	for _, binding := range inputBindings {
		p, ok := daprComponents.InputBindings[binding.Binding]
		if !ok {
			log.Fatal(fmt.Errorf("input binding %q is not configured", binding.Binding))
		}
		if binding.CodecArgs == nil {
			binding.CodecArgs = []interface{}{}
		}
		c, ok := codecs[binding.Codec]
		if !ok {
			log.Fatal(fmt.Errorf("codec %q is not configured", binding.Codec))
		}

		go func(p bindings.InputBinding, binding InputBinding, c codec.Codec) {
			log.Printf("reading from input binding %q", binding.Binding)
			if err = p.Read(inputBindingHandler(binding.Function, c, binding.CodecArgs)); err != nil {
				log.Println(err)
			}
		}(p, binding, c)
	}

	daprComponents.InputBindingHandler(func(ctx context.Context, event *embedded.BindingEvent) ([]byte, error) {
		var input interface{}
		if codec, ok := bindingDecs[event.BindingName]; ok {
			decoded, typeName, err := codec.codec.Decode(event.Data, codec.args...)
			if err != nil {
				log.Println("error decoding input binding payload", err)
				return nil, err
			}
			input = map[string]interface{}{
				"data": decoded,
				"type": typeName,
			}
		} else if err := json.Unmarshal(event.Data, &input); err != nil {
			return nil, err
		}

		actionData := actions.Data{
			"input":    input,
			"metadata": event.Metadata,
			"env":      env,
		}

		output, err := rt.processor.Inbound(ctx, event.BindingName, actionData)
		if err != nil {
			return nil, err
		}

		var respData []byte
		if output != nil {
			if respData, err = json.Marshal(&output); err != nil {
				return nil, err
			}
		}

		return respData, err
	})

	pubsubHandler := func(function string, codec codec.Codec, args []interface{}) pubsub.Handler {
		return func(ctx context.Context, msg *pubsub.NewMessage) error {
			target := function

			decoded, typeName, err := codec.Decode(msg.Data, args...)
			if err != nil {
				log.Println("error decoding event payload", err)
				return err
			}
			input := map[string]interface{}{
				"data": decoded,
				"type": typeName,
			}

			data := actions.Data{
				"input":    input,
				"metadata": msg.Metadata,
				"env":      env,
			}

			if target == "" {
				target = typeName
			}

			_, err = rt.processor.Inbound(ctx, target, data)
			if err != nil {
				log.Println("error processing event", err)
				return err
			}

			return nil
		}
	}

	type Subscription struct {
		Pubsub    string            `mapstructure:"pubsub"`
		Topic     string            `mapstructure:"topic"`
		Metadata  map[string]string `mapstructure:"metadata"`
		Codec     string            `mapstructure:"codec"`
		CodecArgs []interface{}     `mapstructure:"codecArgs"`
		Function  string            `mapstructure:"function"`
	}

	var subscriptions []Subscription
	if err := mapstructure.Decode(rt.config.Subscriptions, &subscriptions); err != nil {
		log.Fatal(err)
	}

	// Direct subscriptions
	for _, sub := range subscriptions {
		p, ok := daprComponents.PubSubs[sub.Pubsub]
		if !ok {
			log.Fatal(fmt.Errorf("pubsub %q is not configured", sub.Pubsub))
		}
		if sub.CodecArgs == nil {
			sub.CodecArgs = []interface{}{}
		}
		codec, ok := codecs[sub.Codec]
		if !ok {
			log.Fatal(fmt.Errorf("codec %q is not configured", sub.Codec))
		}
		log.Printf("subscribing to pubsub %q, topic %q", sub.Pubsub, sub.Topic)
		if err = p.Subscribe(pubsub.SubscribeRequest{
			Topic:    sub.Topic,
			Metadata: sub.Metadata,
		}, pubsubHandler(sub.Function, codec, sub.CodecArgs)); err != nil {
			log.Fatal(err)
		}
	}

	// Subscriptions via the embedded app channel
	daprComponents.PubSubHandler(func(ctx context.Context, event *embedded.TopicEvent) (embedded.EventResponseStatus, error) {
		function := event.Path

		if jsonBytes, err := json.MarshalIndent(event.CloudEvent, "", "  "); err == nil {
			log.Println("-->", function, string(jsonBytes)+"\n")
		}

		var input interface{} = event.CloudEvent
		if event.RawPayload {
			dataBase64String, ok := event.CloudEvent["data_base64"].(string)
			if !ok {
				log.Println("error decoding raw payload", err)
				return embedded.EventResponseStatusRetry, err
			}
			inputBytes, err := base64.StdEncoding.DecodeString(dataBase64String)
			if err != nil {
				log.Println("error decoding raw payload", err)
				return embedded.EventResponseStatusRetry, err
			}

			if codec, ok := pubsubDecs[pubsubKey{
				PubsubName: event.PubsubName,
				Topic:      event.Topic,
			}]; ok {
				decoded, typeName, err := codec.codec.Decode(inputBytes, codec.args...)
				if err != nil {
					log.Println("error decoding raw payload", err)
					return embedded.EventResponseStatusRetry, err
				}
				input = map[string]interface{}{
					"data": decoded,
					"type": typeName,
				}
			} else {
				input = inputBytes
			}
		}

		data := actions.Data{
			"input":    input,
			"metadata": event.Metadata,
			"env":      env,
		}

		_, err := rt.processor.Inbound(ctx, function, data)
		if err != nil {
			log.Println("error processing event", err)
			return embedded.EventResponseStatusRetry, err
		}

		return embedded.EventResponseStatusSuccess, nil
	})

	transportInvoker := func(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error) {
		if err := coalesceInput(namespaces, namespace, service, function, input); err != nil {
			return nil, err
		}

		if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
			log.Println("-->", namespace+"."+service+"/"+function, string(jsonBytes)+"\n")
		}

		data := actions.Data{
			"input": input,
			"env":   env,
		}

		response, ok, err := rt.processor.Service(ctx, namespace, service, function, data)
		if err != nil {
			return nil, err
		}

		if !ok {
			// No pipeline exits for the operation so invoke directly.
			ns := namespace + "." + service
			if err = invoker.InvokeWithReturn(ctx, ns, function, input, &response); err != nil {
				return nil, err
			}
		}

		return response, err
	}

	var g run.Group
	{
		g.Add(func() error {
			return daprRuntime.WaitUntilShutdown()
		}, func(error) {
			daprRuntime.Shutdown(1 * time.Second)
		})
	}
	{
		// Expose the bus
		r := mux.NewRouter()
		r.HandleFunc("/outbound/{namespace}/{function}", rt.OutboundHandler).Methods("POST")
		r.HandleFunc("/inbound/{function}", rt.InboundHandler).Methods("POST")
		//r.HandleFunc("/dapr/subscribe", rt.SubscriptionsHandler).Methods("GET")

		log.Printf("Bus listening on %s\n", busListenAddr)
		ln, err := net.Listen("tcp", busListenAddr)
		if err != nil {
			log.Fatalln(err)
		}
		g.Add(func() error {
			return http.Serve(ln, r)
		}, func(error) {
			ln.Close()
		})
	}
	{
		// Expose HTTP-RPC
		transport, err := httprpc.New(httpListenAddr, namespaces, transportInvoker, jsoncodec, msgpackcodec)
		if err != nil {
			log.Fatal(err)
		}
		g.Add(func() error {
			log.Printf("HTTP-RPC listening on %s\n", httpListenAddr)
			return transport.Listen()
		}, func(error) {
			transport.Close()
		})
	}
	{
		// Expose REST
		transport, err := rest.New(restListenAddr, namespaces, transportInvoker, jsoncodec, msgpackcodec)
		if err != nil {
			log.Fatal(err)
		}
		g.Add(func() error {
			log.Printf("REST listening on %s\n", restListenAddr)
			return transport.Listen()
		}, func(error) {
			transport.Close()
		})
	}
	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}

	err = g.Run()
	if _, isSignal := err.(run.SignalError); !isSignal {
		log.Fatalln(err)
	}
}

func (rt *Runtime) OutboundHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	namespace := mux.Vars(r)["namespace"]
	function := mux.Vars(r)["function"]

	lastDot := strings.LastIndexByte(namespace, '.')
	if lastDot < 0 {
		handleError(ErrInvalidURISyntax, w, http.StatusBadRequest)
		return
	}
	service := namespace[lastDot+1:]
	namespace = namespace[:lastDot]

	var input interface{}
	if err := msgpack.NewDecoder(r.Body).Decode(&input); err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	data := actions.Data{
		"input": input,
		"env":   rt.env,
	}

	if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
		log.Println("<--", namespace+"."+service+"/"+function, string(jsonBytes)+"\n")
	}

	output, err := rt.processor.Outbound(r.Context(), namespace, service, function, data)
	if err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	if output, err = coalesceOutput(rt.namespaces, namespace, service, function, output); err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/msgpack")
	msgpack.NewEncoder(w).Encode(output)
}

func (rt *Runtime) BusInvoker(ctx context.Context, namespace, service, function string, input interface{}) (interface{}, error) {
	if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
		log.Println("<--", namespace+"."+service+"/"+function, string(jsonBytes)+"\n")
	}

	data := actions.Data{
		"input": input,
		"env":   rt.env,
	}

	output, err := rt.processor.Outbound(ctx, namespace, service, function, data)
	if err != nil {
		return nil, err
	}

	if output, err = coalesceOutput(rt.namespaces, namespace, service, function, output); err != nil {
		return nil, err
	}

	return output, nil
}

func (rt *Runtime) InboundHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	function := mux.Vars(r)["function"]

	var input interface{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}
	input = coalesce.Integers(input)

	if jsonBytes, err := json.MarshalIndent(input, "", "  "); err == nil {
		log.Println("-->", function, string(jsonBytes))
	}

	data := actions.Data{
		"input": input,
		"env":   rt.env,
	}

	output, err := rt.processor.Inbound(r.Context(), function, data)
	if err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (rt *Runtime) SubscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	type Subscription struct {
		Pubsub   string            `mapstructure:"pubsub"`
		Topic    string            `mapstructure:"topic"`
		Metadata map[string]string `mapstructure:"metadata"`
		Function string            `mapstructure:"function"`
	}

	var subscriptions []Subscription
	if err := mapstructure.Decode(rt.config.Subscriptions, &subscriptions); err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	type DaprSupscription struct {
		Pubsubname string            `json:"pubsubname"`
		Topic      string            `json:"topic"`
		Metadata   map[string]string `json:"metadata"`
		Route      string            `json:"route"`
	}

	daprSubs := make([]DaprSupscription, len(subscriptions))
	for i, sub := range subscriptions {
		daprSubs[i] = DaprSupscription{
			Pubsubname: sub.Pubsub,
			Topic:      sub.Topic,
			Metadata:   sub.Metadata,
			Route:      "/inbound/" + sub.Function,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(daprSubs)
}

func handleError(err error, w http.ResponseWriter, status int) {
	log.Println(err)
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %v", err)
}

func loadConfiguration(filename string) (*runtime.Configuration, error) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return runtime.LoadYAML(f)
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

func coalesceInput(namespaces spec.Namespaces, namespace, service, function string, input interface{}) error {
	if oper, ok := namespaces.Operation(namespace, service, function); ok {
		if oper.Parameters != nil {
			inputMap, ok := coalesce.ToMapSI(input)
			if !ok {
				return fmt.Errorf("%w: input is not a map", transport.ErrBadInput)
			}
			input = inputMap
			if err := oper.Parameters.Coalesce(inputMap, true); err != nil {
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
		if oper.Returns != nil {
			outputMap, ok := coalesce.ToMapSI(output)
			if !ok {
				return nil, errors.New("output is not a map")
			}
			output, err = oper.Returns.Coalesce(outputMap, true)
		}
	} else {
		coalesce.Integers(output)
	}
	return output, err
}
