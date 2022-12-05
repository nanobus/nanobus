// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export interface ExecConfig {
  // Resource is the name of the connection resource to use.
  resource: ResourceRef;
  // Data is the input bindings sent.
  data?: DataExpr;
  // SQL is the SQL query to execute.
  sql: string;
  // Args are the evaluations to use as arguments for the SQL query.
  args?: ValueExpr[];
}

export class Exec implements Component<ExecConfig> {
  readonly uses: string = "@postgres/exec";
  readonly with: ExecConfig;

  constructor(config: ExecConfig) {
    this.with = config;
  }
}
export interface ExecMultiConfig {
  // Resource is the name of the connection resource to use.
  resource: ResourceRef;
  // Statements are the statements to execute within a single transaction.
  statements: Statement[];
}

export class ExecMulti implements Component<ExecMultiConfig> {
  readonly uses: string = "@postgres/exec_multi";
  readonly with: ExecMultiConfig;

  constructor(config: ExecMultiConfig) {
    this.with = config;
  }
}
export interface Statement {
  // Data is the input bindings sent.
  data?: DataExpr;
  // SQL is the SQL query to execute.
  sql: string;
  // Args are the evaluations to use as arguments for the SQL query.
  args?: ValueExpr[];
}

export interface FindOneConfig {
  // Resource is the name of the connection resource to use.
  resource: ResourceRef;
  // Namespace is the type namespace to load.
  namespace: string;
  // Type is the type name to load.
  type: string;
  // Preload lists the relationship to expand/load.
  preload?: Preload[];
  // Where list the parts of the where clause.
  where?: Where[];
  // NotFoundError is the error to return if the key is not found.
  notFoundError: string;
}

export class FindOne implements Component<FindOneConfig> {
  readonly uses: string = "@postgres/find_one";
  readonly with: FindOneConfig;

  constructor(config: FindOneConfig) {
    this.with = config;
  }
}
export interface Preload {
  field: string;
  preload: Preload[];
}

export interface Where {
  query: string;
  value: ValueExpr;
}

export interface FindConfig {
  // Resource is the name of the connection resource to use.
  resource: ResourceRef;
  // Namespace is the type namespace to load.
  namespace: string;
  // Type is the type name to load.
  type: string;
  // Preload lists the relationship to expand/load.
  preload?: Preload[];
  // Where list the parts of the where clause.
  where?: Where[];
  // Pagination is the optional fields to wrap the results with.
  pagination?: Pagination;
  // Offset is the query offset.
  offset?: ValueExpr;
  // Limit is the query limit.
  limit?: ValueExpr;
}

export class Find implements Component<FindConfig> {
  readonly uses: string = "@postgres/find";
  readonly with: FindConfig;

  constructor(config: FindConfig) {
    this.with = config;
  }
}
export interface Pagination {
  pageIndex?: string;
  pageCount?: string;
  offset?: string;
  limit?: string;
  count?: string;
  total?: string;
  items: string;
}

export interface LoadConfig {
  // Resource is the name of the connection resource to use.
  resource: ResourceRef;
  // Namespace is the type namespace to load.
  namespace: string;
  // Type is the type name to load.
  type: string;
  // ID is the entity identifier expression.
  key: ValueExpr;
  // Preload lists the relationship to expand/load.
  preload?: Preload[];
  // NotFoundError is the error to return if the key is not found.
  notFoundError?: string;
}

export class Load implements Component<LoadConfig> {
  readonly uses: string = "@postgres/load";
  readonly with: LoadConfig;

  constructor(config: LoadConfig) {
    this.with = config;
  }
}
export interface QueryConfig {
  // Resource is the name of the connection resource to use.
  resource: ResourceRef;
  // SQL is the SQL query to execute.
  sql: string;
  // Args are the evaluations to use as arguments for the SQL query.
  args?: ValueExpr[];
  // Single indicates a single row should be returned if found.
  single?: boolean;
}

export class Query implements Component<QueryConfig> {
  readonly uses: string = "@postgres/query";
  readonly with: QueryConfig;

  constructor(config: QueryConfig) {
    this.with = config;
  }
}
