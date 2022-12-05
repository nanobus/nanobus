// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export type FilePath = string;

export interface MigratePostgresV1Config {
  dataSource: string;
  directory?: FilePath;
  sourceUrl?: string;
}

export class MigratePostgresV1 implements Component<MigratePostgresV1Config> {
  readonly uses: string = "nanobus.migrate.postgres/v1";
  readonly with: MigratePostgresV1Config;

  constructor(config: MigratePostgresV1Config) {
    this.with = config;
  }
}
