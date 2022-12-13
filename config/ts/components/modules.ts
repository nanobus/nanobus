import { Application, Module } from "../nanobus.ts";
import { RestV1 } from "./transport_rest.ts";
import { HttpServerV1 } from "./transport_server.ts";

export const standardErrors = {
  not_found: {
    type: "NotFound",
    code: "not_found",
    title: "Resource not found",
    message: "Resource with id {{ .key }} was not found",
  },
  permission_denied: {
    type: "PermissionDenied",
    code: "permission_denied",
    title: "Permission denied",
    message:
      "You don't have permission to access this resource or to perform the operation.",
  },
  unauthenticated: {
    type: "Unauthenticated",
    code: "unauthenticated",
    title: "Unauthenticated",
    message: "You must be logged in to perform the operation.",
  },
}

export class RestModule implements Module {
  private address: string;

  constructor(address: string) {
    this.address = address;
  }

  initialize(app: Application): void {
    app.transport(
      "http",
      HttpServerV1({
        address: this.address,
        routes: [
          RestV1({
            documentation: {
              swaggerUI: true,
              postman: true,
              restClient: true,
            },
          }),
        ],
      })
    );

    app.errors(standardErrors);
  }
}
