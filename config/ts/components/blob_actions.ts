import { Component, ResourceRef, ValueExpr } from "../nanobus.ts";
import * as blob from "./actions_blob.ts";

export class BlobActions {
  bucket: ResourceRef;

  constructor(bucket: ResourceRef) {
    this.bucket = bucket;
  }

  read(
    key: ValueExpr,
    codec: string,
    ...codecArgs: unknown[]
  ): Component<blob.ReadConfig> {
    return blob.Read({
      resource: this.bucket,
      key,
      codec,
      codecArgs,
    });
  }

  write(
    key: ValueExpr,
    codec: string,
    options: Omit<blob.WriteConfig, "resource" | "key" | "codec">,
  ): Component<blob.WriteConfig> {
    return blob.Write({
      resource: this.bucket,
      key,
      codec,
      ...options,
    });
  }
}
