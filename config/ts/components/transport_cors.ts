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

export interface CorsV0Config {
  // AllowedOrigins is a list of origins a cross-domain request can be executed from.
  // If the special "*" value is present in the list, all origins will be allowed. An
  // origin may contain a wildcard (*) to replace 0 or more characters (i.e.:
  // http://*.domain.com). Usage of wildcards implies a small performance penalty.
  // Only one wildcard can be used per origin. Default value is ["*"]
  allowedOrigins?: string[];
  // AllowedMethods is a list of methods the client is allowed to use with
  // cross-domain requests. Default value is simple methods (HEAD, GET and POST).
  allowedMethods?: string[];
  // AllowedHeaders is list of non simple headers the client is allowed to use with
  // cross-domain requests. If the special "*" value is present in the list, all
  // headers will be allowed. Default value is [] but "Origin" is always appended to
  // the list.
  allowedHeaders?: string[];
  // ExposedHeaders indicates which headers are safe to expose to the API of a CORS
  // API specification
  exposedHeaders?: string[];
  // MaxAge indicates how long (in seconds) the results of a preflight request can be
  // cached
  maxAge?: number;
  // AllowCredentials indicates whether the request can include user credentials like
  // cookies, HTTP authentication or client side SSL certificates.
  allowCredentials: boolean;
  // OptionsPassthrough instructs preflight to let other potential next handlers to
  // process the OPTIONS method. Turn this on if your application handles OPTIONS.
  optionsPassthrough: boolean;
  // Provides a status code to use for successful OPTIONS requests. Default value is
  // http.StatusNoContent (204).
  optionsSuccessStatus: number;
}

export function CorsV0(config: CorsV0Config): Component<CorsV0Config> {
  return {
    uses: "nanobus.transport.http.cors/v0",
    with: config
  };
}
