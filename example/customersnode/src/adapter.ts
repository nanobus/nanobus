import logger from "./lib/logger";
import dotenv from "dotenv";
import {
  HTTPHandlers,
  HTTPInvoker,
  Invoker,
  msgpackCodec,
} from "./lib/nanobus";
import { Customer, Inbound, Outbound } from "./interfaces";

interface GetCustomerArgs {
  id: number;
}

export function registerInboundHanders(h: Inbound): void {
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
        const args = handlers.codec.decoder(input) as GetCustomerArgs;
        return h
          .getCustomer(args.id)
          .then((result) => handlers.codec.encoder(result));
      }
    );
  }
}

export class OutboundImpl implements Outbound {
  private invoker: Invoker;

  constructor(invoker: Invoker) {
    this.invoker = invoker;
  }

  async saveCustomer(customer: Customer): Promise<void> {
    return this.invoker("customers.v1.Outbound", "saveCustomer", customer);
  }

  async fetchCustomer(id: number): Promise<Customer> {
    const args: GetCustomerArgs = {
      id,
    };
    return this.invoker("customers.v1.Outbound", "fetchCustomer", args);
  }

  async customerCreated(customer: Customer): Promise<void> {
    return this.invoker("customers.v1.Outbound", "customerCreated", customer);
  }
}

export const invoker = HTTPInvoker(
  process.env.OUTBOUND_BASE_URL || "http://localhost:32321/outbound",
  msgpackCodec
);

export const outbound = new OutboundImpl(invoker);

const result = dotenv.config();
if (result.error) {
  dotenv.config({ path: ".env.default" });
}

const PORT = parseInt(process.env.PORT) || 9000;
const HOST = process.env.HOST || "localhost";

export const handlers = new HTTPHandlers(msgpackCodec);

export function start(): void {
  handlers.listen(PORT, HOST, () => {
    logger.info(`🌏 Nanoprocess server started at http://${HOST}:${PORT}`);
  });
}
