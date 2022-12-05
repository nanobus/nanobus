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

export class StaticV1 implements Component<StaticV1Config> {
  readonly uses: string = "nanobus.transport.http.static/v1";
  readonly with: StaticV1Config;

  constructor(config: StaticV1Config) {
    this.with = config;
  }
}
export interface StaticPath {
  dir: string;
  path: string;
  strip?: string;
}
