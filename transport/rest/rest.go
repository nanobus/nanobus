package rest

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-logr/logr"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nanobus/go-functions"

	"github.com/nanobus/nanobus/errorz"
	"github.com/nanobus/nanobus/spec"
	"github.com/nanobus/nanobus/transport"
	"github.com/nanobus/nanobus/transport/filter"
)

type Rest struct {
	log           logr.Logger
	address       string
	namespaces    spec.Namespaces
	invoker       transport.Invoker
	errorResolver errorz.Resolver
	codecs        map[string]functions.Codec
	filters       []filter.Filter
	router        *mux.Router
	ln            net.Listener
}

type queryParam struct {
	name         string
	arg          string
	isArray      bool
	required     bool
	typeRef      *spec.TypeRef
	defaultValue interface{}
}

type optionsHolder struct {
	codecs  []functions.Codec
	filters []filter.Filter
}

var rePathParams = regexp.MustCompile(`(?m)\{([^\}]*)\}`)

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
// 	return "rest", New
// }

func New(log logr.Logger, address string, namespaces spec.Namespaces, invoker transport.Invoker, errorResolver errorz.Resolver, options ...Option) (transport.Transport, error) {
	var opts optionsHolder

	for _, opt := range options {
		opt(&opts)
	}

	codecMap := make(map[string]functions.Codec, len(opts.codecs))
	for _, c := range opts.codecs {
		codecMap[c.ContentType()] = c
	}

	r := mux.NewRouter()
	r.Use(handlers.ProxyHeaders)
	r.Use(mux.CORSMethodMiddleware(r))

	swaggerHost := address
	if strings.HasPrefix(swaggerHost, ":") {
		swaggerHost = "localhost" + swaggerHost
	}
	log.Info("Registering Swagger UI", "url", fmt.Sprintf("http://%s/swagger/", swaggerHost))
	if err := RegisterSwaggerRoutes(r, namespaces); err != nil {
		return nil, err
	}

	rest := Rest{
		log:           log,
		address:       address,
		namespaces:    namespaces,
		invoker:       invoker,
		errorResolver: errorResolver,
		codecs:        codecMap,
		filters:       opts.filters,
		router:        r,
	}

	for _, namespace := range namespaces {
		pathNS := ""
		if path, ok := namespace.Annotation("path"); ok {
			if arg, ok := path.Argument("value"); ok {
				pathNS = arg.ValueString()
			}
		}

		for _, service := range namespace.Services {
			_, isService := service.Annotation("service")
			_, isActor := service.Annotation("actor")
			_, isStateful := service.Annotation("stateful")

			isActor = isActor || isStateful
			if !(isService || isActor) {
				continue
			}

			pathSrv := ""
			if path, ok := service.Annotation("path"); ok {
				if arg, ok := path.Argument("value"); ok {
					pathSrv = arg.ValueString()
				}
			}

			for _, operation := range service.Operations {
				pathOper := ""
				if path, ok := operation.Annotation("path"); ok {
					if arg, ok := path.Argument("value"); ok {
						pathOper = arg.ValueString()
					}
				}

				path := pathNS + pathSrv + pathOper
				if path == "" {
					continue
				}

				method := ""
				if _, ok := operation.Annotation("GET"); ok {
					method = "GET"
				} else if _, ok := operation.Annotation("POST"); ok {
					method = "POST"
				} else if _, ok := operation.Annotation("PUT"); ok {
					method = "PUT"
				} else if _, ok := operation.Annotation("PATCH"); ok {
					method = "PATCH"
				} else if _, ok := operation.Annotation("DELETE"); ok {
					method = "DELETE"
				} else {
					continue
				}

				bodyParamName := ""
				hasBody := false
				queryParams := map[string]queryParam{}

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
						} else if _, ok := param.Annotation("query"); ok {
							t := param.Type
							required := true
							isArray := false
							if t.OptionalType != nil {
								required = false
								t = t.OptionalType
							}
							if t.ListType != nil {
								t = t.ListType
								isArray = true
							}
							if t.IsPrimitive() {
								queryParams[param.Name] = queryParam{
									name:         param.Name,
									required:     required,
									isArray:      isArray,
									typeRef:      t,
									defaultValue: param.DefaultValue,
								}
							} else if t.Type != nil {
								for _, f := range param.Type.Type.Fields {
									t := param.Type
									required := true
									isArray := false
									if t.OptionalType != nil {
										t = t.OptionalType
									}
									if t.ListType != nil {
										t = t.ListType
										isArray = true
									}

									queryParams[f.Name] = queryParam{
										name:         f.Name,
										arg:          param.Name,
										required:     required,
										isArray:      isArray,
										typeRef:      t,
										defaultValue: f.DefaultValue,
									}
								}
							}
						} else {
							bodyParamName = param.Name
							hasBody = true
						}
					}
				} else {
					_, hasQuery := operation.Parameters.Annotation("query")
					if hasQuery {
						for _, param := range operation.Parameters.Fields {
							if param.Type.IsPrimitive() {
								queryParams[param.Name] = queryParam{
									name:         param.Name,
									isArray:      false, // TODO
									typeRef:      param.Type,
									defaultValue: param.DefaultValue,
								}
							} else {
								for _, f := range param.Type.Type.Fields {
									queryParams[f.Name] = queryParam{
										name:         f.Name,
										isArray:      false, // TODO
										typeRef:      f.Type,
										defaultValue: f.DefaultValue,
									}
								}
							}
						}
					} else {
						hasBody = true
					}
				}

				log.Info("Registering REST handler", "method", method, "path", path)
				r.HandleFunc(path, rest.handler(
					namespace.Name, service.Name, operation.Name, isActor,
					hasBody, bodyParamName, queryParams)).Methods(method)
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

func (t *Rest) handler(namespace, service, operation string, isActor bool,
	hasBody bool, bodyParamName string, queryParams map[string]queryParam) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()
		vars := mux.Vars(r)
		id := ""
		if isActor {
			id = vars["id"]
		}

		contentType := r.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/json"
		}

		codec, ok := t.codecs[contentType]
		if !ok {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			fmt.Fprintf(w, "%v: %s", ErrUnregisteredContentType, contentType)
			return
		}

		for _, filter := range t.filters {
			var err error
			if ctx, err = filter(ctx, r.Header); err != nil {
				t.handleError(err, codec, r, w, http.StatusInternalServerError)
				return
			}
		}

		requestBytes, err := io.ReadAll(r.Body)
		if err != nil {
			t.handleError(err, codec, r, w, http.StatusInternalServerError)
			return
		}

		var input map[string]interface{}
		if hasBody && len(requestBytes) > 0 {
			if bodyParamName == "" {
				if err := codec.Decode(requestBytes, &input); err != nil {
					t.handleError(err, codec, r, w, http.StatusInternalServerError)
					return
				}
			} else {
				var body interface{}
				if err := codec.Decode(requestBytes, &body); err != nil {
					t.handleError(err, codec, r, w, http.StatusInternalServerError)
					return
				}
				input = map[string]interface{}{
					bodyParamName: body,
				}
			}
		} else {
			input = make(map[string]interface{}, len(vars)+len(queryParams))
		}

		for name, value := range vars {
			input[name] = value
		}

		if len(queryParams) > 0 {
			queryValues, _ := url.ParseQuery(r.URL.RawQuery)
			for name, q := range queryParams {
				if values, ok := queryValues[name]; ok {
					var converted interface{}
					if q.isArray {
						items := make([]interface{}, 0, 100)
						for _, value := range values {
							parts := strings.Split(value, ",")
							for _, v := range parts {
								converted, _, err = q.typeRef.Coalesce(v, false)
								if err != nil {
									t.handleError(err, codec, r, w, http.StatusBadRequest)
									return
								}
								items = append(items, converted)
							}
						}
						converted = items
					} else {
						converted, _, err = q.typeRef.Coalesce(values[0], false)
						if err != nil {
							t.handleError(err, codec, r, w, http.StatusBadRequest)
							return
						}
					}
					wrapper := input
					if q.arg != "" {
						var w interface{}
						found := false
						if w, found = input[q.arg]; found {
							wrapper, found = w.(map[string]interface{})
						}
						if !found {
							wrapper = make(map[string]interface{}, len(queryValues))
							input[q.arg] = wrapper
						}
					}
					wrapper[name] = converted
				} else if q.isArray && q.required {
					input[name] = []interface{}{}
				} else if q.defaultValue != nil {
					input[name] = q.defaultValue
				}
			}
		}

		response, err := t.invoker(ctx, namespace, service, id, operation, input)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, transport.ErrBadInput) {
				code = http.StatusBadRequest
			}
			t.handleError(err, codec, r, w, code)
			return
		}

		w.Header().Set("Content-Type", codec.ContentType())
		responseBytes, err := codec.Encode(response)
		if err != nil {
			t.handleError(err, codec, r, w, http.StatusInternalServerError)
			return
		}

		w.Write(responseBytes)
	}
}

func (t *Rest) handleError(err error, codec functions.Codec, req *http.Request, w http.ResponseWriter, status int) {
	errz := t.errorResolver(err)
	errz.Path = req.RequestURI

	w.Header().Add("Content-Type", codec.ContentType())
	w.WriteHeader(errz.Status)
	payload, err := codec.Encode(errz)
	if err != nil {
		fmt.Fprint(w, "unknown error")
	}

	w.Write(payload)
}
