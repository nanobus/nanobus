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
)

type HTTPRPC struct {
	address string
	invoker transport.Invoker
	codecs  map[string]functions.Codec
	ln      net.Listener
}

var (
	ErrUnregisteredContentType = errors.New("unregistered content type")
	ErrInvalidURISyntax        = errors.New("invalid invocation URI syntax")
)

func Loader() (string, transport.Loader) {
	return "httprpc", New
}

func New(address string, namespaces spec.Namespaces, invoker transport.Invoker, codecs ...functions.Codec) (transport.Transport, error) {
	codecMap := make(map[string]functions.Codec, len(codecs))

	for _, c := range codecs {
		codecMap[c.ContentType()] = c
	}

	return &HTTPRPC{
		address: address,
		invoker: invoker,
		codecs:  codecMap,
	}, nil
}

func (t *HTTPRPC) Listen() error {
	r := mux.NewRouter()
	r.HandleFunc("/{namespace}/{function}", t.handler).Methods("POST")
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
	if err := codec.Decode(requestBytes, &input); err != nil {
		handleError(err, w, http.StatusInternalServerError)
		return
	}

	response, err := t.invoker(r.Context(), namespace, service, function, input)
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
