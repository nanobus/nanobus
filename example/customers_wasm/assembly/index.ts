import { Customer, Handlers, Host } from "./module";

var host = new Host();

export function wapc_init(): void {
  Handlers.registerCreateCustomer(createCustomer);
  Handlers.registerGetCustomer(getCustomer);
}

function createCustomer(customer: Customer): Customer {
  consoleLog("createCustomer called");
  host.saveCustomer(customer);
  host.customerCreated(customer);

  return customer;
}

function getCustomer(id: u64): Customer {
  consoleLog("getCustomer called");
  return host.fetchCustomer(id);
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
