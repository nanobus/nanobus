import { encode, decode } from "@msgpack/msgpack";

export type Encodable = Buffer | Uint8Array;
export type Encoder = (v: any) => Encodable;
export type Decoder = (v: Encodable) => any;
export interface Codec {
  serialize: Encoder;
  deserialize: Decoder;
  contentType(): string;
}

export const msgPackSerializer: Codec = {
  deserialize: (data: Encodable) => {
    if (!data) {
      return undefined;
    }
    return decode(data as Buffer);
  },
  serialize: (data: any): Encodable => {
    if (!data) {
      return undefined;
    }
    return Buffer.from(encode(data));
  },
  contentType: () => "application/msgpack",
};

export const jsonSerializer: Codec = {
  deserialize: (data) => {
    return JSON.parse(Buffer.from(data).toString());
  },
  serialize: (data: Encodable) => Buffer.from(JSON.stringify(data)),
  contentType: () => "application/json",
};

export interface IPublisher<T> {
  subscribe(subscriber?: Partial<ISubscriber<T>>): void;
  map<R>(fn: (data: T) => R): IPublisher<R>;
  forEach(onNextFn: (next: T) => void, requestN?: number): void;
}

export type Source<I> = (subscriber: ISubscriber<I>) => void;

export interface ISubscriber<T> {
  onComplete(): void;
  onError(error: Error): void;
  onNext(value: T): void;
  onSubscribe(subscription: ISubscription): void;
}

export interface ISubscription {
  cancel(): void;
  request(n: number): void;
}

export type Metadata = { [key: string]: string[] };

export type RequestResponseHandler = (md: Metadata, payload: any) => Promise<any>;

export type RequestStreamHandler = (
  md: Metadata,
  payload: any,
  subscriber: ISubscriber<any>
) => void;

export type RequestChannelHandler = (
  md: Metadata,
  payload: any,
  subscriber: ISubscriber<any>
) => IPublisher<any>;

export interface IAdapter {
  registerRequestResponseHandler(
    path: string,
    handler: RequestResponseHandler
  ): void;
  registerRequestStreamHandler(
    path: string,
    handler: RequestStreamHandler
  ): void;
  registerRequestChannelHandler(
    path: string,
    handler: RequestChannelHandler
  ): void;

  requestResponse<R>(path: string, payload: any): Promise<R>;
  requestStream<R>(path: string, payload?: any): IPublisher<R>;
  requestChannel<I, O>(
    path: string,
    payload: any,
    source: Source<I>
  ): IPublisher<O>;

  connect(): Promise<void>;
}
