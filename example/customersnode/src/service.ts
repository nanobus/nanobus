import { start, registerInboundHandlers, outbound } from "./adapter";
import { Customer, CustomerPage, CustomerQuery } from "./interfaces";

class InboundHandlers {
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

registerInboundHandlers(new InboundHandlers());

start();
