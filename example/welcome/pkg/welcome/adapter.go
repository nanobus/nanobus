package welcome

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

	functions "github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/transports/mux"
)

var busURI = lookupEnvOrString("BUS_URI", "http://localhost:32321")

type OutboundImpl struct {
	invoker *functions.Invoker
}

func NewOutboundImpl(invoker *functions.Invoker) *OutboundImpl {
	return &OutboundImpl{
		invoker: invoker,
	}
}

func (m *OutboundImpl) SendEmail(ctx context.Context, email string, message string) error {
	inputArgs := outboundSendEmailArgs{
		Email:   email,
		Message: message,
	}
	return m.invoker.Invoke(ctx, "welcome.v1.Outbound", "sendEmail", inputArgs)
}

func lookupEnvOrInt(key string, defaultVal int) int {
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

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

type Adapter struct {
	mux                *mux.Mux
	codec              functions.Codec
	invoker            *functions.Invoker
	registerFn         functions.Register
	registerStatefulFn functions.RegisterStateful

	ln net.Listener
}

func NewAdapter() *Adapter {
	codec := msgpack.New()
	m := mux.New(busURI+"/outbound/", codec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, codec)

	app := Adapter{
		mux:                m,
		codec:              codec,
		invoker:            invoker,
		registerFn:         m.Register,
		registerStatefulFn: m.RegisterStateful,
	}

	return &app
}

func (a *Adapter) Start() (err error) {
	host := lookupEnvOrString("HOST", "localhost")
	port := lookupEnvOrInt("PORT", 9000)
	httpListenAddr := fmt.Sprintf("%s:%d", host, port)
	a.ln, err = net.Listen("tcp", httpListenAddr)
	if err != nil {
		return err
	}
	return http.Serve(a.ln, a.mux.Router())
}

func (a *Adapter) Stop() (err error) {
	if a.ln != nil {
		err = a.ln.Close()
	}
	return err
}

func (a *Adapter) Run() {
	ctx := context.Background()
	var g run.Group
	{
		g.Add(func() error {
			return a.Start()
		}, func(error) {
			a.Stop()
		})
	}
	{
		g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	}

	if err := g.Run(); err.Error() != "received signal interrupt" {
		log.Fatalln(err)
	}
}

func (a *Adapter) Invoker() *functions.Invoker {
	return a.invoker
}

func (a *Adapter) RegisterInbound(handlers Inbound) *Adapter {
	if handlers.GreetCustomer != nil {
		a.registerFn("welcome.v1.Inbound", "greetCustomer", a.inbound_greetCustomerWrapper(handlers.GreetCustomer))
	}
	return a
}

func (a *Adapter) inbound_greetCustomerWrapper(handler func(ctx context.Context, customer Customer) error) functions.Handler {
	return func(ctx context.Context, payload []byte) ([]byte, error) {
		var request Customer
		if err := a.codec.Decode(payload, &request); err != nil {
			return nil, err
		}
		err := handler(ctx, request)
		if err != nil {
			return nil, err
		}
		return []byte{}, nil
	}
}

func (a *Adapter) NewOutbound() Outbound {
	return NewOutboundImpl(a.invoker)
}

type outboundSendEmailArgs struct {
	Email   string `json:"email" msgpack:"email"`
	Message string `json:"message" msgpack:"message"`
}
