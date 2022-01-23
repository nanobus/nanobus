package customers

import (
	"context"
)

type ns struct{}

func (n *ns) Namespace() string {
	return "customers.v1"
}

func (n *ns) Version() string {
	return "0.1.0"
}

type LogicalAddress struct {
	Type string `json:"type,omitempty" msgpack:"type,omitempty"`
	ID   string `json:"id,omitempty" msgpack:"id,omitempty"`
}

func (a LogicalAddress) String() string {
	return a.Type + "/" + a.ID
}

type Context interface {
	context.Context
	Self() LogicalAddress
	Get(key string, dst interface{}) (bool, error)
	Set(key string, data interface{})
	Remove(key string)
}

// Operations that can be performed on a customer.
type Inbound interface {
	// Creates a new customer.
	CreateCustomer(ctx context.Context, customer Customer) (*Customer, error)
	// Retrieve a customer by id.
	GetCustomer(ctx context.Context, id uint64) (*Customer, error)
	// Return a page of customers using optional search filters.
	ListCustomers(ctx context.Context, query CustomerQuery) (*CustomerPage, error)
}

// Stateful operations that can be performed on a customer.
type CustomerActor interface {
	// Creates the customer state.
	CreateCustomer(ctx Context, customer Customer) (*Customer, error)
	// Retrieve the customer state.
	GetCustomer(ctx Context) (*Customer, error)
}

type CustomerPublisher interface {
	Receive() (*Customer, error)
}

type CustomerSubscriber interface {
	Complete()
	Error(error)
	Next(*Customer)
}

type CustomerSource func(CustomerSubscriber)

type Outbound interface {
	// Saves a customer to the backend database
	SaveCustomer(ctx context.Context, customer Customer) error
	// Fetches a customer from the backend database
	FetchCustomer(ctx context.Context, id uint64) (*Customer, error)
	// Sends a customer creation event
	CustomerCreated(ctx context.Context, customer Customer) error
	// Get customers from the database
	GetCustomers(ctx context.Context) (CustomerPublisher, error)
	// Transform customers
	TransformCustomers(ctx context.Context, prefix string, source CustomerSource) (CustomerPublisher, error)
}

// Customer information.
type Customer struct {
	ns
	// The customer identifer
	ID uint64 `json:"id" msgpack:"id"`
	// The customer's first name
	FirstName string `json:"firstName" msgpack:"firstName"`
	// The customer's middle name
	MiddleName *string `json:"middleName,omitempty" msgpack:"middleName,omitempty"`
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

type CustomerQuery struct {
	ns
	// The customer identifer
	ID *uint64 `json:"id,omitempty" msgpack:"id,omitempty"`
	// The customer's first name
	FirstName *string `json:"firstName,omitempty" msgpack:"firstName,omitempty"`
	// The customer's middle name
	MiddleName *string `json:"middleName,omitempty" msgpack:"middleName,omitempty"`
	// The customer's last name
	LastName *string `json:"lastName,omitempty" msgpack:"lastName,omitempty"`
	// The customer's email address
	Email  *string `json:"email,omitempty" msgpack:"email,omitempty"`
	Offset uint64  `json:"offset" msgpack:"offset"`
	Limit  uint64  `json:"limit" msgpack:"limit"`
}

func (c *CustomerQuery) Type() string {
	return "CustomerQuery"
}

type CustomerPage struct {
	ns
	Offset uint64     `json:"offset" msgpack:"offset"`
	Limit  uint64     `json:"limit" msgpack:"limit"`
	Items  []Customer `json:"items" msgpack:"items"`
}

func (c *CustomerPage) Type() string {
	return "CustomerPage"
}

// Address information.
type Address struct {
	ns
	// The address line 1
	Line1 string `json:"line1" msgpack:"line1"`
	// The address line 2
	Line2 *string `json:"line2,omitempty" msgpack:"line2,omitempty"`
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
