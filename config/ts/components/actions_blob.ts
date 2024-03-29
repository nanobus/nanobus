// Code generated by NanoBus codegen utilities. DO NOT EDIT.

// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
import {
  CodecRef,
  Component,
  DataExpr,
  Handler,
  Entity,
  ResourceRef,
  Step,
  ValueExpr
} from "../nanobus.ts";

export interface ReadConfig {
  // The blob store resource to read.
  resource: ResourceRef;
  // The key to read.
  key: ValueExpr;
  // Codec is the name of the codec to use for decoding.
  codec?: string;
  // codecArgs are the arguments to pass to the decode function.
  codecArgs?: any[];
  offset?: ValueExpr;
  length?: ValueExpr;
  bufferSize?: number;
}

export function Read(config: ReadConfig): Component<ReadConfig> {
  return {
    uses: "@blob/read",
    with: config
  };
}

export interface WriteConfig {
  // The blob store resource to write.
  resource: ResourceRef;
  // The key to write.
  key: ValueExpr;
  // The data to write.
  data?: DataExpr;
  // Codec is the name of the codec to use for decoding.
  codec?: string;
  // codecArgs are the arguments to pass to the decode function.
  codecArgs?: any[];
  delimiterString?: string;
  delimiterBytes?: ArrayBuffer;
}

export function Write(config: WriteConfig): Component<WriteConfig> {
  return {
    uses: "@blob/write",
    with: config
  };
}
