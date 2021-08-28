
import { Customer } from './interfaces';
import { outbound, service } from './dependencies'

async function createCustomer(customer: Customer): Promise<Customer> {
  await outbound.saveCustomer(customer);
  await outbound.customerCreated(customer);

  return customer;
};

async function getCustomer(id: number): Promise<Customer> {
  return outbound.fetchCustomer(id);
};

service.registerInboundHandlers({
  createCustomer: createCustomer,
  getCustomer: getCustomer,
});

service.start();