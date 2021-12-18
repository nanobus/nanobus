package customers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	//err = s.outbound.CustomerCreated(ctx, customer)

	return &customer, err
}

func (s *Service) GetCustomer(ctx context.Context, id uint64) (*Customer, error) {
	log.Printf("RECEIVED: %d\n", id)

	s.outbound.GetCustomers(ctx, func(recv RecvCustomer) error {
		for {
			var customer Customer
			if err := recv(&customer); err != nil {
				if err == io.EOF {
					return nil
				}
				fmt.Println(err)
				return err
			}

			jsonBytes, _ := json.MarshalIndent(&customer, "", "  ")
			fmt.Println(string(jsonBytes))
		}
	})

	// return &Customer{
	// 	ID:        id,
	// 	FirstName: "TEST",
	// }, nil

	return s.outbound.FetchCustomer(ctx, id)
}

func (s *Service) ListCustomers(ctx context.Context, query CustomerQuery) (*CustomerPage, error) {
	if jsonBytes, err := json.MarshalIndent(&query, "", "  "); err == nil {
		log.Printf("RECEIVED: %s\n", string(jsonBytes))
	}

	return &CustomerPage{
		Offset: query.Offset,
		Limit:  query.Limit,
		Items:  []Customer{},
	}, nil
}

type CustomerActorImpl struct{}

func NewCustomerActorImpl() *CustomerActorImpl {
	return &CustomerActorImpl{}
}

func (c *CustomerActorImpl) Activate(ctx Context) error {
	log.Printf("Activated %s", ctx.Self())

	return nil
}

func (c *CustomerActorImpl) Deactivate(ctx Context) error {
	log.Printf("Deactivated %s", ctx.Self())

	return nil
}

func (c *CustomerActorImpl) CreateCustomer(ctx Context, customer Customer) (*Customer, error) {
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("ACTOR RECEIVED: %s\n", string(jsonBytes))
	}

	log.Printf("Actor Type/ID = %s", ctx.Self())

	ctx.Set("customer", &customer)

	return &customer, nil
}

func (c *CustomerActorImpl) GetCustomer(ctx Context) (*Customer, error) {
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
