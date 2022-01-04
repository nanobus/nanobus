import {
  Connection,
  Stream as CStream,
  StreamHandler,
} from "./functions/connection";
import {
  Codec,
  Handler,
  Handlers,
  Invoker,
  StatefulHandler,
  Stream,
} from "./lib/nanobus";

export class NanomsgInvoker implements Invoker {
  private basePath: string;
  private conn: Connection;
  private codec: Codec;

  constructor(basePath: string, conn: Connection, codec: Codec) {
    this.basePath = basePath;
    this.conn = conn;
    this.codec = codec;
  }

  async unary<S, R>(
    namespace: string,
    operation: string,
    payload?: S
  ): Promise<R | undefined> {
    const s = this.conn.newStream();
    const data = payload ? this.codec.encoder(payload) : new ArrayBuffer(0);
    const path = this.basePath + namespace + "/" + operation;

    s.sendUnary({ ":path": [path] }, data);

    const resp = await s.receiveData();
    if (resp === undefined) {
      return; // TODO
    }

    const statusAry = s.md[":status"];
    const status = parseInt(statusAry.length === 1 ? statusAry[0] : "500");

    if (status / 100 !== 2) {
      return; // TODO
    }

    return this.codec.decoder(resp);
  }

  stream<S, R>(
    namespace: string,
    operation: string,
    sendEnd: boolean = false
  ): Stream<S, R> {
    const s = this.conn.newStream();
    const path = this.basePath + namespace + "/" + operation;

    s.sendMetadata({ ":path": [path] }, sendEnd);

    return new NanomsgStream(s, this.codec);
  }
}

export class NanomsgStream<S, R> implements Stream<S, R> {
  private s: CStream;
  private codec: Codec;

  constructor(s: CStream, codec: Codec) {
    this.s = s;
    this.codec = codec;
  }

  async receive(): Promise<R | undefined> {
    const buffer = await this.s.receiveData();
    if (buffer === undefined) {
      return undefined;
    }

    return this.codec.decoder(buffer) as R;
  }

  async forEach(cb: (ab: R) => Promise<void>): Promise<void> {
    let item = await this.receive();
    while (item !== undefined) {
      await cb(item);
      item = await this.receive();
    }
  }

  send(out: S, end: boolean = false): void {
    const buffer = this.codec.encoder(out);
    this.s.sendData(buffer, end);
  }

  end(): void {
    this.s.sendMetadata({}, true);
  }
}

export class NanomsgHandlers implements Handlers {
  readonly codec: Codec;
  private handlers: Map<string, Handler> = new Map();
  private statefulHandlers: Map<string, StatefulHandler> = new Map();
  private conn: Connection;

  constructor(conn: Connection, codec: Codec) {
    this.codec = codec;

    conn.setHandler(this.streamHandler.bind(this));
  }

  private async streamHandler(stream: CStream): Promise<void> {
    const pathAry = stream.md[":path"];
    const path =
      pathAry !== undefined && pathAry.length === 1 ? pathAry[0] : "";
    const parts = path.split("/");
    var handler: Handler = undefined;

    if (parts.length === 3) {
      handler = this.handlers.get(parts[1] + "/" + parts[2]);
    } else if (parts.length === 4) {
      const id = parts[2];
      const shandler = this.statefulHandlers.get(parts[1] + "/" + parts[3]);
      if (shandler) {
        handler = (payload) => {
          return shandler(id, payload);
        };
      }
    }

    if (!handler) {
      stream.sendMetadata({ ":status": ["404"] }, true);
      return;
    }

    var input = (await stream.receiveData()) || new ArrayBuffer(0);

    try {
      const output = await handler(input);
      stream.sendUnary({ ":status": ["200"] }, output);
    } catch (e) {
      console.log(e);
      stream.sendUnary({ ":status": ["500"] }, e.message);
    }
  }

  registerHandler(
    namespace: string,
    operation: string,
    handler: Handler
  ): void {
    if (handler) {
      this.handlers.set(namespace + "/" + operation, handler);
    }
  }

  registerStatefulHandler(
    namespace: string,
    operation: string,
    handler: StatefulHandler
  ): void {
    if (handler) {
      this.statefulHandlers.set(namespace + "/" + operation, handler);
    }
  }
}
