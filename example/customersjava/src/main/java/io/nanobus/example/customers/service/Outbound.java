package io.nanobus.example.customers.service;

import reactor.core.publisher.*;

import io.nanobus.example.customers.models.*;

public interface Outbound {
  Mono<Void> saveCustomer(Customer customer);
  Mono<Customer> fetchCustomer(long id);
  Mono<Void> customerCreated(Customer customer);
}
