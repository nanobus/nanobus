import { Framer } from "./framer";
import { Headers } from "./headers";
import { Data } from "./data";
import { queue } from "./queue";
import { Type } from "./frame";

export type StreamHandler = (stream: Stream) => Promise<void>;

export interface Socket {
  send(msg: ArrayBuffer): Promise<void>;
  onData(cb: (msg: ArrayBuffer) => void): void;
}

export const CLIENT_STARTING_STREAM_ID = 1;
export const SERVER_STARTING_STREAM_ID = 2;

export class Connection {
  private socket: Socket;
  private framer: Framer;
  private nextStreamId: number;
  private handler?: StreamHandler;
  private streams: Map<number, Stream> = new Map();

  constructor(socket: Socket, nextStreamId: number, handler?: StreamHandler) {
    this.socket = socket;
    this.framer = new Framer();
    this.nextStreamId = nextStreamId;
    this.handler = handler;

    socket.onData(this.onData.bind(this));
  }

  setHandler(handler: StreamHandler) {
    this.handler = handler;
  }

  onData(buffer: ArrayBuffer) {
    const frame = this.framer.readFrame(buffer);
    switch (frame.type) {
      case Type.HEADERS:
        this.handleHeaders(frame as Headers);
        break;

      case Type.DATA:
        this.handleData(frame as Data);
        break;
    }
  }

  async send(buffer: ArrayBuffer): Promise<void> {
    this.socket.send(buffer);
  }

  handleHeaders(headers: Headers): void {
    const isNew = !this.streams.has(headers.streamId);
    const stream = this.getStream(headers.streamId, true);

    if (headers.blockFragment.byteLength > 0) {
      var enc = new TextDecoder("utf-8");
      stream.md = JSON.parse(enc.decode(headers.blockFragment));
    }

    if (headers.endStream) {
      stream.endStream();
    }

    if (isNew && this.handler !== undefined) {
      this.handler(stream);
    }
  }

  handleData(data: Data): void {
    const stream = this.getStream(data.streamId, true);
    if (stream === undefined) {
      throw new Error(`stream ID ${data.streamId} not found`);
    }
    stream.enqueue(data.data);
    if (data.endStream) {
      stream.endStream();
    }
  }

  getStream(streamId: number, create: boolean): Stream | undefined {
    var stream = this.streams.get(streamId);
    var isNew = false;
    if (stream === undefined && create) {
      stream = new Stream(this, this.framer, streamId);
      this.streams.set(streamId, stream);
      isNew = true;
    }
    return stream;
  }

  removeStream(streamId: number): void {
    this.streams.delete(streamId);
  }

  newStream(parentStreamId: number = 0): Stream {
    const streamID = this.nextStreamId;
    this.nextStreamId += 2;
    const s = new Stream(this, this.framer, streamID, parentStreamId);
    this.streams.set(streamID, s);
    return s;
  }
}

export class Stream {
  private conn: Connection;
  private framer: Framer;

  readonly streamId: number;
  readonly parentStreamId: number;
  private selfClosed: boolean = false;
  private otherClosed: boolean = false;
  private closed: boolean = false;
  md: { [key: string]: string[] };
  private q: queue<ArrayBuffer>;

  constructor(
    conn: Connection,
    framer: Framer,
    streamId: number,
    parentStreamId: number = 0
  ) {
    this.conn = conn;
    this.framer = framer;
    this.streamId = streamId;
    this.parentStreamId = parentStreamId;
    this.q = new queue();
  }

  close(): void {
    if (this.selfClosed) {
      return;
    }
    const headers = new Headers(
      this.streamId,
      true,
      true,
      false,
      0,
      0,
      false,
      new ArrayBuffer(0)
    );
    const frame = this.framer.writeFrame(headers);
    this.conn.send(frame);
    this.selfClosed = true;
    this.checkFullyClosed();
  }

  private checkFullyClosed() {
    if (this.closed) {
      return;
    }
    if (this.selfClosed && this.otherClosed) {
      this.conn.removeStream(this.streamId);
      this.closed = true;
    }
  }

  endStream(): void {
    this.q.close();
    this.otherClosed = true;
    this.checkFullyClosed();
  }

  enqueue(buffer: ArrayBuffer): void {
    this.q.push(buffer);
  }

  sendMetadata(
    metadata: { [key: string]: string[] },
    end: boolean = false
  ): void {
    var blockFragment: ArrayBuffer;
    if (Object.keys(metadata).length > 0) {
      const mdJson = JSON.stringify(metadata);
      blockFragment = toArrayBuffer(Buffer.from(mdJson));
    } else {
      blockFragment = new ArrayBuffer(0);
    }

    const headers = new Headers(
      this.streamId,
      end,
      true,
      false,
      0,
      0,
      false,
      blockFragment
    );
    const frame = this.framer.writeFrame(headers);
    this.conn.send(frame);

    if (end) {
      this.selfClosed = true;
      this.checkFullyClosed();
    }
  }

  sendData(data: ArrayBuffer, end: boolean = false): void {
    const dataFrame = new Data(this.streamId, data, end);
    const frame = this.framer.writeFrame(dataFrame);
    this.conn.send(frame);

    if (end) {
      this.selfClosed = true;
      this.checkFullyClosed();
    }
  }

  sendUnary(metadata: { [key: string]: string[] }, data: ArrayBuffer): void {
    const frameEnd = data === undefined || data.byteLength === 0;
    this.sendMetadata(metadata, frameEnd);
    if (!frameEnd) {
      this.sendData(data, true);
    }
  }

  async receiveData(): Promise<ArrayBuffer | undefined> {
    return this.q.receive();
  }

  async forEach(cb: (ab: ArrayBuffer) => Promise<void>): Promise<void> {
    let buffer = await this.receiveData();
    while (buffer !== undefined) {
      await cb(buffer);
      buffer = await this.receiveData();
    }
  }

  newStream(): Stream {
    return this.conn.newStream(this.streamId);
  }
}

function toArrayBuffer(buf: Buffer | Uint8Array) {
  const ab = new ArrayBuffer(buf.length);
  const view = new Uint8Array(ab);
  view.set(buf);
  return ab;
}
