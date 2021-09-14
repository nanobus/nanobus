package welcome

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/oklog/run"

	functions "github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/codecs/msgpack"
	"github.com/nanobus/go-functions/stateful"
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

type Storage struct {
	codec functions.Codec
}

func NewStorage(codec functions.Codec) *Storage {
	return &Storage{
		codec: codec,
	}
}

func (s *Storage) Get(namespace, id, key string) (stateful.RawItem, bool, error) {
	var item stateful.RawItem
	url := busURI + "/state/" + namespace + "/" + id + "/" + key
	resp, err := http.Get(url)
	if err != nil {
		return item, false, err
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return item, false, err
	}

	if len(payload) == 0 {
		return item, false, nil
	}
	if err = s.codec.Decode(payload, &item); err != nil {
		return item, false, err
	}

	return item, true, nil
}

type Adapter struct {
	mux                *mux.Mux
	stateManager       *stateful.Manager
	codec              functions.Codec
	invoker            *functions.Invoker
	registerFn         functions.Register
	registerStatefulFn functions.RegisterStateful

	ln net.Listener
}

func NewAdapter(stateManager *stateful.Manager) *Adapter {
	codec := msgpack.New()
	m := mux.New(busURI+"/outbound/", codec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, codec)

	app := Adapter{
		mux:                m,
		stateManager:       stateManager,
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
