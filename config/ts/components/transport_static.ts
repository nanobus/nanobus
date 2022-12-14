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

export interface StaticV1Config {
  paths: StaticPath[];
}

export function StaticV1(config: StaticV1Config): Component<StaticV1Config> {
  return {
    uses: "nanobus.transport.http.static/v1",
    with: config
  };
}

export interface StaticPath {
  path: string;
  dir?: string;
  file?: string;
  strip?: string;
}
