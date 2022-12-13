// Code generated by NanoBus codegen utilities. DO NOT EDIT.

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

export function HttpServerV1(
  config: HttpServerV1Config
): Component<HttpServerV1Config> {
  return {
    uses: "nanobus.transport.http.server/v1",
    with: config
  };
}