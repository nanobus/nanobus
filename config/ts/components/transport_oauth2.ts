// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  Component,
  DataExpr,
  Handler,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export interface OAuth2V1Config {
  loginPath?: string;
  callbackPath?: string;
  clientId: string;
  clientSecret: string;
  endpoint: Endpoint;
  redirectUrl: string;
  scopes?: string[];
  handler?: Handler;
}

export function OAuth2V1(config: OAuth2V1Config): Component<OAuth2V1Config> {
  return {
    uses: "nanobus.transport.http.oauth2/v1",
    with: config
  };
}

export interface Endpoint {
  authUrl: string;
  tokenUrl: string;
  userInfoUrl: string;
  // AuthStyle optionally specifies how the endpoint wants the client ID & client
  // secret sent.
  authStyle?: AuthStyle;
}

// AuthStyle represents how requests for tokens are authenticated to the server.
export enum AuthStyle {
  // AuthStyleAutoDetect means to auto-detect which authentication style the
  // provider wants by trying both ways and caching the successful way for the
  // future.
  AutoDetect = "auto-detect",
  // AuthStyleInParams sends the "client_id" and "client_secret" in the POST body as
  // application/x-www-form-urlencoded parameters.
  InParams = "inparams",
  // AuthStyleInHeader sends the client_id and client_password using HTTP Basic
  // Authorization. This is an optional style described in the OAuth2 RFC 6749
  // section 2.3.1.
  InHeader = "inheader"
}
