package customers

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/oklog/run"

	functions "github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/stateful"
	"github.com/nanobus/go-functions/transports/mux"
)

var busURI = lookupEnvOrString("BUS_URI", "http://localhost:32321")

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

func (m *OutboundImpl) FetchCustomer(ctx context.Context, id uint64) (*Customer, error) {
	var ret Customer
	inputArgs := outboundFetchCustomerArgs{
		ID: id,
	}
	err := m.invoker.InvokeWithReturn(ctx, "customers.v1.Outbound", "fetchCustomer", inputArgs, &ret)
	return &ret, err
}

func (m *OutboundImpl) CustomerCreated(ctx context.Context, customer Customer) error {
	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "customerCreated", customer)
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

type Storage struct {
	codec functions.Codec
}

func NewStorage(codec functions.Codec) *Storage {
	return &Storage{
		codec: codec,
	}
}

func (s *Storage) Get(namespace, id, key string) (stateful.RawItem, bool, error) {
	var item stateful.RawItem
	url := busURI + "/state/" + namespace + "/" + id + "/" + key
	resp, err := http.Get(url)
	if err != nil {
		return item, false, err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return item, false, err
	}

	if len(payload) == 0 {
		return item, false, nil
	}
	if err = s.codec.Decode(payload, &item); err != nil {
		return item, false, err
	}

	return item, true, nil
}

type Adapter struct {
	mux                *mux.Mux
	stateManager       *stateful.Manager
	codec              functions.Codec
	invoker            *functions.Invoker
	registerFn         functions.Register
	registerStatefulFn functions.RegisterStateful

	ln net.Listener
}

func NewAdapter(stateManager *stateful.Manager) *Adapter {
	codec := msgpack.New()
	m := mux.New(busURI+"/outbound/", codec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, codec)

	app := Adapter{
		mux:                m,
		stateManager:       stateManager,
		codec:              codec,
		invoker:            invoker,
		registerFn:         m.Register,
		registerStatefulFn: m.RegisterStateful,
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
		a.registerFn("customers.v1.Inbound", "createCustomer", a.inbound_createCustomerWrapper(handlers.CreateCustomer))
	}
	if handlers.GetCustomer != nil {
		a.registerFn("customers.v1.Inbound", "getCustomer", a.inbound_getCustomerWrapper(handlers.GetCustomer))
	}
	if handlers.ListCustomers != nil {
		a.registerFn("customers.v1.Inbound", "listCustomers", a.inbound_listCustomersWrapper(handlers.ListCustomers))
	}
	return a
}

func (a *Adapter) RegisterCustomerActor(stateful CustomerActor) *Adapter {
	a.registerStatefulFn("customers.v1.CustomerActor", "deactivate", a.stateManager.DeactivateHandler("customers.v1.CustomerActor", stateful))
	a.registerStatefulFn("customers.v1.CustomerActor", "createCustomer", a.customerActor_createCustomerWrapper(stateful))
	a.registerStatefulFn("customers.v1.CustomerActor", "getCustomer", a.customerActor_getCustomerWrapper(stateful))
	return a
}

func (a *Adapter) inbound_createCustomerWrapper(handler func(ctx context.Context, customer Customer) (*Customer, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var request Customer
		if err := a.codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		response, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(response)
	}
}

func (a *Adapter) inbound_getCustomerWrapper(handler func(ctx context.Context, id uint64) (*Customer, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var inputArgs inboundGetCustomerArgs
		if err := a.codec.Decode(payload, &inputArgs); err != nil {
			return nil, err
		}
		response, err := handler(ctx, inputArgs.ID)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(response)
	}
}

func (a *Adapter) inbound_listCustomersWrapper(handler func(ctx context.Context, query CustomerQuery) (*CustomerPage, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var request CustomerQuery
		if err := a.codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		response, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(response)
	}
}

func (a *Adapter) customerActor_createCustomerWrapper(stateful CustomerActor) functions.StatefulHandler {
	return func(ctx context.Context, id string, payload []byte) ([]byte, error) {
		var request Customer
		if err := a.codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		sctx, err := a.stateManager.ToContext(ctx, "customers.v1.CustomerActor", id, stateful)
		if err != nil {
			return nil, err
		}
		response, err := stateful.CreateCustomer(sctx, request)
		if err != nil {
			return nil, err
		}
		statefulResponse, err := sctx.Response(response)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(statefulResponse)
	}
}

func (a *Adapter) customerActor_getCustomerWrapper(stateful CustomerActor) functions.StatefulHandler {
	return func(ctx context.Context, id string, payload []byte) ([]byte, error) {
		sctx, err := a.stateManager.ToContext(ctx, "customers.v1.CustomerActor", id, stateful)
		if err != nil {
			return nil, err
		}
		response, err := stateful.GetCustomer(sctx)
		if err != nil {
			return nil, err
		}
		statefulResponse, err := sctx.Response(response)
		if err != nil {
			return nil, err
		}
		return a.codec.Encode(statefulResponse)
	}
}

type inboundGetCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}

func (a *Adapter) NewOutbound() Outbound {
	return NewOutboundImpl(a.invoker)
}

type outboundFetchCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}
