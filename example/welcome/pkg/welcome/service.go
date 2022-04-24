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

package welcome

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/go-logr/logr"
)

type Service struct {
	log            logr.Logger
	outbound       Outbound
	receiveCounter uint64
}

func NewService(log logr.Logger, outbound Outbound) *Service {
	return &Service{
		log:      log,
		outbound: outbound,
	}
}

func (s *Service) GreetCustomer(ctx context.Context, customer *Customer) error {
	counter := atomic.AddUint64(&s.receiveCounter, 1)
	if counter%2 == 0 {
		log.Printf("RETURNING SIMULATED ERROR")
		return errors.New("simulated error")
	}

	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("RECEIVED: %s", string(jsonBytes))
	}

	message := fmt.Sprintf("Hello, %s %s", customer.FirstName, customer.LastName)

	return s.outbound.SendEmail(ctx, customer.Email, message)
}
