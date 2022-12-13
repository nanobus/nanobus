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
  data?: DataExpr;
  to?: string;
}

export function Assign(config: AssignConfig): Component<AssignConfig> {
  return {
    uses: "assign",
    with: config
  };
}

export interface AuthorizeConfig {
  // Condition is the predicate expression for authorization.
  condition?: ValueExpr;
  has?: string[];
  check?: { [key: string]: any };
  error?: string;
}

export function Authorize(config: AuthorizeConfig): Component<AuthorizeConfig> {
  return {
    uses: "authorize",
    with: config
  };
}

export interface CallInterfaceConfig {
  handler: Handler;
}

export function CallInterface(
  config: CallInterfaceConfig
): Component<CallInterfaceConfig> {
  return {
    uses: "call_interface",
    with: config
  };
}

export interface CallProviderConfig {
  handler: Handler;
}

export function CallProvider(
  config: CallProviderConfig
): Component<CallProviderConfig> {
  return {
    uses: "call_provider",
    with: config
  };
}

export interface DecodeConfig {
  typeField: string;
  dataField: string;
  // Codec is the name of the codec to use for decoding.
  codec: string;
  // codecArgs are the arguments to pass to the decode function.
  codecArgs?: any[];
}

export function Decode(config: DecodeConfig): Component<DecodeConfig> {
  return {
    uses: "decode",
    with: config
  };
}

export interface FilterConfig {
  // Condition is the predicate expression for filtering.
  condition: ValueExpr;
}

export function Filter(config: FilterConfig): Component<FilterConfig> {
  return {
    uses: "filter",
    with: config
  };
}

export interface HTTPResponseConfig {
  status?: number;
  headers?: HTTPResponseHeader[];
}

export function HTTPResponse(
  config: HTTPResponseConfig
): Component<HTTPResponseConfig> {
  return {
    uses: "http_response",
    with: config
  };
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
  body?: DataExpr;
  // Metadata is the input binding metadata.
  headers?: DataExpr;
  // Output is an optional transformation to be applied to the response.
  output?: DataExpr;
  // Codec is the name of the codec to use for decoing.
  codec: string;
  // Args are the arguments to pass to the decode function.
  codecArgs?: any[];
}

export function HTTP(config: HTTPConfig): Component<HTTPConfig> {
  return {
    uses: "http",
    with: config
  };
}

export interface InvokeConfig {
  // Name of the interface to invoke.
  interface?: string;
  // Operation of the interface to invoke.
  operation?: string;
  // Input optionally transforms the input sent to the function.
  input?: DataExpr;
}

export function Invoke(config: InvokeConfig): Component<InvokeConfig> {
  return {
    uses: "invoke",
    with: config
  };
}

export interface JMESPathConfig {
  // Path is the predicate expression for filtering.
  path: string;
  // Data is the optional data expression to pass to jq.
  data?: DataExpr;
  // Var, if set, is the variable that is set with the result.
  var?: string;
}

export function JMESPath(config: JMESPathConfig): Component<JMESPathConfig> {
  return {
    uses: "jmespath",
    with: config
  };
}

export interface JQConfig {
  // Query is the predicate expression for filtering.
  query: string;
  // Data is the optional data expression to pass to jq.
  data?: DataExpr;
  // Single, if true, returns the first result.
  single?: boolean;
  // Var, if set, is the variable that is set with the result.
  var?: string;
}

export function JQ(config: JQConfig): Component<JQConfig> {
  return {
    uses: "jq",
    with: config
  };
}

export interface LogConfig {
  format: string;
  // Args are the evaluations to use as arguments into the string format.
  args?: ValueExpr[];
}

export function Log(config: LogConfig): Component<LogConfig> {
  return {
    uses: "log",
    with: config
  };
}

export interface ReCaptchaConfig {
  siteVerifyUrl?: string;
  secret: string;
  response: ValueExpr;
  score?: number;
  action?: string;
}

export function ReCaptcha(config: ReCaptchaConfig): Component<ReCaptchaConfig> {
  return {
    uses: "recaptcha",
    with: config
  };
}

export interface RouteConfig {
  // Selection defines the selection mode: single or multi.
  selection?: SelectionMode;
  // Routes are the possible runnable routes which conditions for selection.
  routes: RouteCondition[];
}

export function Route(config: RouteConfig): Component<RouteConfig> {
  return {
    uses: "route",
    with: config
  };
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
