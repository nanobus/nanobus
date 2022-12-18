import { Component, Handler } from "../nanobus.ts";
import {
  CallInterface,
  CallInterfaceConfig,
  CallProvider,
  CallProviderConfig,
  Log,
  LogConfig,
} from "./actions_core.ts";

export function log(format: string, ...args: unknown[]): Component<LogConfig> {
  return Log({
    format: format,
    args: args,
  });
}

export function callInterface(
  handler: Handler,
): Component<CallInterfaceConfig> {
  return CallInterface({
    handler,
  });
}

export function callProvider(handler: Handler): Component<CallProviderConfig> {
  return CallProvider({
    handler,
  });
}
