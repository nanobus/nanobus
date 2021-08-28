package io.nanobus.example.customers;

import reactor.core.publisher.*;

public class InboundImpl implements Inbound {
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
