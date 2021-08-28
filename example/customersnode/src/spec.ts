import { handlers, invoker, Invoker, Handler, Codec } from './nanobus';

export interface GetCustomerArgs {
  id: number;
}

export interface Customer {
  id: number;
  firstName: string;
  middleName?: string;
  lastName: string;
  email: string;
  address?: Address;
}

export interface Address {
  line1: string;
  line2?: string;
  city: string;
  state: string;
  zip: string;
}

export interface Outbound {
  saveCustomer(customer: Customer): Promise<void>;
  fetchCustomer(id: number): Promise<Customer>;
  customerCreated(customer: Customer): Promise<void>;
}

export class OutboundImpl implements Outbound {
  private invoker: Invoker;

  constructor(invoker: Invoker) {
    this.invoker = invoker;
  }

  async saveCustomer(customer: Customer): Promise<void> {
    await this.invoker('/customers.v1.Outbound/saveCustomer', customer);
  }

  async fetchCustomer(id: number): Promise<Customer> {
    return this.invoker('/customers.v1.Outbound/fetchCustomer', {
      id: id
    });
  }

  async customerCreated(customer: Customer): Promise<void> {
    await this.invoker('/customers.v1.Outbound/customerCreated', customer);
  }
}

export type createCustomerHandler = (customer: Customer) => Promise<Customer>;
export type getCustomerHandler = (id: number) => Promise<Customer>;


function createCustomerWrapper(codec: Codec, handler?: createCustomerHandler): Handler {
  return (input: ArrayBuffer): Promise<ArrayBuffer> => {
    const payload = codec.decoder(input) as Customer;
    return handler(payload).then(result => codec.encoder(result));
  };
}

function getCustomerWrapper(codec: Codec, handler?: getCustomerHandler): Handler {
  return (input: ArrayBuffer): Promise<ArrayBuffer> => {
    const args = codec.decoder(input) as GetCustomerArgs;
    return handler(args.id).then(result => codec.encoder(result));
  };
}

export interface InboundHanders {
  createCustomer?: createCustomerHandler;
  getCustomer?: getCustomerHandler;
}

export function registerInboundHanders(params: InboundHanders): void {
  if (params.createCustomer) {
    handlers.registerHandler(
        '/customers.v1.Inbound/createCustomer',
        createCustomerWrapper(handlers.codec, params.createCustomer));
  }
  if (params.getCustomer) {
    handlers.registerHandler(
      '/customers.v1.Inbound/getCustomer',
        getCustomerWrapper(handlers.codec, params.getCustomer));
  }
}

export const outbound = new OutboundImpl(invoker);