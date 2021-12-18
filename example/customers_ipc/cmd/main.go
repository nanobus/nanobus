package main

import (
	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

func main() {
	adapter, err := customers.NewAdapter()
	if err != nil {
		panic(err)
	}
	outbound := adapter.NewOutbound()
	service := customers.NewService(outbound)

	adapter.RegisterInbound(customers.Inbound{
		CreateCustomer: service.CreateCustomer,
		GetCustomer:    service.GetCustomer,
		ListCustomers:  service.ListCustomers,
	})
	adapter.RegisterCustomerActor(customers.NewCustomerActorImpl())

	adapter.Run()
}
