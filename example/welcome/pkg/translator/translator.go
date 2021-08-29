package translator

import (
	"context"

	functions "github.com/nanobus/go-functions"
	"github.com/nanobus/nanobus/example/welcome/pkg/welcome"
)

type Functions struct {
	invoker *functions.Invoker
}

func New(invoker *functions.Invoker) *Functions {
	return &Functions{
		invoker: invoker,
	}
}

func (m *Functions) SendEmail(ctx context.Context, email string, message string) error {
	inputArgs := hostSendEmailArgs{
		Email:   email,
		Message: message,
	}
	return m.invoker.Invoke(ctx, "welcome.v1.Outbound", "sendEmail", inputArgs)
}

type InboudHandlers struct {
	GreetCustomer func(ctx context.Context, customer welcome.Customer) error
}

func (h InboudHandlers) Register(codec functions.Codec, registerFn functions.Register) {
	if h.GreetCustomer != nil {
		registerFn("welcome.v1.Inbound", "greetCustomer", greetCustomerWrapper(codec, h.GreetCustomer))
	}
}

func greetCustomerWrapper(codec functions.Codec, handler func(ctx context.Context, customer welcome.Customer) error) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var request welcome.Customer
		if err := codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		err := handler(ctx, request)
		if err != nil {
			return nil, err
		}
		return []byte{}, nil
	}
}

type hostSendEmailArgs struct {
	Email   string `json:"email" msgpack:"email"`
	Message string `json:"message" msgpack:"message"`
}
