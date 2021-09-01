package io.nanobus.example.customers.adapter;

import io.nanobus.functions.Invoker;
import reactor.core.publisher.*;

import io.nanobus.example.customers.models.*;
import io.nanobus.example.customers.service.Outbound;

public class OutboundImpl implements Outbound {
  private Invoker invoker;

  public OutboundImpl(Invoker invoker) {
    this.invoker = invoker;
  }

  @Override
  public Mono<Void> saveCustomer(Customer customer) {
    return this.invoker.invoke( //
        "customers.v1.Outbound", "saveCustomer", customer);
  }

  @Override
  public Mono<Customer> fetchCustomer(long id) {
    GetCustomerArgs args = new GetCustomerArgs();
    args.setId(id);

    return this.invoker.invokeWithReturn( //
        "customers.v1.Outbound", "fetchCustomer", args, Customer.class);
  }

  @Override
  public Mono<Void> customerCreated(Customer customer) {
    return this.invoker.invoke( //
        "customers.v1.Outbound", "customerCreated", customer);
  }
}
