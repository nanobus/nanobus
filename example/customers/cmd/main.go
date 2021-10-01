package main

import (
	"github.com/nanobus/go-functions/codecs/json"
	"github.com/nanobus/go-functions/stateful"
	"github.com/nanobus/nanobus/example/customers/pkg/customers"
)

func main() {
	cache, err := stateful.NewLRUCache(200)
	if err != nil {
		panic(err)
	}
	codec := json.New()
	storage := customers.NewStorage(codec)
	manager := stateful.NewManager(cache, storage, codec)
	adapter := customers.NewAdapter(manager)
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
