// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export interface AssignConfig {
  value?: ValueExpr;
  data?: ValueExpr;
  to: string;
}

export class Assign implements Component<AssignConfig> {
  readonly uses: string = "assign";
  readonly with: AssignConfig;

  constructor(config: AssignConfig) {
    this.with = config;
  }
}
export interface AuthorizeConfig {
  // Condition is the predicate expression for authorization.
  condition: ValueExpr;
  has: string[];
  check: { [key: string]: any };
  error?: string;
}

export class Authorize implements Component<AuthorizeConfig> {
  readonly uses: string = "authorize";
  readonly with: AuthorizeConfig;

  constructor(config: AuthorizeConfig) {
    this.with = config;
  }
}
export interface CallInterfaceConfig {
  handler: Handler;
}

export class CallInterface implements Component<CallInterfaceConfig> {
  readonly uses: string = "call_interface";
  readonly with: CallInterfaceConfig;

  constructor(config: CallInterfaceConfig) {
    this.with = config;
  }
}
export interface CallProviderConfig {
  handler: Handler;
}

export class CallProvider implements Component<CallProviderConfig> {
  readonly uses: string = "call_provider";
  readonly with: CallProviderConfig;

  constructor(config: CallProviderConfig) {
    this.with = config;
  }
}
export interface DecodeConfig {
  typeField: string;
  dataField: string;
  // Codec is the name of the codec to use for decoding.
  codec: string;
  // codecArgs are the arguments to pass to the decode function.
  codecArgs?: any[];
}

export class Decode implements Component<DecodeConfig> {
  readonly uses: string = "decode";
  readonly with: DecodeConfig;

  constructor(config: DecodeConfig) {
    this.with = config;
  }
}
export interface FilterConfig {
  // Condition is the predicate expression for filtering.
  condition: ValueExpr;
}

export class Filter implements Component<FilterConfig> {
  readonly uses: string = "filter";
  readonly with: FilterConfig;

  constructor(config: FilterConfig) {
    this.with = config;
  }
}
export interface HTTPResponseConfig {
  status?: number;
  headers?: HTTPResponseHeader[];
}

export class HTTPResponse implements Component<HTTPResponseConfig> {
  readonly uses: string = "http_response";
  readonly with: HTTPResponseConfig;

  constructor(config: HTTPResponseConfig) {
    this.with = config;
  }
}
export interface HTTPResponseHeader {
  name: string;
  value: ValueExpr;
}

export interface HTTPConfig {
  // URL is HTTP URL to request.
  url: string;
  // Method is the HTTP method.
  method: string;
  // Body is the data to sent as the body payload.
  body: DataExpr;
  // Metadata is the input binding metadata.
  headers?: DataExpr;
  // Output is an optional transformation to be applied to the response.
  output?: DataExpr;
  // Codec is the name of the codec to use for decoing.
  codec: string;
  // Args are the arguments to pass to the decode function.
  codecArgs?: any[];
}

export class HTTP implements Component<HTTPConfig> {
  readonly uses: string = "http";
  readonly with: HTTPConfig;

  constructor(config: HTTPConfig) {
    this.with = config;
  }
}
export interface InvokeConfig {
  // Name of the interface to invoke.
  interface: string;
  // Operation of the interface to invoke.
  operation: string;
  // Input optionally transforms the input sent to the function.
  input: DataExpr;
}

export class Invoke implements Component<InvokeConfig> {
  readonly uses: string = "invoke";
  readonly with: InvokeConfig;

  constructor(config: InvokeConfig) {
    this.with = config;
  }
}
export interface JMESPathConfig {
  // Path is the predicate expression for filtering.
  path: string;
  // Data is the optional data expression to pass to jq.
  data: DataExpr;
  // Var, if set, is the variable that is set with the result.
  var?: string;
}

export class JMESPath implements Component<JMESPathConfig> {
  readonly uses: string = "jmespath";
  readonly with: JMESPathConfig;

  constructor(config: JMESPathConfig) {
    this.with = config;
  }
}
export interface JQConfig {
  // Query is the predicate expression for filtering.
  query: string;
  // Data is the optional data expression to pass to jq.
  data?: DataExpr;
  // Single, if true, returns the first result.
  single: boolean;
  // Var, if set, is the variable that is set with the result.
  var?: string;
}

export class JQ implements Component<JQConfig> {
  readonly uses: string = "jq";
  readonly with: JQConfig;

  constructor(config: JQConfig) {
    this.with = config;
  }
}
export interface LogConfig {
  format: string;
  // Args are the evaluations to use as arguments into the string format.
  args?: ValueExpr[];
}

export class Log implements Component<LogConfig> {
  readonly uses: string = "log";
  readonly with: LogConfig;

  constructor(config: LogConfig) {
    this.with = config;
  }
}
export interface ReCaptchaConfig {
  siteVerifyUrl?: string;
  secret: string;
  response: ValueExpr;
  score?: number;
  action?: string;
}

export class ReCaptcha implements Component<ReCaptchaConfig> {
  readonly uses: string = "recaptcha";
  readonly with: ReCaptchaConfig;

  constructor(config: ReCaptchaConfig) {
    this.with = config;
  }
}
export interface RouteConfig {
  // Selection defines the selection mode: single or multi.
  selection: SelectionMode;
  // Routes are the possible runnable routes which conditions for selection.
  routes: RouteCondition[];
}

export class Route implements Component<RouteConfig> {
  readonly uses: string = "route";
  readonly with: RouteConfig;

  constructor(config: RouteConfig) {
    this.with = config;
  }
}
export interface RouteCondition {
  // Name if the overall summary of this route.
  name: string;
  // When is the predicate expression for filtering.
  when: ValueExpr;
  // Then is the steps to process.
  then: Step[];
}

// SelectionMode indicates how many routes can be selected.
export enum SelectionMode {
  Single = "single",
  Multi = "multi"
}
