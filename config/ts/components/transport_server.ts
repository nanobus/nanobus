// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export interface HttpServerV1Config {
  address: string;
  routes?: Component<any>[];
  middleware?: Component<any>[];
}

export class HttpServerV1 implements Component<HttpServerV1Config> {
  readonly uses: string = "nanobus.transport.http.server/v1";
  readonly with: HttpServerV1Config;

  constructor(config: HttpServerV1Config) {
    this.with = config;
  }
}
