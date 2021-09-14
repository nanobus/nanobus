package customers

import (
	"context"

	"github.com/nanobus/go-functions/stateful"
)

type ns struct{}

func (n *ns) Namespace() string {
	return "customers.v1"
}

func (n *ns) Version() string {
	return "0.1.0"
}

// Operations that can be performed on a customer.
type Inbound struct {
	// Creates a new customer.
	CreateCustomer func(ctx context.Context, customer Customer) (*Customer, error)
	// Retrieve a customer by id.
	GetCustomer func(ctx context.Context, id uint64) (*Customer, error)
}

type CustomerActor interface {
	// Creates the customer state.
	CreateCustomer(ctx stateful.Context, customer Customer) (*Customer, error)
	// Retrieve the customer state.
	GetCustomer(ctx stateful.Context) (*Customer, error)
}

type Outbound interface {
	SaveCustomer(ctx context.Context, customer Customer) error
	FetchCustomer(ctx context.Context, id uint64) (*Customer, error)
	CustomerCreated(ctx context.Context, customer Customer) error
}

// Customer information.
type Customer struct {
	ns
	// The customer identifer
	ID uint64 `json:"id" msgpack:"id"`
	// The customer's first name
	FirstName string `json:"firstName" msgpack:"firstName"`
	// The customer's middle name
	MiddleName *string `json:"middleName" msgpack:"middleName"`
	// The customer's last name
	LastName string `json:"lastName" msgpack:"lastName"`
	// The customer's email address
	Email string `json:"email" msgpack:"email"`
	// The customer's address
	Address Address `json:"address" msgpack:"address"`
}

func (c *Customer) Type() string {
	return "Customer"
}

type Nested struct {
	ns
	Foo string `json:"foo" msgpack:"foo"`
	Bar string `json:"bar" msgpack:"bar"`
}

func (n *Nested) Type() string {
	return "Nested"
}

// Address information.
type Address struct {
	ns
	// The address line 1
	Line1 string `json:"line1" msgpack:"line1"`
	// The address line 2
	Line2 *string `json:"line2" msgpack:"line2"`
	// The city
	City string `json:"city" msgpack:"city"`
	// The state
	State string `json:"state" msgpack:"state"`
	// The zipcode
	Zip string `json:"zip" msgpack:"zip"`
}

func (a *Address) Type() string {
	return "Address"
}

// Error response.
type Error struct {
	ns
	// The detailed error message
	Message string `json:"message" msgpack:"message"`
}

func (e *Error) Type() string {
	return "Error"
}
