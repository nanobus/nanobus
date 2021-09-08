import { registerInboundHanders, start, outbound } from "./adapter";
import { Customer } from "./interfaces";

class InboundHandlers {
  async createCustomer(customer: Customer): Promise<Customer> {
    await outbound.saveCustomer(customer);
    await outbound.customerCreated(customer);

    return customer;
  }

  async getCustomer(id: number): Promise<Customer> {
    return outbound.fetchCustomer(id);
  };
}

registerInboundHanders(new InboundHandlers());

start();
