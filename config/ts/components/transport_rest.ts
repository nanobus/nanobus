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

export class RestV1 implements Component<RestV1Config> {
  readonly uses: string = "nanobus.transport.http.rest/v1";
  readonly with: RestV1Config;

  constructor(config: RestV1Config) {
    this.with = config;
  }
}
export interface Documentation {
  swaggerUI?: boolean;
  postman?: boolean;
  restClient?: boolean;
}
