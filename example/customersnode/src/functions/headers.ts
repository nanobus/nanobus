import { Flag, isFlagSet } from "./flags";
import { Frame, FrameHeader, Type, stripPadding } from "./frame";

const PRIORITY: Flag = 0x20;

export class Headers implements Frame {
  readonly type: Type = Type.HEADERS;
  readonly streamId: number;
  readonly endStream: boolean;
  readonly endHeaders: boolean;
  readonly priority: boolean;
  readonly streamDependencyId: number;
  readonly weight: number;
  readonly exclusive: boolean;
  readonly blockFragment: ArrayBuffer;

  static decode(header: FrameHeader, buffer: ArrayBuffer) {
    const endStream = isFlagSet(header.flags, Flag.END_STREAM);
    const endHeaders = isFlagSet(header.flags, Flag.END_HEADERS);
    const padded = isFlagSet(header.flags, Flag.PADDED);
    const priority = isFlagSet(header.flags, PRIORITY);

    var payload = padded ? stripPadding(buffer) : buffer;
    var streamDependencyId: number = 0;
    var weight: number = 0;
    var exclusive: boolean = false;

    if (priority) {
      const dv = new DataView(payload);
      if (payload.byteLength <= 5) {
        throw new Error(
          "Invalid HEADERS frame: Priority flag set, but payload is too short"
        );
      }
      exclusive = (dv.getUint8(0) & 0x80) === 1;
      streamDependencyId = dv.getUint32(0);
      weight = dv.getUint8(4);
      payload = payload.slice(5);
    }

    return new Headers(
      header.streamId,
      endStream,
      endHeaders,
      priority,
      streamDependencyId,
      weight,
      exclusive,
      payload
    );
  }

  constructor(
    streamId: number,
    endStream: boolean,
    endHeaders: boolean,
    priority: boolean,
    streamDependencyId: number,
    weight: number,
    exclusive: boolean,
    blockFragment: ArrayBuffer
  ) {
    this.streamId = streamId;
    this.endStream = endStream;
    this.endHeaders = endHeaders;
    this.priority = priority;
    this.streamDependencyId = streamDependencyId;
    this.weight = weight;
    this.exclusive = exclusive;
    this.blockFragment = blockFragment;
  }

  size(): number {
    return 9 + this.blockFragment.byteLength + (this.priority ? 4 : 0);
  }

  encode(buffer: ArrayBuffer): void {
    var flags: Flag = 0;
    if (this.endStream) {
      flags |= Flag.END_STREAM;
    }
    if (this.endHeaders) {
      flags |= Flag.END_HEADERS;
    }
    if (this.priority) {
      flags |= PRIORITY;
    }
    const header = new FrameHeader(
      this.blockFragment.byteLength,
      Type.HEADERS,
      flags,
      this.streamId
    );
    header.encode(buffer);
    var data = buffer;
    var start = 9;
    if (this.priority) {
      start += 4;
      const dv = new DataView(buffer);
      var payload = this.streamDependencyId;
      if (this.exclusive) {
        payload |= 0x80000000;
      }
      dv.setUint32(0, payload);
    }

    new Uint8Array(data, start).set(new Uint8Array(this.blockFragment));
  }
}
