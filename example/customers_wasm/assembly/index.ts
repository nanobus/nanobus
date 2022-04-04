import { registerService, RepositoryImpl, PublisherImpl } from "./adapter";
import { ServiceImpl } from "./service";

export function wapc_init(): void {
  const repository = new RepositoryImpl();
  const publisher = new PublisherImpl();
  registerService(new ServiceImpl(repository, publisher));
}

// Boilerplate code for waPC.  Do not remove.

import { handleCall, handleAbort } from "@wapc/as-guest";

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
