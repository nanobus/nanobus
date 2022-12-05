// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export type RouterV1Config = Array<AddRoute>;

export class RouterV1 implements Component<RouterV1Config> {
  readonly uses: string = "nanobus.transport.http.router/v1";
  readonly with: RouterV1Config;

  constructor(config: RouterV1Config) {
    this.with = config;
  }
}
export interface AddRoute {
  methods: string;
  uri: string;
  encoding?: string;
  handler: Handler;
}
