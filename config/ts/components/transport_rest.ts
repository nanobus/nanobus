// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export interface RestV1Config {
  documentation: Documentation;
}

export function RestV1(config: RestV1Config): Component<RestV1Config> {
  return {
    uses: "nanobus.transport.http.rest/v1",
    with: config
  };
}

export interface Documentation {
  swaggerUI?: boolean;
  postman?: boolean;
  restClient?: boolean;
}
