
import { start } from './nanobus';
import { Customer, registerInboundHanders, outbound } from './spec';

async function createCustomer(customer: Customer): Promise<Customer> {
  await outbound.saveCustomer(customer);
  await outbound.customerCreated(customer);

  return customer;
};

async function getCustomer(id: number): Promise<Customer> {
  return outbound.fetchCustomer(id);
};

registerInboundHanders({
  createCustomer: createCustomer,
  getCustomer: getCustomer,
});

start();
