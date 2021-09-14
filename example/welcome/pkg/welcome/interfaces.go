package welcome

import (
	"context"
)

type ns struct{}

func (n *ns) Namespace() string {
	return "welcome.v1"
}

func (n *ns) Version() string {
	return "0.1.0"
}

type Inbound struct {
	GreetCustomer func(ctx context.Context, customer Customer) error
}

type Outbound interface {
	SendEmail(ctx context.Context, email string, message string) error
}

type Customer struct {
	ns
	FirstName string `json:"firstName" msgpack:"firstName"`
	LastName  string `json:"lastName" msgpack:"lastName"`
	Email     string `json:"email" msgpack:"email"`
}

func (c *Customer) Type() string {
	return "Customer"
}
