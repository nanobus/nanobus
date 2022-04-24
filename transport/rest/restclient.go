/*
Copyright 2022 The NanoBus Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

				pathParams, raw := exampleOperationRequestBody(p, service, oper, 2)
				if len(pathParams) > 0 {
					for _, param := range pathParams {
						sb.WriteString("@")
						sb.WriteString(param)
						sb.WriteString(" = string\n")
					}
					sb.WriteString("\n")
				}

				parts := strings.Split(p, "/")
				for i, part := range parts {
					if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") &&
						!strings.HasPrefix(part, "{{") && !strings.HasSuffix(part, "}}") {
						parts[i] = "{" + part + "}"
					}
				}
				p = strings.Join(parts, "/")

				sb.WriteString(method + " {{host}}" + p + " HTTP/1.1\n")
				// TODO: query parameters
				sb.WriteString("Accept: {{accept}}\n")

				if len(raw) > 0 {
					sb.WriteString("Content-Type: {{contentType}}\n")
					sb.WriteString("\n")
					sb.WriteString(raw)
					sb.WriteString("\n")
				}

				sb.WriteString("\n")
			}
		}
	}

	return []byte(sb.String()), nil
}
