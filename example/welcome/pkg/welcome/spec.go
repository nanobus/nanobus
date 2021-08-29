package welcome

import (
	"context"
)

type Inboud struct {
	GreetCustomer func(ctx context.Context, customer Customer) error
}

type Outbound interface {
	SendEmail(ctx context.Context, email string, message string) error
}

type Customer struct {
	FirstName string `json:"firstName" msgpack:"firstName"`
	LastName  string `json:"lastName" msgpack:"lastName"`
	Email     string `json:"email" msgpack:"email"`
}
