package customers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

type MyApp struct {
	customers.Outbound
}

func New(outbound customers.Outbound) *MyApp {
	return &MyApp{
		Outbound: outbound,
	}
}

func (s *MyApp) CreateCustomer(ctx context.Context, customer customers.Customer) (customers.Customer, error) {
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("RECEIVED: %s\n", string(jsonBytes))
	}

	err := s.SaveCustomer(ctx, customer)
	if err != nil {
		return customer, err
	}
	err = s.CustomerCreated(ctx, customer)

	return customer, err
}

func (s *MyApp) GetCustomer(ctx context.Context, id uint64) (customers.Customer, error) {
	log.Printf("RECEIVED: %d\n", id)

	return s.FetchCustomer(ctx, id)
}
