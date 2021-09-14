package customers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nanobus/go-functions/stateful"
)

type Service struct {
	outbound Outbound
}

func NewService(outbound Outbound) *Service {
	return &Service{
		outbound: outbound,
	}
}

func (s *Service) CreateCustomer(ctx context.Context, customer Customer) (*Customer, error) {
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("RECEIVED: %s\n", string(jsonBytes))
	}

	err := s.outbound.SaveCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}
	err = s.outbound.CustomerCreated(ctx, customer)

	return &customer, err
}

func (s *Service) GetCustomer(ctx context.Context, id uint64) (*Customer, error) {
	log.Printf("RECEIVED: %d\n", id)

	return s.outbound.FetchCustomer(ctx, id)
}

type CustomerActorImpl struct{}

func NewCustomerActorImpl() *CustomerActorImpl {
	return &CustomerActorImpl{}
}

func (c *CustomerActorImpl) CreateCustomer(ctx stateful.Context, customer Customer) (*Customer, error) {
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("ACTOR RECEIVED: %s\n", string(jsonBytes))
	}

	log.Printf("Actor Type/ID = %s", &ctx.Self)

	ctx.Set("customer", &customer)

	return &customer, nil
}

func (c *CustomerActorImpl) GetCustomer(ctx stateful.Context) (*Customer, error) {
	log.Printf("RECEIVED\n")

	var customer Customer
	if _, err := ctx.Get("customer", &customer); err != nil {
		return nil, err
	}
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("ACTOR RETURNING: %s\n", string(jsonBytes))
	}

	return &customer, nil
}
