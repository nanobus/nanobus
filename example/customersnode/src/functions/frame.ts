export enum Type {
  DATA = 0x00,
  HEADERS = 0x01,
  PING = 0x06,
}

export interface Frame {
  readonly type: Type;
  readonly streamId: number;
  readonly endStream: boolean;
  size(): number;
  encode(buffer: ArrayBuffer): void;
}

// Flags is a bitwise mask to determine is a flag is set in a frame.
export enum Flag {
  END_STREAM = 0x01,
  END_HEADERS = 0x04,
  PADDED = 0x08,
}

export function isFlagSet(f: Flag, flagsByte: Flag): boolean {
  return (flagsByte & f) !== 0;
}

export function hasFlag(f: Flag, other: Flag): boolean {
  return (other & f) !== 0;
}

export class FrameHeader {
  readonly length: number;
  readonly type: Type;
  readonly flags: Flag;
  readonly streamId: number;

  static decode(buffer: ArrayBuffer) {
    const dv = new DataView(buffer);
    return new FrameHeader(
      uint32IgnoreFirstBit(new Uint8Array(buffer, 0, 3)),
      dv.getUint8(3),
      dv.getUint8(4),
      dv.getUint32(5)
    );
  }

  constructor(length: number, type: Type, flags: Flag, streamId: number) {
    this.length = length;
    this.type = type;
    this.flags = flags;
    this.streamId = streamId;
  }

  encode(buffer: ArrayBuffer): void {
    const buf4 = new ArrayBuffer(4); // implicitly initialized with zero bytes
    const buf = new DataView(buf4);
    buf.setUint32(0, this.length);
    const dv = new DataView(buffer);
    dv.setUint8(0, buf.getUint8(1));
    dv.setUint8(1, buf.getUint8(2));
    dv.setUint8(2, buf.getUint8(3));
    dv.setUint8(3, this.type);
    dv.setUint8(4, this.flags);
    dv.setUint32(5, this.streamId);
  }
}

function uint32IgnoreFirstBit(src: Uint8Array): number {
  const buffer = new Uint8Array(4); // implicitly initialized with zero bytes
  for (var i = 0; i < Math.min(4, src.length); i++) {
    buffer[4 - (1 + i)] = src[4 - (1 + i)];
  }
  buffer[0] &= 0x7f; // clear first bit
  const buf = new DataView(buffer.buffer);
  return buf.getUint32(0);
}

export function stripPadding(payload: ArrayBuffer): ArrayBuffer {
  const dv = new DataView(payload);
  const padLength = dv.getUint8(0);
  if (payload.byteLength <= padLength) {
    throw new Error("Invalid HEADERS: padding >= payload.");
  }
  return payload.slice(1, payload.byteLength - padLength);
}
