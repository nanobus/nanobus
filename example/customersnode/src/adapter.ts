import logger from './lib/logger';
import dotenv from 'dotenv';
import { Codec, Handler, HTTPHandlers, HTTPInvoker, Invoker, msgpackCodec } from './lib/nanobus';
import { createCustomerHandler, Customer, getCustomerHandler, InboundHanders, Outbound } from './interfaces';


interface GetCustomerArgs {
  id: number;
}

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

export function registerInboundHanders(params: InboundHanders): void {
  if (params.createCustomer) {
    handlers.registerHandler(
      'customers.v1.Inbound', 'createCustomer',
      createCustomerWrapper(handlers.codec, params.createCustomer));
  }
  if (params.getCustomer) {
    handlers.registerHandler(
      'customers.v1.Inbound', 'getCustomer',
      getCustomerWrapper(handlers.codec, params.getCustomer));
  }
}

export class OutboundImpl implements Outbound {
  private invoker: Invoker;

  constructor(invoker: Invoker) {
    this.invoker = invoker;
  }

  async saveCustomer(customer: Customer): Promise<void> {
    await this.invoker('customers.v1.Outbound', 'saveCustomer', customer);
  }

  async fetchCustomer(id: number): Promise<Customer> {
    const args: GetCustomerArgs = {
      id: id
    };
    return this.invoker('customers.v1.Outbound', 'fetchCustomer', args);
  }

  async customerCreated(customer: Customer): Promise<void> {
    await this.invoker('customers.v1.Outbound', 'customerCreated', customer);
  }
}

export const invoker = HTTPInvoker(
    process.env.OUTBOUND_BASE_URL || 'http://localhost:32321/outbound',
    msgpackCodec
  );

export const outbound = new OutboundImpl(invoker);

const result = dotenv.config();
if (result.error) {
  dotenv.config({ path: '.env.default' });
}

const PORT = parseInt(process.env.PORT) || 9000;
const HOST = process.env.HOST || 'localhost';

export const handlers = new HTTPHandlers(msgpackCodec);

export function start(): void {
  handlers.listen(PORT, HOST, () => {
    logger.info(`🌏 Nanoprocess server started at http://${HOST}:${PORT}`);
  });
}
