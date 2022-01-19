import logger from "./lib/logger";
import dotenv from "dotenv";
import { Manager, LRUCache, Deactivator, isDeactivator } from "./lib/stateful";
import { Expose } from "class-transformer";
import {
  Inbound,
  Customer,
  CustomerQuery,
  CustomerActor,
  Outbound,
} from "./interfaces";
import {
  jsonSerializer,
  msgPackSerializer,
  IAdapter,
  Metadata,
  Source,
  IPublisher,
} from "./lib/nanobus";
import { RSocketAdapter, RSocketStorage } from "./lib/rsocket";

const result = dotenv.config();
if (result.error) {
  dotenv.config({ path: ".env.default" });
}

const HOST = process.env.RSOCKET_HOST || "127.0.0.1";
const PORT = parseInt(process.env.RSOCKET_PORT) || 7878;

const adapter = new RSocketAdapter(msgPackSerializer, {
  host: HOST,
  port: PORT,
});

const busUrl = process.env.BUS_URL || "http://127.0.0.1:32321";
const cache = new LRUCache();
const storage = new RSocketStorage(adapter);
const stateManager = new Manager(cache, storage, jsonSerializer);

export function registerInbound(h: Inbound): void {
  if (h.createCustomer) {
    adapter.registerRequestResponseHandler(
      "/customers.v1.Inbound/createCustomer",
      async (_: Metadata, input: any): Promise<any> => {
        const payload = input as Customer;
        return h.createCustomer(payload);
      }
    );
  }
  if (h.getCustomer) {
    adapter.registerRequestResponseHandler(
      "/customers.v1.Inbound/getCustomer",
      async (_: Metadata, input: any): Promise<any> => {
        const inputArgs = input as InboundGetCustomerArgs;
        return h.getCustomer(inputArgs.id);
      }
    );
  }
  if (h.listCustomers) {
    adapter.registerRequestResponseHandler(
      "/customers.v1.Inbound/listCustomers",
      async (_: Metadata, input: any): Promise<any> => {
        const payload = input as CustomerQuery;
        return h.listCustomers(payload);
      }
    );
  }
}

class InboundGetCustomerArgs {
  @Expose() id: number;

  constructor({ id = 0 }: { id?: number } = {}) {
    this.id = id;
  }
}

export function registerCustomerActor(actor: CustomerActor): void {
  adapter.registerRequestResponseHandler(
    "/customers.v1.CustomerActor/deactivate",
    async (md: Metadata, _: any): Promise<void> => {
      const id = md[":id"][0];
      const sctx = stateManager.toContext(
        "customers.v1.CustomerActor",
        id,
        actor
      );
      if (isDeactivator(actor)) {
        (actor as Deactivator).deactivate(sctx);
      }
      stateManager.deactivate(sctx.self);
    }
  );
  adapter.registerRequestResponseHandler(
    "/customers.v1.CustomerActor/createCustomer",
    async (md: Metadata, input: any): Promise<any> => {
      const id = md[":id"][0];
      const payload = input as Customer;
      const sctx = stateManager.toContext(
        "customers.v1.CustomerActor",
        id,
        actor
      );
      return actor
        .createCustomer(sctx, payload)
        .then((result) => sctx.response(result));
    }
  );
  adapter.registerRequestResponseHandler(
    "/customers.v1.CustomerActor/getCustomer",
    async (md: Metadata, _: any): Promise<any> => {
      const id = md[":id"][0];
      const sctx = stateManager.toContext(
        "customers.v1.CustomerActor",
        id,
        actor
      );
      return actor.getCustomer(sctx).then((result) => sctx.response(result));
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
  private adapter: IAdapter;

  constructor(adapter: IAdapter) {
    this.adapter = adapter;
  }

  // Saves a customer to the backend database
  async saveCustomer(customer: Customer): Promise<void> {
    return this.adapter
      .requestResponse<void>("/customers.v1.Outbound/saveCustomer", customer)
      .then();
  }

  // Fetches a customer from the backend database
  async fetchCustomer(id: number): Promise<Customer> {
    const inputArgs: OutboundFetchCustomerArgs = {
      id,
    };
    return this.adapter
      .requestResponse<Customer>(
        "/customers.v1.Outbound/fetchCustomer",
        inputArgs
      )
      .then();
  }

  // Queries customers from the backend database
  getCustomers(): IPublisher<Customer> {
    return this.adapter.requestStream<Customer>(
      "/customers.v1.Outbound/getCustomers"
    );
  }

  // Sends a customer creation event
  async customerCreated(customer: Customer): Promise<void> {
    return this.adapter
      .requestResponse("/customers.v1.Outbound/customerCreated", customer)
      .then();
  }

  transformCustomers(
    prefix: string,
    source: Source<Customer>
  ): IPublisher<Customer> {
    return this.adapter.requestChannel(
      "/customers.v1.Outbound/transformCustomers",
      {
        prefix: prefix,
      },
      source
    );
  }
}

export var outbound = new OutboundImpl(adapter);

export function start(): void {
  Promise.resolve(adapter.connect()).then(
    () => {},
    (error) => {
      console.error(error.stack);
      process.exit(1);
    }
  );
}
