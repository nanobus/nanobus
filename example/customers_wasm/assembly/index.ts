import { OutboundImpl, registerInbound } from "./adapters";
import {
  Customer,
  CustomerQuery,
  CustomerPage,
  Inbound,
  Outbound,
} from "./interfaces";

var outbound: Outbound;

export function wapc_init(): void {
  outbound = new OutboundImpl()
  registerInbound(new InboundImpl());
}

class InboundImpl implements Inbound {
  createCustomer(customer: Customer): Customer {
    consoleLog("createCustomer called");
    outbound.saveCustomer(customer);
    outbound.customerCreated(customer);

    return customer;
  }

  getCustomer(id: u64): Customer {
    consoleLog("getCustomer called");
    return outbound.fetchCustomer(id);
  }

  listCustomers(query: CustomerQuery): CustomerPage {
    consoleLog("listCustomers called");
    return CustomerPage.newBuilder()
      .withOffset(query.offset)
      .withLimit(query.limit)
      .build();
  }
}

// Boilerplate code for waPC.  Do not remove.

import { handleCall, handleAbort, consoleLog } from "@wapc/as-guest";

export function __guest_call(operation_size: usize, payload_size: usize): bool {
  return handleCall(operation_size, payload_size);
}

// Abort function
function abort(
  message: string | null,
  fileName: string | null,
  lineNumber: u32,
  columnNumber: u32
): void {
  handleAbort(message, fileName, lineNumber, columnNumber);
}
