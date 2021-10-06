package main

import (
	"github.com/nanobus/nanobus/example/welcome/pkg/welcome"
)

func main() {
	adapter := welcome.NewAdapter()
	outbound := adapter.NewOutbound()
	service := welcome.NewService(outbound)

	adapter.RegisterInbound(welcome.Inbound{
		GreetCustomer: service.GreetCustomer,
	})

	adapter.Run()
}
