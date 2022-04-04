import { consoleLog as log } from "@wapc/as-guest";
import {
  Service,
  Customer,
  CustomerQuery,
  CustomerPage,
  Repository,
  Publisher,
} from "./interfaces";

export class ServiceImpl implements Service {
  private repository: Repository;
  private publisher: Publisher;

  constructor(repository: Repository, publisher: Publisher) {
    this.repository = repository;
    this.publisher = publisher;
  }

  createCustomer(customer: Customer): Customer {
    log("createCustomer called");
    this.repository.saveCustomer(customer);
    this.publisher.customerCreated(customer);

    return customer;
  }

  getCustomer(id: u64): Customer {
    log("getCustomer called");
    return this.repository.fetchCustomer(id);
  }

  listCustomers(query: CustomerQuery): CustomerPage {
    log("listCustomers called");
    return CustomerPage.newBuilder()
      .withOffset(query.offset)
      .withLimit(query.limit)
      .build();
  }
}
