package httprpc

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nanobus/go-functions"

	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/transport"
	"github.com/nanobus/nanobus/transport/filter"
)

type HTTPRPC struct {
	address string
	invoker transport.Invoker
	codecs  map[string]functions.Codec
	filters []filter.Filter
	ln      net.Listener
}

type optionsHolder struct {
	codecs  []functions.Codec
	filters []filter.Filter
}

var (
	ErrUnregisteredContentType = errors.New("unregistered content type")
	ErrInvalidURISyntax        = errors.New("invalid invocation URI syntax")
)

type Option func(opts *optionsHolder)

func WithCodecs(codecs ...functions.Codec) Option {
	return func(opts *optionsHolder) {
		opts.codecs = codecs
	}
}

func WithFilters(filters ...filter.Filter) Option {
	return func(opts *optionsHolder) {
		opts.filters = filters
	}
}

// func Loader() (string, transport.Loader) {
// 	return "httprpc", New
// }

func New(address string, namespaces spec.Namespaces, invoker transport.Invoker, options ...Option) (transport.Transport, error) {
	var opts optionsHolder

	for _, opt := range options {
		opt(&opts)
	}

	codecMap := make(map[string]functions.Codec, len(opts.codecs))
	for _, c := range opts.codecs {
		codecMap[c.ContentType()] = c
	}

	return &HTTPRPC{
		address: address,
		invoker: invoker,
		codecs:  codecMap,
		filters: opts.filters,
	}, nil
}

func (t *HTTPRPC) Listen() error {
	r := mux.NewRouter()
	r.HandleFunc("/{namespace}/{function}", t.handler).Methods("POST")
	r.HandleFunc("/{namespace}/{id}/{function}", t.handler).Methods("POST")
	r.Use(mux.CORSMethodMiddleware(r))
	ln, err := net.Listen("tcp", t.address)
	if err != nil {
		return err
	}
	t.ln = ln

	return http.Serve(ln, r)
}

func (t *HTTPRPC) Close() (err error) {
	if t.ln != nil {
		err = t.ln.Close()
		t.ln = nil
	}

	return err
}

func (t *HTTPRPC) handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	namespace := mux.Vars(r)["namespace"]
	function := mux.Vars(r)["function"]
	id := mux.Vars(r)["id"]

	lastDot := strings.LastIndexByte(namespace, '.')
	if lastDot < 0 {
		handleError(ErrInvalidURISyntax, w, http.StatusBadRequest)
		return
	}
	service := namespace[lastDot+1:]
	namespace = namespace[:lastDot]

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	ctx := r.Context()

	for _, filter := range t.filters {
		var err error
		if ctx, err = filter(ctx, r); err != nil {
			handleError(err, w, http.StatusInternalServerError)
			return
		}
	}

	codec, ok := t.codecs[contentType]
	if !ok {
		handleError(ErrUnregisteredContentType, w, http.StatusUnsupportedMediaType)
		return
	}

	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	var input interface{}
	if len(requestBytes) > 0 {
		if err := codec.Decode(requestBytes, &input); err != nil {
			handleError(err, w, http.StatusInternalServerError)
			return
		}
	} else {
		input = map[string]interface{}{}
	}

	response, err := t.invoker(ctx, namespace, service, id, function, input)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, transport.ErrBadInput) {
			code = http.StatusBadRequest
		}
		handleError(err, w, code)
		return
	}

	w.Header().Set("Content-Type", codec.ContentType())
	responseBytes, err := codec.Encode(response)
	if err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	w.Write(responseBytes)
}

func handleError(err error, w http.ResponseWriter, status int) {
	log.Println(err)
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %v", err)
}
