package customers

import (
	"context"
)

// Operations that can be performed on a customer.
type Inbound interface {
	// Creates a new customer.
	CreateCustomer(ctx context.Context, customer Customer) (Customer, error)
	// Retrieve a customer by id.
	GetCustomer(ctx context.Context, id uint64) (Customer, error)
}

type Outbound interface {
	SaveCustomer(ctx context.Context, customer Customer) error
	FetchCustomer(ctx context.Context, id uint64) (Customer, error)
	CustomerCreated(ctx context.Context, customer Customer) error
}

// Customer information.
type Customer struct {
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

type Nested struct {
	Foo string `json:"foo" msgpack:"foo"`
	Bar string `json:"bar" msgpack:"bar"`
}

// Address information.
type Address struct {
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

// Error response.
type Error struct {
	// The detailed error message
	Message string `json:"message" msgpack:"message"`
}
