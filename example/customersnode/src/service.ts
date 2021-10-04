import {
  start,
  registerInbound,
  registerCustomerActor,
  outbound,
} from "./adapter";
import { Context, Customer, CustomerPage, CustomerQuery } from "./interfaces";

class InboundImpl {
  async createCustomer(customer: Customer): Promise<Customer> {
    await outbound.saveCustomer(customer);
    await outbound.customerCreated(customer);

    return customer;
  }

  async getCustomer(id: number): Promise<Customer> {
    return outbound.fetchCustomer(id);
  }

  async listCustomers(query: CustomerQuery): Promise<CustomerPage> {
    return new CustomerPage({
      offset: query.offset,
      limit: query.limit,
      items: [new Customer()],
    });
  }
}

registerInbound(new InboundImpl());

class CustomerActorImpl {
  async createCustomer(ctx: Context, customer: Customer): Promise<Customer> {
    ctx.set("customer", customer);
    return customer;
  }

  async getCustomer(ctx: Context): Promise<Customer> {
    const customer: Customer = await ctx.get("customer");
    return customer;
  }
}

registerCustomerActor(new CustomerActorImpl());

start();
