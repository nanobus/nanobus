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

export function RouterV1(config: RouterV1Config): Component<RouterV1Config> {
  return {
    uses: "nanobus.transport.http.router/v1",
    with: config
  };
}

export interface AddRoute {
  methods: string;
  uri: string;
  encoding?: string;
  handler: Handler;
}
