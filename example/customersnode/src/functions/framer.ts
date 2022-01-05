import { Data } from "./data";
import { Headers } from "./headers";
import { Frame, FrameHeader, Type } from "./frame";

export class Framer {
  private maxReadSize: number = 1 << (24 - 1);

  writeFrame(frame: Frame): ArrayBuffer {
    const size = frame.size();
    const buffer = new ArrayBuffer(size);
    frame.encode(buffer);
    return buffer;
  }

  readFrame(buffer: ArrayBuffer): Frame {
    const frame = FrameHeader.decode(buffer);
    const body = buffer.slice(9);

    if (frame.length > this.maxReadSize) {
      throw new Error("frame too large");
    }

    return typeFrameParser(frame.type)(frame, body);
  }
}

type FrameParser = (fh: FrameHeader, payload: ArrayBuffer) => Frame;

const frameParsers: FrameParser[] = [Data.decode, Headers.decode];
// parseUnknownFrame,
// parseUnknownFrame,
// parseUnknownFrame,
// parseUnknownFrame,
// parsePingFrame,

function typeFrameParser(t: Type): FrameParser {
  if (t >= frameParsers.length) {
    // return parseUnknownFrame;
    throw new Error("TODO");
  }
  return frameParsers[t];
}

// function parsePingFrame(fh: FrameHeader, payload: ArrayBuffer): Frame {
// 	err := fr.pingFrame.Decode(&fh, payload)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &fr.pingFrame, nil
// }

// function parseUnknownFrame(fh: FrameHeader, payload: ArrayBuffer): Frame {
// 	err := fr.rawFrame.Decode(&fh, payload)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &fr.rawFrame, nil
// }
