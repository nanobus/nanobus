package welcome

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
)

type Service struct {
	Outbound
	receiveCounter uint64
}

func New(outbound Outbound) *Service {
	return &Service{
		Outbound: outbound,
	}
}

func (s *Service) GreetCustomer(ctx context.Context, customer Customer) error {
	counter := atomic.AddUint64(&s.receiveCounter, 1)
	if counter%2 == 0 {
		log.Printf("RETURNING SIMULATED ERROR")
		return errors.New("simulated error")
	}

	if jsonBytes, err := json.MarshalIndent(&customer, "", "  "); err == nil {
		log.Printf("RECEIVED: %s", string(jsonBytes))
	}

	message := fmt.Sprintf("Hello, %s %s", customer.FirstName, customer.LastName)

	return s.SendEmail(ctx, customer.Email, message)
}
