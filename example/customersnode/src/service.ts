import {
  Context,
  Customer,
  CustomerPage,
  CustomerQuery,
  Inbound,
  Outbound,
} from "./interfaces";

export class InboundImpl implements Inbound {
  private outbound: Outbound;

  constructor(outbound: Outbound) {
    this.outbound = outbound;
  }

  async createCustomer(customer: Customer): Promise<Customer> {
    await this.outbound.saveCustomer(customer);
    await this.outbound.customerCreated(customer);

    return customer;
  }

  async getCustomer(id: number): Promise<Customer> {
    const stream = this.outbound.getCustomers();
    await stream.forEach(async (customer) => {
      console.log(customer);
    });
    return this.outbound.fetchCustomer(id);
  }

  async listCustomers(query: CustomerQuery): Promise<CustomerPage> {
    return new CustomerPage({
      offset: query.offset,
      limit: query.limit,
      items: [new Customer()],
    });
  }
}

export class CustomerActorImpl {
  async createCustomer(ctx: Context, customer: Customer): Promise<Customer> {
    ctx.set("customer", customer);
    return customer;
  }

  async getCustomer(ctx: Context): Promise<Customer> {
    const customer = await ctx.get("customer", Customer);
    return customer;
  }
}
