import { Invoker, Handler, Handlers, Codec } from "./functions";

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

export type createCustomerHandler = (customer: Customer) => Promise<Customer>
export type getCustomerHandler = (id: number) => Promise<Customer>

export interface InboundHandlers {
  createCustomer?: createCustomerHandler
  getCustomer?: getCustomerHandler
}

export class Service {
  private handlers: Handlers;
  private codec: Codec;
  private starter: () => void;

  constructor(handlers: Handlers, codec: Codec, starter: () => void) {
    this.handlers = handlers;
    this.codec = codec;
    this.starter = starter;
  }

  registerInboundHandlers(handlers: InboundHandlers): Service {
    this.handlers.registerHandlers({
      '/customers.v1.Inbound/createCustomer': this.createCustomerWrapper(handlers.createCustomer),
      '/customers.v1.Inbound/getCustomer': this.getCustomerWrapper(handlers.getCustomer)
    })
    return this;
  }

  createCustomerWrapper(handler?: createCustomerHandler): Handler {
    if (!handler) {
      return undefined;
    }
    return (input: ArrayBuffer): Promise<ArrayBuffer> => {
      const payload = this.codec.decoder(input) as Customer;
      return handler(payload)
        .then((result) => this.codec.encoder(result));
    }
  }

  getCustomerWrapper(handler?: getCustomerHandler): Handler {
    if (!handler) {
      return undefined;
    }
    return (input: ArrayBuffer): Promise<ArrayBuffer> => {
      const args = this.codec.decoder(input) as GetCustomerArgs;
      return handler(args.id)
        .then((result) => this.codec.encoder(result));
    }
  }

  start(): void {
    this.starter();
  }
}
