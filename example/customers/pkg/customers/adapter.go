package customers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/oklog/run"

	functions "github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/transports/mux"
)

type OutboundImpl struct {
	invoker *functions.Invoker
}

func NewOutboundImpl(invoker *functions.Invoker) *OutboundImpl {
	return &OutboundImpl{
		invoker: invoker,
	}
}

func (m *OutboundImpl) SaveCustomer(ctx context.Context, customer Customer) error {
	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "saveCustomer", customer)
}

func (m *OutboundImpl) FetchCustomer(ctx context.Context, id uint64) (Customer, error) {
	var ret Customer
	inputArgs := hostFetchCustomerArgs{
		ID: id,
	}
	err := m.invoker.InvokeWithReturn(ctx, "customers.v1.Outbound", "fetchCustomer", inputArgs, &ret)
	return ret, err
}

func (m *OutboundImpl) CustomerCreated(ctx context.Context, customer Customer) error {
	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "customerCreated", customer)
}

type handlerGetCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}

type hostFetchCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}

type Adapter struct {
	mux        *mux.Mux
	codec      functions.Codec
	invoker    *functions.Invoker
	registerFn functions.Register

	ln net.Listener
}

func NewAdapter() *Adapter {
	outboundBaseURI := lookupEnvOrString("OUTBOUND_BASE_URI", "http://localhost:9000/outbound/")
	codec := msgpack.New()
	m := mux.New(outboundBaseURI, codec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, codec)

	app := Adapter{
		mux:        m,
		codec:      codec,
		invoker:    invoker,
		registerFn: m.Register,
	}

	return &app
}

func (a *Adapter) Start() (err error) {
	host := lookupEnvOrString("HOST", "localhost")
	port := lookupEnvOrInt("PORT", 9000)
	httpListenAddr := fmt.Sprintf("%s:%d", host, port)
	a.ln, err = net.Listen("tcp", httpListenAddr)
	if err != nil {
		return err
	}
	return http.Serve(a.ln, a.mux.Router())
}

func (a *Adapter) Stop() (err error) {
	if a.ln != nil {
		err = a.ln.Close()
	}
	return err
}

func (a *Adapter) Run() {
	ctx := context.Background()
	var g run.Group
	{
		g.Add(func() error {
			return a.Start()
		}, func(error) {
			a.Stop()
		})
	}
	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}

	if err := g.Run(); err.Error() != "received signal interrupt" {
		log.Fatalln(err)
	}
}

func (a *Adapter) Invoker() *functions.Invoker {
	return a.invoker
}

func (a *Adapter) RegisterInbound(handlers Inbound) *Adapter {
	if handlers.CreateCustomer != nil {
		a.registerFn("customers.v1.Inbound", "createCustomer", a.createCustomerWrapper(handlers.CreateCustomer))
	}
	if handlers.GetCustomer != nil {
		a.registerFn("customers.v1.Inbound", "getCustomer", a.getCustomerWrapper(handlers.GetCustomer))
	}
	return a
}

func (a *Adapter) createCustomerWrapper(handler func(ctx context.Context, customer Customer) (Customer, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var request Customer
		if err := a.codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		response, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(&response)
	}
}

func (a *Adapter) getCustomerWrapper(handler func(ctx context.Context, id uint64) (Customer, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var inputArgs handlerGetCustomerArgs
		if err := a.codec.Decode(payload, &inputArgs); err != nil {
			return nil, err
		}
		response, err := handler(ctx, inputArgs.ID)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(&response)
	}
}

func (a *Adapter) NewOutbound() Outbound {
	return NewOutboundImpl(a.invoker)
}

func lookupEnvOrInt(key string, defaultVal int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
