package io.nanobus.example.customers;

import io.netty.buffer.*;
import reactor.core.publisher.*;

import reactor.netty.http.client.HttpClient;

public class OutboundImpl implements Outbound {
  private HttpClient client;
  private Mapper mapper;

  public OutboundImpl(HttpClient client, Mapper mapper) {
    this.client = client;
    this.mapper = mapper;
  }

  @Override
  public Mono<Void> saveCustomer(Customer customer) {
    return this.client
            .post()
            .uri("/customers.v1.Outbound/saveCustomer")
            .send((request, outbound) -> outbound.sendByteArray(Mono.just(mapper.encode(customer))))
            .response().then();
  }

  @Override
  public Mono<Customer> fetchCustomer(long id) {
    GetCustomerArgs args = new GetCustomerArgs();
    args.setId(id);

    return this.client
            .post()
            .uri("/customers.v1.Outbound/fetchCustomer")
            .send((request, outbound) -> outbound.sendByteArray(Mono.just(mapper.encode(args))))
            .responseContent().retain().aggregate()
            .map(res -> mapper.decode(ByteBufUtil.getBytes(res), Customer.class));
  }

  @Override
  public Mono<Void> customerCreated(Customer customer) {
    return this.client
            .post()
            .uri("/customers.v1.Outbound/customerCreated")
            .send((request, outbound) -> outbound.sendByteArray(Mono.just(mapper.encode(customer))))
            .responseContent().retain().aggregate()
            .then();
  }
}
