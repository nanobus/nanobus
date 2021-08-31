package customers

import (
	"context"
	"encoding/json"
	"log"
)

type Service struct {
	Outbound
}

func NewService(outbound Outbound) *Service {
	return &Service{
		Outbound: outbound,
	}
}

func (s *Service) CreateCustomer(ctx context.Context, customer Customer) (Customer, error) {
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

func (s *Service) GetCustomer(ctx context.Context, id uint64) (Customer, error) {
	log.Printf("RECEIVED: %d\n", id)

	return s.FetchCustomer(ctx, id)
}
