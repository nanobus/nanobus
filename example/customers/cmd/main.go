package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/oklog/run"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/transports/mux"
	customers "github.com/nanobus/nanobus/example/customers"
	"github.com/nanobus/nanobus/example/customers/pkg/translator"
)

func main() {
	host := LookupEnvOrString("HOST", "localhost")
	port := LookupEnvOrInt("PORT", 9000)
	outboundBaseURI := LookupEnvOrString("OUTBOUND_BASE_URI", "http://localhost:9000/outbound/")

	inCodec := msgpack.New()
	outCodec := msgpack.New()
	m := mux.New(outboundBaseURI, outCodec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, outCodec)
	s := customers.New(translator.New(invoker))

	translator.Handlers{
		CreateCustomer: s.CreateCustomer,
		GetCustomer:    s.GetCustomer,
	}.Register(inCodec, m.Register)

	ctx := context.Background()
	var g run.Group
	{
		httpListenAddr := fmt.Sprintf("%s:%d", host, port)
		ln, err := net.Listen("tcp", httpListenAddr)
		if err != nil {
			log.Fatalln(err)
		}
		g.Add(func() error {
			return http.Serve(ln, m.Router())
		}, func(error) {
			ln.Close()
		})
	}
	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}

	if err := g.Run(); err.Error() != "received signal interrupt" {
		log.Fatalln(err)
	}
}

func LookupEnvOrInt(key string, defaultVal int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
