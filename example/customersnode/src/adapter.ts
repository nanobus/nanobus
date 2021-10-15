import logger from "./lib/logger";
import dotenv from "dotenv";
import {
  HTTPHandlers,
  HTTPInvoker,
  Invoker,
  jsonCodec,
  msgpackCodec,
} from "./lib/nanobus";
import { Manager, Storage, LRUCache } from "./lib/stateful";
import { Expose } from "class-transformer";
import {
  Inbound,
  Customer,
  CustomerQuery,
  CustomerActor,
  Outbound,
} from "./interfaces";

const busUrl = process.env.BUS_URL || "http://127.0.0.1:32321";

const invoker = HTTPInvoker(busUrl + "/providers", msgpackCodec);
const handlers = new HTTPHandlers(msgpackCodec);
const cache = new LRUCache();
const storage = new Storage(busUrl, jsonCodec);
const stateManager = new Manager(cache, storage, jsonCodec);

class InboundGetCustomerArgs {
  @Expose() id: number;

  constructor({ id = 0 }: { id?: number } = {}) {
    this.id = id;
  }
}

export function registerInbound(h: Inbound): void {
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

export function registerCustomerActor(h: CustomerActor): void {
  handlers.registerStatefulHandler(
    "customers.v1.CustomerActor",
    "deactivate",
    stateManager.deactivateHandler("customers.v1.CustomerActor", h)
  );
  handlers.registerStatefulHandler(
    "customers.v1.CustomerActor",
    "createCustomer",
    (id: string, input: ArrayBuffer): Promise<ArrayBuffer> => {
      const payload = handlers.codec.decoder(input) as Customer;
      const sctx = stateManager.toContext("customers.v1.CustomerActor", id, h);
      return h
        .createCustomer(sctx, payload)
        .then((result) => sctx.response(result))
        .then((result) => handlers.codec.encoder(result));
    }
  );
  handlers.registerStatefulHandler(
    "customers.v1.CustomerActor",
    "getCustomer",
    (id: string, input: ArrayBuffer): Promise<ArrayBuffer> => {
      const sctx = stateManager.toContext("customers.v1.CustomerActor", id, h);
      return h
        .getCustomer(sctx)
        .then((result) => sctx.response(result))
        .then((result) => handlers.codec.encoder(result));
    }
  );
}

class OutboundFetchCustomerArgs {
  @Expose() id: number;

  constructor({ id = 0 }: { id?: number } = {}) {
    this.id = id;
  }
}

export class OutboundImpl implements Outbound {
  private invoker: Invoker;

  constructor(invoker: Invoker) {
    this.invoker = invoker;
  }

  // Saves a customer to the backend database
  async saveCustomer(customer: Customer): Promise<void> {
    return this.invoker(
      "customers.v1.Outbound",
      "saveCustomer",
      customer
    ).then();
  }

  // Fetches a customer from the backend database
  async fetchCustomer(id: number): Promise<Customer> {
    const inputArgs: OutboundFetchCustomerArgs = {
      id,
    };
    return this.invoker("customers.v1.Outbound", "fetchCustomer", inputArgs);
  }

  // Sends a customer creation event
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
const HOST = process.env.HOST || "127.0.0.1";

export function start(): void {
  handlers.listen(PORT, HOST, () => {
    logger.info(`🌏 Nanoserver started at http://${HOST}:${PORT}`);
  });
}
