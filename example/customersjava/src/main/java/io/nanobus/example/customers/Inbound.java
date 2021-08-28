package io.nanobus.example.customers;

import reactor.core.publisher.*;

public interface Inbound {
  Mono<Customer> createCustomer(Customer customer);
  Mono<Customer> getCustomer(long id);
}
