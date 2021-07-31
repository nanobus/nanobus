package translator

import (
	"context"

	functions "github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

type Functions struct {
	invoker *functions.Invoker
}

func New(invoker *functions.Invoker) *Functions {
	return &Functions{
		invoker: invoker,
	}
}

func (m *Functions) SaveCustomer(ctx context.Context, customer customers.Customer) error {
	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "saveCustomer", customer)
}

func (m *Functions) FetchCustomer(ctx context.Context, id uint64) (customers.Customer, error) {
	var ret customers.Customer
	inputArgs := hostFetchCustomerArgs{
		ID: id,
	}
	err := m.invoker.InvokeWithReturn(ctx, "customers.v1.Outbound", "fetchCustomer", inputArgs, &ret)
	return ret, err
}

func (m *Functions) CustomerCreated(ctx context.Context, customer customers.Customer) error {
	return m.invoker.Invoke(ctx, "customers.v1.Outbound", "customerCreated", customer)
}

type Handlers struct {
	// Creates a new customer.
	CreateCustomer func(ctx context.Context, customer customers.Customer) (customers.Customer, error)
	// Retrieve a customer by id.
	GetCustomer func(ctx context.Context, id uint64) (customers.Customer, error)
}

func (h Handlers) Register(codec functions.Codec, registerFn functions.Register) {
	if h.CreateCustomer != nil {
		registerFn("customers.v1.Inbound", "createCustomer", createCustomerWrapper(codec, h.CreateCustomer))
	}
	if h.GetCustomer != nil {
		registerFn("customers.v1.Inbound", "getCustomer", getCustomerWrapper(codec, h.GetCustomer))
	}
}

func createCustomerWrapper(codec functions.Codec, handler func(ctx context.Context, customer customers.Customer) (customers.Customer, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var request customers.Customer
		if err := codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		response, err := handler(ctx, request)
		if err != nil {
			return nil, err
		}
		return codec.Encode(&response)
	}
}

func getCustomerWrapper(codec functions.Codec, handler func(ctx context.Context, id uint64) (customers.Customer, error)) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var inputArgs handlerGetCustomerArgs
		if err := codec.Decode(payload, &inputArgs); err != nil {
			return nil, err
		}
		response, err := handler(ctx, inputArgs.ID)
		if err != nil {
			return nil, err
		}
		return codec.Encode(&response)
	}
}

type handlerGetCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}

type hostFetchCustomerArgs struct {
	ID uint64 `json:"id" msgpack:"id"`
}
