package io.nanobus.example.customers;

import reactor.core.publisher.*;

public interface Outbound {
  Mono<Void> saveCustomer(Customer customer);
  Mono<Customer> fetchCustomer(long id);
  Mono<Void> customerCreated(Customer customer);
}
