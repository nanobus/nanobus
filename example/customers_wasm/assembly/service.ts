import { consoleLog } from "@wapc/as-guest";
import {
  Inbound,
  CustomerActor,
  Outbound,
  Customer,
  CustomerQuery,
  CustomerPage,
  Address,
  Error,
} from "./interfaces";

export class InboundImpl implements Inbound {
  private outbound: Outbound;

  constructor(outbound: Outbound) {
    this.outbound = outbound;
  }

  createCustomer(customer: Customer): Customer {
    consoleLog("createCustomer called");
    this.outbound.saveCustomer(customer);
    this.outbound.customerCreated(customer);

    return customer;
  }

  getCustomer(id: u64): Customer {
    consoleLog("getCustomer called");
    return this.outbound.fetchCustomer(id);
  }

  listCustomers(query: CustomerQuery): CustomerPage {
    consoleLog("listCustomers called");
    return CustomerPage.newBuilder()
      .withOffset(query.offset)
      .withLimit(query.limit)
      .build();
  }
}
