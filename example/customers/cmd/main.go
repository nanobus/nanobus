package main

import (
	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

func main() {
	adapter := customers.NewAdapter()
	outbound := adapter.NewOutbound()
	service := customers.NewService(outbound)

	adapter.RegisterInbound(customers.Inbound{
		CreateCustomer: service.CreateCustomer,
		GetCustomer:    service.GetCustomer,
	}).Run()
}
