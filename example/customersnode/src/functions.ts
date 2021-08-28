import url from 'url';
import http, { RequestListener, Server } from 'http';
import { encode, decode } from '@msgpack/msgpack';

export type Encoder = (v: any) => ArrayBuffer
export type Decoder = (v: ArrayBuffer) => any;
export interface Codec {
  encoder: Encoder
  decoder: Decoder
}

export type Handler = (payload: ArrayBuffer) => Promise<ArrayBuffer>

export interface Handlers {
  registerHandler(uri: string, handler: Handler): void
  registerHandlers(handlers: { [uri: string]: Handler }): void
}

export class HTTPHandlers implements Handlers {
  private handlers: Map<string, Handler> = new Map();
  private server: Server;

  constructor() {
    const requestListener: RequestListener = async function (req, res) {
      const handler: Handler = this.handlers.get(req.url);
      if (!handler) {
        res.writeHead(404);
        res.end('Not found');
        return;
      }

      try {
        const buffers = [];
        for await (const chunk of req) {
          buffers.push(chunk);
        }
        const input = Buffer.concat(buffers);

        const output = await handler(input);

        var responseBuffer: Buffer;
        if (output.byteLength > 0) {
          responseBuffer = Buffer.from(output);
        } else {
          responseBuffer = Buffer.alloc(0);
        }

        res.setHeader('Content-Type', 'application/msgpack');
        res.writeHead(200);
        res.end(responseBuffer);
      } catch (e) {
        console.log(e);
        res.writeHead(500);
        res.end(e.message);
      }
    }

    this.server = http.createServer(requestListener.bind(this));
  }

  registerHandler(uri: string, handler: Handler): void {
    if (handler) {
      this.handlers.set(uri, handler);
    }
  }

  registerHandlers(handlers: { [uri: string]: Handler }): void {
    Object.keys(handlers).map(uri => this.registerHandler(uri, handlers[uri]));
  }

  listen(port?: number, hostname?: string, listeningListener?: () => void): void {
    this.server.listen(port, hostname, listeningListener);
  }
}

export type Invoker = (operation: string, payload: any) => Promise<any>;

export function HTTPInvoker(baseURL: string, codec: Codec): Invoker {
  const u = url.parse(baseURL);

  return async (operation: string, payload: any): Promise<any> => {
    return new Promise((resolve, reject) => {
      const data = codec.encoder(payload);
      const options = {
        hostname: u.hostname,
        port: u.port,
        path: u.path + operation,
        method: 'POST',
        headers: {
          'Content-Type': 'application/msgpack',
          'Content-Length': data.byteLength
        }
      }

      const req = http.request(options, res => {
        const buffers: Uint8Array[] = [];
        res.on('data', chunk => {
          buffers.push(chunk);
        })

        res.on('end', () => {
          try {
            if (buffers.length == 0) {
              resolve(null);
              return;
            }

            const data = Buffer.concat(buffers);
            const resp = codec.decoder(data);
            resolve(resp);
          } catch (error) {
            console.error(error);
            reject(error);
          }
        })
      })

      req.on('error', error => {
        console.error(error);
        reject(error);
      })

      req.write(Buffer.from(data));
      req.end();
    });
  }
};

export const msgpackCodec: Codec = {
  encoder: (data) => encode(data).buffer,
  decoder: (data) => decode(data)
}