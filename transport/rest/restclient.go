package rest

import (
	"bytes"
	"net/http"
	"path"
	"strings"

	"github.com/gorilla/mux"

	"github.com/nanobus/nanobus/spec"
)

func RegisterRESTClientRoutes(r *mux.Router, namespaces spec.Namespaces) error {
	specData, err := SpecToRESTClient(namespaces)
	if err != nil {
		return err
	}
	r.HandleFunc("/rest-client/service.http", func(w http.ResponseWriter, req *http.Request) {
		var v string
		if req.TLS != nil {
			v = "https://" + req.Host
		} else {
			v = "http://" + req.Host
		}
		replaced := bytes.Replace(specData, []byte("[REPLACE_HOST]"), []byte(v), 1)
		w.Write(replaced)
	})

	return nil
}

func SpecToRESTClient(namespaces spec.Namespaces) ([]byte, error) {
	var sb strings.Builder

	sb.WriteString("@host = [REPLACE_HOST]\n")
	sb.WriteString("@accept = application/json\n")
	sb.WriteString("@contentType = application/json\n")
	sb.WriteString("\n")

	for _, ns := range namespaces {
		nsPath := getAnotationString(ns, "path")
		for _, service := range ns.Services {
			_, isService := service.Annotation("service")
			_, isActor := service.Annotation("actor")
			_, isStateful := service.Annotation("stateful")
			_, isWorkflow := service.Annotation("workflow")

			isActor = isActor || isStateful || isWorkflow
			if !(isService || isActor) {
				continue
			}

			servicePath := getAnotationString(service, "path")

			for _, oper := range service.Operations {
				operPath := getAnotationString(oper, "path")
				p := path.Clean(path.Join(nsPath, servicePath, operPath))

				var method string
				if _, ok := oper.Annotation("GET"); ok {
					method = "GET"
				} else if _, ok := oper.Annotation("OPTIONS"); ok {
					method = "OPTIONS"
				} else if _, ok := oper.Annotation("HEAD"); ok {
					method = "HEAD"
				} else if _, ok := oper.Annotation("PATCH"); ok {
					method = "PATCH"
				} else if _, ok := oper.Annotation("POST"); ok {
					method = "POST"
				} else if _, ok := oper.Annotation("PUT"); ok {
					method = "PUT"
				} else if _, ok := oper.Annotation("DELETE"); ok {
					method = "DELETE"
				} else {
					continue
				}

				sb.WriteString("### " + service.Name + " - " + oper.Name + "\n")
				sb.WriteString("\n")
				sb.WriteString(method + " {{host}}" + p + " HTTP/1.1\n")
				sb.WriteString("Accept: {{accept}}\n")
				// TODO: query parameters

				if len(oper.Parameters.Fields) > 0 {
					sb.WriteString("Content-Type: {{contentType}}\n")
					sb.WriteString("\n")
					sb.WriteString(exampleOperationRequestBody(p, service, oper, 2))
					sb.WriteString("\n")
				}

				sb.WriteString("\n")
			}
		}
	}

	return []byte(sb.String()), nil
}
