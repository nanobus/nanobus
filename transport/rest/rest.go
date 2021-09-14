package rest

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/nanobus/go-functions"

	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/transport"
)

type Rest struct {
	address    string
	namespaces spec.Namespaces
	invoker    transport.Invoker
	codecs     map[string]functions.Codec
	router     *mux.Router
	ln         net.Listener
}

var rePathParams = regexp.MustCompile(`(?m)\{([^\}]*)\}`)

var (
	ErrUnregisteredContentType = errors.New("unregistered content type")
	ErrInvalidURISyntax        = errors.New("invalid invocation URI syntax")
)

func Loader() (string, transport.Loader) {
	return "rest", New
}

func New(address string, namespaces spec.Namespaces, invoker transport.Invoker, codecs ...functions.Codec) (transport.Transport, error) {
	codecMap := make(map[string]functions.Codec, len(codecs))

	for _, c := range codecs {
		codecMap[c.ContentType()] = c
	}

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	rest := Rest{
		address:    address,
		namespaces: namespaces,
		invoker:    invoker,
		codecs:     codecMap,
		router:     r,
	}

	for _, namespace := range namespaces {
		pathNS := ""
		if path, ok := namespace.Annotations["path"]; ok {
			pathNS = path.Arguments["value"].ValueString()
		}

		for _, service := range namespace.Services {
			if _, ok := service.Annotations["service"]; !ok {
				continue
			}
			pathSrv := ""
			if path, ok := service.Annotations["path"]; ok {
				pathSrv = path.Arguments["value"].ValueString()
			}

			for _, operation := range service.Operations {
				pathOper := ""
				if path, ok := operation.Annotations["path"]; ok {
					pathOper = path.Arguments["value"].ValueString()
				}

				path := pathNS + pathSrv + pathOper
				if path == "" {
					continue
				}

				method := ""
				if _, ok := operation.Annotations["GET"]; ok {
					method = "GET"
				} else if _, ok := operation.Annotations["POST"]; ok {
					method = "POST"
				} else if _, ok := operation.Annotations["PUT"]; ok {
					method = "PUT"
				} else if _, ok := operation.Annotations["PATCH"]; ok {
					method = "PATCH"
				} else if _, ok := operation.Annotations["DELETE"]; ok {
					method = "DELETE"
				} else {
					continue
				}

				bodyParamName := ""
				hasBody := operation.Unary
				if !operation.Unary {
					pathParams := map[string]struct{}{}
					for _, match := range rePathParams.FindAllString(path, -1) {
						match = strings.TrimPrefix(match, "{")
						match = strings.TrimSuffix(match, "}")
						pathParams[match] = struct{}{}
					}

					for _, param := range operation.Parameters.Fields {
						if _, ok := pathParams[param.Name]; ok {
							continue
						}
						if _, ok := param.Annotations["query"]; ok {
							continue
						}
						bodyParamName = param.Name
						hasBody = true
					}
				}

				r.HandleFunc(path, rest.handler(namespace.Name, service.Name, operation.Name, hasBody, bodyParamName)).Methods(method)
			}
		}
	}

	return &rest, nil
}

func (t *Rest) Listen() error {
	ln, err := net.Listen("tcp", t.address)
	if err != nil {
		return err
	}
	t.ln = ln

	return http.Serve(ln, t.router)
}

func (t *Rest) Close() (err error) {
	if t.ln != nil {
		err = t.ln.Close()
		t.ln = nil
	}

	return err
}

func (t *Rest) handler(namespace, service, operation string, hasBody bool, bodyParamName string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		vars := mux.Vars(r)

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

		var input map[string]interface{}
		if hasBody {
			if bodyParamName == "" {
				if err := codec.Decode(requestBytes, &input); err != nil {
					handleError(err, w, http.StatusInternalServerError)
					return
				}
			} else {
				var body interface{}
				if err := codec.Decode(requestBytes, &body); err != nil {
					handleError(err, w, http.StatusInternalServerError)
					return
				}
				input = map[string]interface{}{
					bodyParamName: body,
				}
			}
		} else {
			input = make(map[string]interface{}, len(vars))
		}

		for name, value := range vars {
			input[name] = value
		}

		response, err := t.invoker(r.Context(), namespace, service, "", operation, input)
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
}

func handleError(err error, w http.ResponseWriter, status int) {
	log.Println(err)
	w.WriteHeader(status)
	fmt.Fprintf(w, "error: %v", err)
}
