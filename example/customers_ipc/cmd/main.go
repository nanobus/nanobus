package main

import (
	"context"
	"log"

	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

func main() {
	adapter := customers.NewAdapter()
	outbound := adapter.NewOutbound()
	service := customers.NewService(outbound)

	adapter.RegisterInbound(service)
	adapter.RegisterCustomerActor(customers.NewCustomerActorImpl())

	log.Fatal(adapter.Start(context.Background()))
}
