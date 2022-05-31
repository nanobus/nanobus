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

package mux

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	"github.com/gorilla/mux"

	functions "github.com/nanobus/nanobus/channel"
)

type Mux struct {
	r           *mux.Router
	baseURL     string
	contentType string
}

// Ensure `Invoke` conforms to `functions.Invoke`
var _ = (functions.Invoke)(((*Mux)(nil)).Invoke)

// Ensure `Register` conforms to `functions.Register`
var _ = (functions.Register)(((*Mux)(nil)).Register)

// Ensure `RegisterStateful` conforms to `functions.RegisterStateful`
var _ = (functions.RegisterStateful)(((*Mux)(nil)).RegisterStateful)

func New(baseURL string, contentType string) *Mux {
	return &Mux{
		r:           mux.NewRouter(),
		baseURL:     baseURL,
		contentType: contentType,
	}
}

func (m *Mux) Router() *mux.Router {
	return m.r
}

func (m *Mux) Invoke(ctx context.Context, receiver functions.Receiver, payload []byte) ([]byte, error) {
	u, err := url.Parse(m.baseURL)
	if err != nil {
		return nil, err
	}
	operationPath := receiver.Namespace + "/" + receiver.Operation
	u.Path = path.Join(u.Path, operationPath)
	resp, err := http.Post(
		u.String(),
		m.contentType,
		bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if response == nil {
		response = []byte{}
	}
	if resp.StatusCode/100 != 2 {
		return nil, errors.New(string(response))
	}

	return response, err
}

func (m *Mux) Register(namespace, operation string, handler functions.Handler) {
	m.r.HandleFunc("/"+namespace+"/"+operation, wrap(handler))
}

func (m *Mux) RegisterStateful(namespace, operation string, handler functions.StatefulHandler) {
	m.r.HandleFunc("/"+namespace+"/{id}/"+operation, wrapMethod(handler))
}

func wrap(handler functions.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		response, err := handler(r.Context(), body)
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func wrapMethod(handler functions.StatefulHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		body, err := io.ReadAll(r.Body)
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		response, err := handler(r.Context(), id, body)
		if err != nil {
			handleError(err, w, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func handleError(err error, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	fmt.Fprint(w, err.Error())
}
