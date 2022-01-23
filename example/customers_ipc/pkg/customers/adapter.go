package customers

import (
	"context"
	"errors"
	"io"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/json"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/metadata"
	"github.com/nanobus/go-functions/stateful"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

var ErrNotFound = errors.New("not_found")

type (
	RequestResponseHandler func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink)
	RequestStreamHandler   func(ctx context.Context, md metadata.MD, request payload.Payload, sink flux.Sink)
	RequestChannelHandler  func(ctx context.Context, md metadata.MD, requests flux.Flux, sink flux.Sink)
)

type Adapter struct {
	starter      rsocket.ClientStarter
	client       rsocket.Client
	codec        functions.Codec
	stateManager *stateful.Manager
	handlersRR   map[string]RequestResponseHandler
	handlersRS   map[string]RequestStreamHandler
	handlersRC   map[string]RequestChannelHandler
	done         chan struct{}
}

func NewAdapter() *Adapter {
	a := Adapter{
		codec:      msgpack.New(),
		handlersRR: make(map[string]RequestResponseHandler),
		handlersRS: make(map[string]RequestStreamHandler),
		handlersRC: make(map[string]RequestChannelHandler),
		done:       make(chan struct{}),
	}

	cache, _ := stateful.NewLRUCache(200)
	storage := NewStorage(&a)
	a.stateManager = stateful.NewManager(cache, storage, json.New())

	// Start a client connection
	contentType := a.codec.ContentType()
	tp := rsocket.TCPClient().SetHostAndPort("127.0.0.1", 7878).Build()

	a.starter = rsocket.Connect().
		OnClose(func(error) {
			close(a.done)
		}).
		DataMimeType(contentType).
		MetadataMimeType(contentType).
		Acceptor(func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
			return rsocket.NewAbstractSocket(
				rsocket.RequestResponse(a.requestResponseHandler),
				rsocket.RequestStream(a.requestStreamHandler),
				rsocket.RequestChannel(a.requestChannelHandler),
			)
		}).
		Transport(tp)

	return &a
}

func (a *Adapter) Start(ctx context.Context) error {
	client, err := a.starter.Start(ctx)
	if err != nil {
		return err
	}
	a.client = client
	<-a.done
	return nil
}

func (a *Adapter) parseMetadata(pl payload.Payload) (metadata.MD, bool) {
	if mdBytes, ok := pl.Metadata(); ok {
		var md metadata.MD
		if err := a.codec.Decode(mdBytes, &md); err == nil {
			return md, true
		}
	}
	return nil, false
}

func (a *Adapter) requestResponseHandler(request payload.Payload) mono.Mono {
	return mono.Create(func(ctx context.Context, sink mono.Sink) {
		md, ok := a.parseMetadata(request)
		if !ok {
			sink.Error(ErrNotFound)
			return
		}
		path, _ := md.Scalar(":path")
		handler, ok := a.handlersRR[path]
		if !ok {
			sink.Error(ErrNotFound)
			return
		}

		handler(ctx, md, request, sink)
	})
}

func (a *Adapter) requestStreamHandler(request payload.Payload) (responses flux.Flux) {
	return flux.Create(func(ctx context.Context, sink flux.Sink) {
		md, ok := a.parseMetadata(request)
		if !ok {
			sink.Error(ErrNotFound)
			return
		}
		path, _ := md.Scalar(":path")
		handler, ok := a.handlersRS[path]
		if !ok {
			sink.Error(ErrNotFound)
			return
		}

		handler(ctx, md, request, sink)
	})
}

func (a *Adapter) requestChannelHandler(requests flux.Flux) (responses flux.Flux) {
	return flux.Create(func(ctx context.Context, sink flux.Sink) {
		request, err := requests.BlockFirst(ctx)
		if err != nil {
			sink.Error(err)
			return
		}

		md, ok := a.parseMetadata(request)
		if !ok {
			sink.Error(ErrNotFound)
			return
		}
		path, _ := md.Scalar(":path")
		handler, ok := a.handlersRS[path]
		if !ok {
			sink.Error(ErrNotFound)
			return
		}

		handler(ctx, md, request, sink)
	})
}

func (a *Adapter) NewOutbound() Outbound {
	return NewOutboundImpl(a)
}

func (a *Adapter) RegisterRR(path string, handler RequestResponseHandler) {
	a.handlersRR[path] = handler
}

func (a *Adapter) RegisterRS(path string, handler RequestStreamHandler) {
	a.handlersRS[path] = handler
}

func (a *Adapter) RegisterRC(path string, handler RequestChannelHandler) {
	a.handlersRC[path] = handler
}

func (a *Adapter) RegisterInbound(handlers Inbound) {
	a.RegisterRR("/customers.v1.Inbound/createCustomer", a.inbound_createCustomerWrapper(handlers.CreateCustomer))
	a.RegisterRR("/customers.v1.Inbound/getCustomer", a.inbound_getCustomerWrapper(handlers.GetCustomer))
	a.RegisterRR("/customers.v1.Inbound/listCustomers", a.inbound_listCustomersWrapper(handlers.ListCustomers))
}

func (a *Adapter) inbound_createCustomerWrapper(handler func(ctx context.Context, customer Customer) (*Customer, error)) RequestResponseHandler {
	return func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink) {
		var customer Customer
		a.invokeRR(request, sink, &customer,
			func() (interface{}, error) {
				return handler(ctx, customer)
			},
		)
	}
}

func (a *Adapter) inbound_getCustomerWrapper(handler func(ctx context.Context, id uint64) (*Customer, error)) RequestResponseHandler {
	return func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink) {
		var inputArgs inboundGetCustomerArgs
		a.invokeRR(request, sink, &inputArgs,
			func() (interface{}, error) {
				return handler(ctx, inputArgs.ID)
			},
		)
	}
}

func (a *Adapter) inbound_listCustomersWrapper(handler func(ctx context.Context, query CustomerQuery) (*CustomerPage, error)) RequestResponseHandler {
	return func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink) {
		var query CustomerQuery
		a.invokeRR(request, sink, &query,
			func() (interface{}, error) {
				return handler(ctx, query)
			},
		)
	}
}

func (a Adapter) invokeRR(request payload.Payload, sink mono.Sink, input interface{}, fn func() (interface{}, error)) {
	if err := a.codec.Decode(request.Data(), input); err != nil {
		sink.Error(err)
		return
	}

	response, err := fn()
	if err != nil {
		sink.Error(err)
		return
	}

	responseBytes, err := a.codec.Encode(response)
	if err != nil {
		sink.Error(err)
		return
	}

	sink.Success(payload.New(responseBytes, nil))
}

func (a Adapter) requestPayload(path string, data interface{}) (payload.Payload, error) {
	metadataBytes, err := a.codec.Encode(metadata.MD{
		":path": []string{path},
	})
	if err != nil {
		return nil, err
	}
	var requestBytes []byte
	if data != nil {
		requestBytes, err = a.codec.Encode(data)
		if err != nil {
			return nil, err
		}
	}

	return payload.New(requestBytes, metadataBytes), nil
}

type OutboundImpl struct {
	a *Adapter
}

func NewOutboundImpl(a *Adapter) *OutboundImpl {
	return &OutboundImpl{
		a: a,
	}
}

func (p *OutboundImpl) SaveCustomer(ctx context.Context, customer Customer) error {
	request, err := p.requestPayload("/customers.v1.Outbound/getCustomers", customer)
	if err != nil {
		return err
	}

	_, err = p.a.client.RequestResponse(request).Block(ctx)
	return err
}

func (p *OutboundImpl) FetchCustomer(ctx context.Context, id uint64) (*Customer, error) {
	inputArgs := outboundFetchCustomerArgs{
		ID: id,
	}
	var result Customer
	if err := p.requestResponse(ctx, "/customers.v1.Outbound/fetchCustomer", &inputArgs, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *OutboundImpl) CustomerCreated(ctx context.Context, customer Customer) error {
	return p.handleFF(ctx, "/customers.v1.Outbound/customerCreated", &customer)
}

func (p *OutboundImpl) GetCustomers(ctx context.Context) (CustomerPublisher, error) {
	request, err := p.requestPayload("/customers.v1.Outbound/getCustomers", nil)
	if err != nil {
		return nil, err
	}

	f := p.a.client.RequestStream(request)
	source := newCustomerPublisher(ctx, p.a.codec)
	source.subscribe(ctx, f)

	return source, nil
}

func (p *OutboundImpl) TransformCustomers(ctx context.Context, prefix string, source CustomerSource) (CustomerPublisher, error) {
	request, err := p.requestPayload("/customers.v1.Outbound/transformCustomers", &transformCustomersArgs{
		Prefix: prefix,
	})
	if err != nil {
		return nil, err
	}

	in := flux.Create(func(ctx context.Context, emitter flux.Sink) {
		emitter.Next(request)
		source(newCustomerSubscriber(p.a.codec, emitter))
	})

	f := p.a.client.RequestChannel(in)
	publisher := newCustomerPublisher(ctx, p.a.codec)
	publisher.subscribe(ctx, f)

	return publisher, nil
}

type transformCustomersArgs struct {
	Prefix string `json:"prefix" msgpack:"prefix"`
}

func (p *OutboundImpl) handleFF(ctx context.Context, path string, input interface{}) error {
	request, err := p.requestPayload(path, input)
	if err != nil {
		return err
	}

	_, err = p.a.client.RequestResponse(request).Block(ctx)
	return err
}

func (p *OutboundImpl) requestResponse(ctx context.Context, path string, input interface{}, dst interface{}) error {
	request, err := p.requestPayload(path, input)
	if err != nil {
		return err
	}

	resp, err := p.a.client.RequestResponse(request).Block(ctx)
	if err != nil {
		return err
	}

	return p.a.codec.Decode(resp.Data(), dst)
}

func (p *OutboundImpl) requestPayload(path string, data interface{}) (payload.Payload, error) {
	metadataBytes, err := p.a.codec.Encode(metadata.MD{
		":path": []string{path},
	})
	if err != nil {
		return nil, err
	}
	var requestBytes []byte
	if data != nil {
		requestBytes, err = p.a.codec.Encode(data)
		if err != nil {
			return nil, err
		}
	}

	return payload.New(requestBytes, metadataBytes), nil
}

type customerPublisher struct {
	ctx   context.Context
	codec functions.Codec
	c     chan *Customer
	e     chan error
}

func newCustomerPublisher(ctx context.Context, codec functions.Codec) *customerPublisher {
	return &customerPublisher{
		ctx:   ctx,
		codec: codec,
		c:     make(chan *Customer, 100),
		e:     make(chan error, 1),
	}
}

func (s *customerPublisher) Receive() (customer *Customer, err error) {
	var ok bool
	select {
	case customer, ok = <-s.c:
	case err, ok = <-s.e:
	case <-s.ctx.Done():
		return nil, context.Canceled
	}
	if !ok {
		return nil, io.EOF
	}
	return customer, err
}

func (s *customerPublisher) subscribe(ctx context.Context, f flux.Flux) {
	f.DoOnNext(func(input payload.Payload) error {
		var customer Customer
		if err := s.codec.Decode(input.Data(), &customer); err != nil {
			return err
		}
		s.c <- &customer
		return nil
	}).DoOnError(func(err error) {
		s.e <- err
	}).DoOnComplete(func() {
		close(s.c)
		close(s.e)
	}).Subscribe(ctx)
}

type customerSubscriber struct {
	codec functions.Codec
	sink  flux.Sink
}

func newCustomerSubscriber(codec functions.Codec, sink flux.Sink) *customerSubscriber {
	return &customerSubscriber{
		codec: codec,
		sink:  sink,
	}
}

func (s *customerSubscriber) Next(customer *Customer) {
	responseBytes, err := s.codec.Encode(customer)
	if err != nil {
		s.sink.Error(err)
		return
	}
	s.sink.Next(payload.New(responseBytes, nil))
}

func (s *customerSubscriber) Complete() {
	s.sink.Complete()
}

func (s *customerSubscriber) Error(err error) {
	s.sink.Error(err)
}

type inboundGetCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}

type outboundFetchCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}

type AdapterContext struct {
	*stateful.Context
}

func (c *AdapterContext) Self() LogicalAddress {
	self := c.Context.Self()
	return LogicalAddress{
		Type: self.Type,
		ID:   self.ID,
	}
}

type Storage struct {
	a *Adapter
}

func NewStorage(a *Adapter) *Storage {
	return &Storage{
		a: a,
	}
}

func (s *Storage) Get(namespace, id, key string) (stateful.RawItem, bool, error) {
	ctx := context.Background()
	var item stateful.RawItem

	type Args struct {
		Namespace string `json:"namespace" msgpack:"namespace"`
		ID        string `json:"id" msgpack:"id"`
		Key       string `json:"key" msgpack:"key"`
	}

	request, err := s.a.requestPayload("nanobus:state/get", &Args{
		Namespace: namespace,
		ID:        id,
		Key:       key,
	})
	if err != nil {
		return item, false, err
	}

	resp, err := s.a.client.RequestResponse(request).Block(ctx)
	if err != nil {
		return item, false, err
	}

	if err := s.a.codec.Decode(resp.Data(), &item); err != nil {
		return item, false, err
	}

	return item, true, nil
}

func (a *Adapter) RegisterCustomerActor(stateful CustomerActor) *Adapter {
	a.RegisterRR("/customers.v1.CustomerActor/deactivate", a.deactivateHandler("customers.v1.CustomerActor", stateful))
	a.RegisterRR("/customers.v1.CustomerActor/createCustomer", a.customerActor_createCustomerWrapper(stateful))
	a.RegisterRR("/customers.v1.CustomerActor/getCustomer", a.customerActor_getCustomerWrapper(stateful))
	return a
}

func (a *Adapter) deactivateHandler(actorType string, actor interface{}) RequestResponseHandler {
	return func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink) {
		id, ok := md.Scalar(":id")
		if !ok {
			sink.Error(ErrNotFound)
			return
		}

		sctx, err := a.stateManager.ToContext(ctx, actorType, id, actor)
		if err != nil {
			sink.Error(ErrNotFound)
			return
		}
		if deactivator, ok := actor.(stateful.Deactivator); ok {
			deactivator.Deactivate(sctx)
		}
		a.stateManager.Deactivate(stateful.LogicalAddress{
			Type: actorType,
			ID:   id,
		})

		sink.Success(payload.Empty())
	}
}

func (a *Adapter) customerActor_createCustomerWrapper(stateful CustomerActor) RequestResponseHandler {
	return func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink) {
		var input Customer
		id, ok := md.Scalar(":id")
		if !ok {
			sink.Error(ErrNotFound)
			return
		}

		if err := a.codec.Decode(request.Data(), &input); err != nil {
			sink.Error(err)
			return
		}

		sctx, err := a.stateManager.ToContext(ctx, "customers.v1.CustomerActor", id, stateful)
		if err != nil {
			sink.Error(err)
			return
		}
		result, err := stateful.CreateCustomer(&AdapterContext{&sctx}, input)
		if err != nil {
			sink.Error(err)
			return
		}

		response, err := sctx.Response(result)
		if err != nil {
			sink.Error(err)
			return
		}

		responseBytes, err := a.codec.Encode(response)
		if err != nil {
			sink.Error(err)
			return
		}

		sink.Success(payload.New(responseBytes, nil))
	}
}

func (a *Adapter) customerActor_getCustomerWrapper(stateful CustomerActor) RequestResponseHandler {
	return func(ctx context.Context, md metadata.MD, request payload.Payload, sink mono.Sink) {
		id, ok := md.Scalar(":id")
		if !ok {
			sink.Error(ErrNotFound)
			return
		}

		sctx, err := a.stateManager.ToContext(ctx, "customers.v1.CustomerActor", id, stateful)
		if err != nil {
			sink.Error(err)
			return
		}
		response, err := stateful.GetCustomer(&AdapterContext{&sctx})
		if err != nil {
			sink.Error(err)
			return
		}

		resp, err := sctx.Response(response)
		if err != nil {
			sink.Error(err)
			return
		}

		responseBytes, err := a.codec.Encode(resp)
		if err != nil {
			sink.Error(err)
			return
		}

		sink.Success(payload.New(responseBytes, nil))
	}
}
