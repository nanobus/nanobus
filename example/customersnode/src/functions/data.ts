import {
  Frame,
  FrameHeader,
  Type,
  Flag,
  isFlagSet,
  stripPadding,
} from "./frame";

export class Data implements Frame {
  readonly type: Type = Type.DATA;
  readonly streamId: number;
  readonly endStream: boolean;
  readonly data: ArrayBuffer;

  static decode(header: FrameHeader, buffer: ArrayBuffer) {
    const padded = isFlagSet(header.flags, Flag.PADDED);
    return new Data(
      header.streamId,
      padded ? stripPadding(buffer) : buffer,
      isFlagSet(header.flags, Flag.END_STREAM)
    );
  }

  constructor(streamId: number, data: ArrayBuffer, endStream: boolean = false) {
    this.streamId = streamId;
    this.data = data;
    this.endStream = endStream;
  }

  size(): number {
    return 9 + this.data.byteLength;
  }

  encode(buffer: ArrayBuffer): void {
    var flags: Flag = this.endStream ? Flag.END_STREAM : 0;
    const header = new FrameHeader(
      this.data.byteLength,
      Type.DATA,
      flags,
      this.streamId
    );
    header.encode(buffer);
    new Uint8Array(buffer, 9).set(new Uint8Array(this.data));
  }
}
