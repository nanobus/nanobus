import logger from "./lib/logger";
import dotenv from "dotenv";
import {
  HTTPHandlers,
  HTTPInvoker,
  Invoker,
  msgpackCodec,
} from "./lib/nanobus";
import {
  Inbound,
  Customer,
  CustomerPage,
  CustomerQuery,
  Outbound,
} from "./interfaces";

export const invoker = HTTPInvoker(
  process.env.OUTBOUND_BASE_URL || "http://localhost:32321/outbound",
  msgpackCodec
);

export const handlers = new HTTPHandlers(msgpackCodec);

class InboundGetCustomerArgs {
  id: number;

  constructor({ id = 0 }: { id?: number } = {}) {
    this.id = id;
  }
}

export function registerInboundHandlers(h: Inbound): void {
  if (h.createCustomer) {
    handlers.registerHandler(
      "customers.v1.Inbound",
      "createCustomer",
      (input: ArrayBuffer): Promise<ArrayBuffer> => {
        const payload = handlers.codec.decoder(input) as Customer;
        return h
          .createCustomer(payload)
          .then((result) => handlers.codec.encoder(result));
      }
    );
  }
  if (h.getCustomer) {
    handlers.registerHandler(
      "customers.v1.Inbound",
      "getCustomer",
      (input: ArrayBuffer): Promise<ArrayBuffer> => {
        const inputArgs = handlers.codec.decoder(
          input
        ) as InboundGetCustomerArgs;
        return h
          .getCustomer(inputArgs.id)
          .then((result) => handlers.codec.encoder(result));
      }
    );
  }
  if (h.listCustomers) {
    handlers.registerHandler(
      "customers.v1.Inbound",
      "listCustomers",
      (input: ArrayBuffer): Promise<ArrayBuffer> => {
        const payload = handlers.codec.decoder(input) as CustomerQuery;
        return h
          .listCustomers(payload)
          .then((result) => handlers.codec.encoder(result));
      }
    );
  }
}

class OutboundFetchCustomerArgs {
  id: number;

  constructor({ id = 0 }: { id?: number } = {}) {
    this.id = id;
  }
}

export class OutboundImpl implements Outbound {
  private invoker: Invoker;

  constructor(invoker: Invoker) {
    this.invoker = invoker;
  }

  async saveCustomer(customer: Customer): Promise<void> {
    return this.invoker(
      "customers.v1.Outbound",
      "saveCustomer",
      customer
    ).then();
  }

  async fetchCustomer(id: number): Promise<Customer> {
    const inputArgs: OutboundFetchCustomerArgs = {
      id,
    };
    return this.invoker("customers.v1.Outbound", "fetchCustomer", inputArgs);
  }

  async customerCreated(customer: Customer): Promise<void> {
    return this.invoker(
      "customers.v1.Outbound",
      "customerCreated",
      customer
    ).then();
  }
}

export var outbound = new OutboundImpl(invoker);

const result = dotenv.config();
if (result.error) {
  dotenv.config({ path: ".env.default" });
}

const PORT = parseInt(process.env.PORT) || 9000;
const HOST = process.env.HOST || "localhost";

export function start(): void {
  handlers.listen(PORT, HOST, () => {
    logger.info(`🌏 Nanoprocess server started at http://${HOST}:${PORT}`);
  });
}
