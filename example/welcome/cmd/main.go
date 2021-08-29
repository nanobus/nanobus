package main

import (
	"github.com/nanobus/nanobus/example/welcome/pkg/welcome"
)

func main() {
	app := welcome.NewApplication()
	outbound := app.NewOutbound()
	service := welcome.New(outbound)

	app.RegisterInbound(welcome.Inboud{
		GreetCustomer: service.GreetCustomer,
	}).Run()
}
