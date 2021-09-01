package io.nanobus.example.customers.service;

import reactor.core.publisher.*;

import io.nanobus.example.customers.models.*;

public class InboundImpl extends Inbound {
  private Outbound outbound;

  public InboundImpl(Outbound outbound) {
    this.outbound = outbound;
  }

  @Override
  public Mono<Customer> createCustomer(Customer customer) {
    return outbound.saveCustomer(customer)
            .then(outbound.customerCreated(customer))
            .then(Mono.just(customer));
  }

  @Override
  public Mono<Customer> getCustomer(long id) {
    return outbound.fetchCustomer(id);
  }
}
