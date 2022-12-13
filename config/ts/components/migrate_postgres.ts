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

export function MigratePostgresV1(
  config: MigratePostgresV1Config
): Component<MigratePostgresV1Config> {
  return {
    uses: "nanobus.migrate.postgres/v1",
    with: config
  };
}
