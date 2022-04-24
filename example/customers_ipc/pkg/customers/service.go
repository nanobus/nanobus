/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package customers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/go-logr/logr"

	"github.com/nanobus/adapter-go/stateful"
)

type Service struct {
	log      logr.Logger
	outbound Outbound
}

var _ = (Inbound)((*Service)(nil))

func NewService(log logr.Logger, outbound Outbound) *Service {
	return &Service{
		log:      log,
		outbound: outbound,
	}
}

func (s *Service) CreateCustomer(ctx context.Context, customer *Customer) (*Customer, error) {
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("RECEIVED: %s\n", string(jsonBytes))
	}

	err := s.outbound.SaveCustomer(ctx, customer)
	if err != nil {
		return nil, err
	}
	//err = s.outbound.CustomerCreated(ctx, customer)

	return customer, err
}

func (s *Service) GetCustomer(ctx context.Context, id uint64) (*Customer, error) {
	log.Printf("RECEIVED: %d\n", id)

	stream, err := s.outbound.GetCustomers(ctx)
	if err != nil {
		return nil, err
	}
	for {
		customer, err := stream.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return nil, err
		}

		jsonBytes, err := json.MarshalIndent(customer, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(jsonBytes))
	}

	stream, err = s.outbound.TransformCustomers(ctx, "TEST_", func(sub CustomerSubscriber) {
		sub.Next(&Customer{
			ID:        1234,
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@email.com",
		})
		sub.Complete()
	})
	if err != nil {
		return nil, err
	}
	for {
		customer, err := stream.Receive()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			return nil, err
		}

		jsonBytes, err := json.MarshalIndent(customer, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(jsonBytes))
	}

	return s.outbound.FetchCustomer(ctx, id)
}

func (s *Service) ListCustomers(ctx context.Context, query *CustomerQuery) (*CustomerPage, error) {
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

func NewCustomerActor() *CustomerActorImpl {
	return &CustomerActorImpl{}
}

func (c *CustomerActorImpl) Activate(ctx stateful.Context) error {
	log.Printf("Activated %s", ctx.Self())

	return nil
}

func (c *CustomerActorImpl) Deactivate(ctx stateful.Context) error {
	log.Printf("Deactivated %s", ctx.Self())

	return nil
}

func (c *CustomerActorImpl) CreateCustomer(ctx stateful.Context, customer *Customer) (*Customer, error) {
	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("ACTOR RECEIVED: %s\n", string(jsonBytes))
	}

	log.Printf("Actor Type/ID = %s", ctx.Self())

	ctx.Set("customer", &customer)

	return customer, nil
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
