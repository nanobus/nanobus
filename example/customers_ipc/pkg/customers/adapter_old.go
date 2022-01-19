package customers

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"syscall"

// 	"github.com/oklog/run"

// 	functions "github.com/nanobus/go-functions"
// 	"github.com/nanobus/go-functions/codecs/msgpack"
// 	"github.com/nanobus/go-functions/stateful"
// )

// var busSocketAddress = lookupEnvOrString("BUS_SOCKET_ADDR", "ipc://bus.sock")

// type AdapterContext struct {
// 	*stateful.Context
// }

// func (c *AdapterContext) Self() LogicalAddress {
// 	self := c.Context.Self()
// 	return LogicalAddress{
// 		Type: self.Type,
// 		ID:   self.ID,
// 	}
// }

// type OutboundImpl struct {
// 	invoker *functions.Invoker
// 	codec   functions.Codec
// }

// func NewOutboundImpl(invoker *functions.Invoker, codec functions.Codec) *OutboundImpl {
// 	return &OutboundImpl{
// 		invoker: invoker,
// 		codec:   codec,
// 	}
// }

// // Saves a customer to the backend database
// func (m *OutboundImpl) SaveCustomer(ctx context.Context, customer Customer) error {
// 	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "saveCustomer", customer)
// }

// type customerRecv struct {
// 	s *functions.Stream
// }

// func (s customerRecv) Receive() (*Customer, error) {
// 	var customer Customer
// 	if err := s.s.RecvData(&customer); err != nil {
// 		return nil, err
// 	}
// 	return &customer, nil
// }

// func (s customerRecv) Subscribe(func(customer *Customer), func(error)) {
// 	//return s.s.RecvData(customer)
// }

// type customerSend struct {
// 	s *functions.Stream
// }

// func (s customerSend) Send(customer *Customer) error {
// 	return s.s.SendData(customer)
// }

// func (s customerSend) End() error {
// 	return s.s.Close()
// }

// type getCustomersStream struct {
// 	customerRecv
// 	customerSend
// }

// func (m *OutboundImpl) GetCustomers(ctx context.Context) (CustomerSource, error) {
// 	s, err := m.invoker.InvokeStream(ctx, "customers.v1.Outbound", "getCustomers")
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err = s.Close(); err != nil {
// 		return nil, err
// 	}

// 	//return customerRecv{s}, nil
// 	return getCustomersStream{
// 		customerRecv{s},
// 		customerSend{s},
// 	}, nil
// }

// // Fetches a customer from the backend database
// func (m *OutboundImpl) FetchCustomer(ctx context.Context, id uint64) (*Customer, error) {
// 	var ret Customer
// 	inputArgs := outboundFetchCustomerArgs{
// 		ID: id,
// 	}
// 	if err := m.invoker.InvokeWithReturn(ctx, "customers.v1.Outbound", "fetchCustomer", &inputArgs, &ret); err != nil {
// 		return nil, err
// 	}
// 	return &ret, nil
// }

// // Sends a customer creation event
// func (m *OutboundImpl) CustomerCreated(ctx context.Context, customer Customer) error {
// 	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "customerCreated", customer)
// }

// func lookupEnvOrString(key string, defaultVal string) string {
// 	if val, ok := os.LookupEnv(key); ok {
// 		return val
// 	}
// 	return defaultVal
// }

// type Storage struct {
// 	invoker *functions.Invoker
// }

// func NewStorage(invoker *functions.Invoker) *Storage {
// 	return &Storage{
// 		invoker: invoker,
// 	}
// }

// func (s *Storage) Get(namespace, id, key string) (stateful.RawItem, bool, error) {
// 	var item stateful.RawItem

// 	type Args struct {
// 		Namespace string `json:"namespace" msgpack:"namespace"`
// 		ID        string `json:"id" msgpack:"id"`
// 		Key       string `json:"key" msgpack:"key"`
// 	}

// 	if err := s.invoker.InvokeWithReturn(context.Background(), "nanobus:state", "get", &Args{
// 		Namespace: namespace,
// 		ID:        id,
// 		Key:       key,
// 	}, &item); err != nil {
// 		return item, false, err
// 	}

// 	return item, true, nil
// }

// type Adapter struct {
// 	stateManager       *stateful.Manager
// 	codec              functions.Codec
// 	invoker            *functions.Invoker
// 	registerFn         functions.Register
// 	registerStatefulFn functions.RegisterStateful
// 	listen             func(context.Context) error
// 	close              func() error
// }

// func NewAdapter() (*Adapter, error) {
// 	codec := msgpack.New()
// 	conn := NewConnection("/", codec)

// 	cache, err := stateful.NewLRUCache(200)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not create cache: %w", err)
// 	}
// 	stateInvoker := functions.NewInvoker(conn.Invoke, conn.InvokeStream, codec)
// 	storage := NewStorage(stateInvoker)
// 	stateManager := stateful.NewManager(cache, storage, codec)
// 	invoker := functions.NewInvoker(conn.Invoke, conn.InvokeStream, codec)

// 	app := Adapter{
// 		stateManager:       stateManager,
// 		codec:              codec,
// 		invoker:            invoker,
// 		registerFn:         conn.Register,
// 		registerStatefulFn: conn.RegisterStateful,
// 		listen:             conn.Connect,
// 		close:              conn.Close,
// 	}

// 	return &app, nil
// }

// func (a *Adapter) Start(ctx context.Context) (err error) {
// 	fmt.Printf("🌏 Nanoserver connected to %s\n", busSocketAddress)
// 	err = a.listen(ctx)
// 	return err
// }

// func (a *Adapter) Stop() (err error) {
// 	if a.close != nil {
// 		err = a.close()
// 	}
// 	return err
// }

// func (a *Adapter) Run() {
// 	ctx := context.Background()
// 	var g run.Group
// 	{
// 		g.Add(func() error {
// 			//return nil
// 			return a.Start(ctx)
// 		}, func(error) {
// 			a.Stop()
// 		})
// 	}
// 	{
// 		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
// 	}

// 	if err := g.Run(); err.Error() != "received signal interrupt" {
// 		log.Fatalln(err)
// 	}
// }

// func (a *Adapter) Invoker() *functions.Invoker {
// 	return a.invoker
// }

// func (a *Adapter) RegisterInbound(handlers Inbound) *Adapter {
// 	if handlers.CreateCustomer != nil {
// 		a.registerFn("customers.v1.Inbound", "createCustomer", a.inbound_createCustomerWrapper(handlers.CreateCustomer))
// 	}
// 	if handlers.GetCustomer != nil {
// 		a.registerFn("customers.v1.Inbound", "getCustomer", a.inbound_getCustomerWrapper(handlers.GetCustomer))
// 	}
// 	if handlers.ListCustomers != nil {
// 		a.registerFn("customers.v1.Inbound", "listCustomers", a.inbound_listCustomersWrapper(handlers.ListCustomers))
// 	}
// 	return a
// }

// func (a *Adapter) RegisterCustomerActor(stateful CustomerActor) *Adapter {
// 	a.registerStatefulFn("customers.v1.CustomerActor", "deactivate", a.stateManager.DeactivateHandler("customers.v1.CustomerActor", stateful))
// 	a.registerStatefulFn("customers.v1.CustomerActor", "createCustomer", a.customerActor_createCustomerWrapper(stateful))
// 	a.registerStatefulFn("customers.v1.CustomerActor", "getCustomer", a.customerActor_getCustomerWrapper(stateful))
// 	return a
// }

// func (a *Adapter) inbound_createCustomerWrapper(handler func(ctx context.Context, customer Customer) (*Customer, error)) functions.Handler {
// 	return func(ctx context.Context, payload []byte) ([]byte, error) {
// 		var request Customer
// 		if err := a.codec.Decode(payload, &request); err != nil {
// 			return nil, err
// 		}
// 		response, err := handler(ctx, request)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return a.codec.Encode(response)
// 	}
// }

// func (a *Adapter) inbound_getCustomerWrapper(handler func(ctx context.Context, id uint64) (*Customer, error)) functions.Handler {
// 	return func(ctx context.Context, payload []byte) ([]byte, error) {
// 		var inputArgs inboundGetCustomerArgs
// 		if err := a.codec.Decode(payload, &inputArgs); err != nil {
// 			return nil, err
// 		}
// 		response, err := handler(ctx, inputArgs.ID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return a.codec.Encode(response)
// 	}
// }

// func (a *Adapter) inbound_listCustomersWrapper(handler func(ctx context.Context, query CustomerQuery) (*CustomerPage, error)) functions.Handler {
// 	return func(ctx context.Context, payload []byte) ([]byte, error) {
// 		var request CustomerQuery
// 		if err := a.codec.Decode(payload, &request); err != nil {
// 			return nil, err
// 		}
// 		response, err := handler(ctx, request)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return a.codec.Encode(response)
// 	}
// }

// func (a *Adapter) customerActor_createCustomerWrapper(stateful CustomerActor) functions.StatefulHandler {
// 	return func(ctx context.Context, id string, payload []byte) ([]byte, error) {
// 		var request Customer
// 		if err := a.codec.Decode(payload, &request); err != nil {
// 			return nil, err
// 		}
// 		sctx, err := a.stateManager.ToContext(ctx, "customers.v1.CustomerActor", id, stateful)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response, err := stateful.CreateCustomer(&AdapterContext{&sctx}, request)
// 		if err != nil {
// 			return nil, err
// 		}
// 		statefulResponse, err := sctx.Response(response)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return a.codec.Encode(statefulResponse)
// 	}
// }

// func (a *Adapter) customerActor_getCustomerWrapper(stateful CustomerActor) functions.StatefulHandler {
// 	return func(ctx context.Context, id string, payload []byte) ([]byte, error) {
// 		sctx, err := a.stateManager.ToContext(ctx, "customers.v1.CustomerActor", id, stateful)
// 		if err != nil {
// 			return nil, err
// 		}
// 		response, err := stateful.GetCustomer(&AdapterContext{&sctx})
// 		if err != nil {
// 			return nil, err
// 		}
// 		statefulResponse, err := sctx.Response(response)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return a.codec.Encode(statefulResponse)
// 	}
// }

// type inboundGetCustomerArgs struct {
// 	ID uint64 `json:"id" msgpack:"id"`
// }

// func (a *Adapter) NewOutbound() Outbound {
// 	return NewOutboundImpl(a.invoker, a.codec)
// }

// type outboundFetchCustomerArgs struct {
// 	ID uint64 `json:"id" msgpack:"id"`
// }
