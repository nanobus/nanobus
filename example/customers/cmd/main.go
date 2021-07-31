package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/oklog/run"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/transports/mux"
	customers "github.com/nanobus/nanobus/example/customers"
	"github.com/nanobus/nanobus/example/customers/pkg/translator"
)

func main() {
	var httpListenAddr string
	var outboundBaseURI string
	flag.StringVar(
		&httpListenAddr,
		"http-listen-addr",
		LookupEnvOrString("HTTP_LISTEN_ADDR", ":8000"),
		"http listen address",
	)
	flag.StringVar(
		&outboundBaseURI,
		"outbound-base-uri",
		LookupEnvOrString("OUTBOUND_BASE_URI", "http://localhost:9000/outbound/"),
		"outbound base uri",
	)
	flag.Parse()

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

func LookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
