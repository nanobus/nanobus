import { Component, ResourceRef } from "../nanobus.ts";
import { ExecConfig, QueryConfig } from "./actions_postgres.ts";
import * as postgres from "./actions_postgres.ts";

export class PostgresActions {
  db: ResourceRef;

  constructor(db: ResourceRef) {
    this.db = db;
  }

  queryOne(sql: string, ...args: unknown[]): Component<QueryConfig> {
    return postgres.Query({
      resource: this.db,
      single: true,
      sql: sql,
      args: args,
    });
  }

  query(sql: string, ...args: unknown[]): Component<QueryConfig> {
    return postgres.Query({
      resource: this.db,
      sql: sql,
      args: args,
    });
  }

  exec(sql: string, ...args: unknown[]): Component<ExecConfig> {
    return postgres.Exec({
      resource: this.db,
      sql: sql,
      args: args,
    });
  }
}
