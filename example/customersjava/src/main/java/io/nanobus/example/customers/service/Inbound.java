package io.nanobus.example.customers.service;

import reactor.core.publisher.*;

import io.nanobus.example.customers.models.*;

public abstract class Inbound {
  public Mono<Customer> createCustomer(Customer customer) {
    throw new UnsupportedOperationException("createCustomer is not implemented");
  }

  public Mono<Customer> getCustomer(long id) {
    throw new UnsupportedOperationException("getCustomer is not implemented");
  }
}
