import { Expose } from "class-transformer";
import { BufferEncoders, RSocketClient } from "rsocket-core";
import RSocketTCPClient from "rsocket-tcp-client";
import { TcpSocketConnectOpts } from "net";
import { Payload, ReactiveSocket } from "rsocket-types";
import { Flowable, Single } from "rsocket-flowable";
import {
  Codec,
  IAdapter,
  IPublisher,
  ISubscriber,
  ISubscription,
  Metadata,
  RequestChannelHandler,
  RequestResponseHandler,
  RequestStreamHandler,
} from "./nanobus";
import { RawItem, Store } from "./stateful";

export class RSocketAdapter implements IAdapter {
  private client: RSocketClient<any, Metadata>;
  private requestResponseHandlers: Map<string, RequestResponseHandler> =
    new Map();
  private requestStreamHandlers: Map<string, RequestStreamHandler> = new Map();
  private requestChannelHandlers: Map<string, RequestChannelHandler> =
    new Map();
  private socket: ReactiveSocket<any, Metadata>;

  constructor(serializer: Codec, options: TcpSocketConnectOpts) {
    const that = this;
    const contentType = serializer.contentType();
    this.client = new RSocketClient({
      serializers: {
        data: serializer,
        metadata: serializer,
      },
      setup: {
        dataMimeType: contentType,
        keepAlive: 60000, // avoid sending during test
        lifetime: 180000,
        metadataMimeType: contentType,
      },
      transport: new RSocketTCPClient(options, BufferEncoders),
      responder: {
        requestResponse(
          payload: Payload<any, Metadata>
        ): Single<Payload<any, Metadata>> {
          const path = payload.metadata[":path"][0];
          const handler = that.requestResponseHandlers.get(path);
          if (!handler) {
            return Single.error(new Error("not_found"));
          }

          return new Single<Payload<any, Metadata>>(async (subscriber) => {
            handler(payload.metadata, payload.data).then(
              (data) => subscriber.onComplete({ data: data }),
              (error) => subscriber.onError(error)
            );
            subscriber.onSubscribe(undefined);
          });
        },
        requestStream<R>(
          payload: Payload<any, Metadata>
        ): Flowable<Payload<any, Metadata>> {
          const path = payload.metadata[":path"][0];
          const handler = that.requestStreamHandlers.get(path);
          if (!handler) {
            return Flowable.error(new Error("not_found"));
          }

          return new Flowable<Payload<any, Metadata>>((subscriber) => {
            const subscriberProxy: ISubscriber<any> = {
              onComplete(): void {
                subscriber.onComplete();
              },
              onError(error: Error): void {
                subscriber.onError(error);
              },
              onSubscribe(subcription: ISubscription) {
                subscriber.onSubscribe(subcription);
              },
              onNext(value: any): void {
                subscriber.onNext({ data: value });
              },
            };
            handler(payload.metadata, payload.data, subscriberProxy);
          });
        },
      },
    });
  }

  registerRequestResponseHandler(
    path: string,
    handler: RequestResponseHandler
  ): void {
    if (handler) {
      this.requestResponseHandlers.set(path, handler);
    }
  }

  registerRequestStreamHandler(
    path: string,
    handler: RequestStreamHandler
  ): void {
    if (handler) {
      this.requestStreamHandlers.set(path, handler);
    }
  }

  registerRequestChannelHandler(
    path: string,
    handler: RequestChannelHandler
  ): void {
    if (handler) {
      this.requestChannelHandlers.set(path, handler);
    }
  }

  async requestResponse<R>(path: string, payload: any): Promise<R> {
    return new Promise((resolve, reject) => {
      this.socket
        .requestResponse({
          data: payload,
          metadata: {
            ":path": [path],
          },
        })
        .then(
          (data) => resolve(data.data as R),
          (error) => reject(error)
        );
    });
  }

  requestStream<R>(path: string, payload?: any): IPublisher<R> {
    return new PublisherImpl(
      this.socket
        .requestStream({
          data: payload,
          metadata: {
            ":path": [path],
          },
        })
        .map((payload) => payload.data)
    );
  }

  requestChannel<I, O>(
    path: string,
    payload: any,
    source: (subscriber: ISubscriber<I>) => void
  ): IPublisher<O> {
    const payloads = new Flowable((subscriber) => {
      var initial = false;
      const subscriberProxy: ISubscriber<any> = {
        onSubscribe(subsciption: ISubscription): void {
          subscriber.onSubscribe({
            cancel() {
              subsciption.cancel();
            },
            request: (n: number) => {
              if (n > 0 && !initial) {
                subscriber.onNext({
                  data: payload,
                  metadata: {
                    ":path": [path],
                  },
                });
                initial = true;
                n--;
              }
              if (n > 0) {
                subsciption.request(n);
              }
            },
          });
        },
        onComplete(): void {
          subscriber.onComplete();
        },
        onError(error: Error): void {
          subscriber.onError(error);
        },
        onNext(value: any): void {
          subscriber.onNext({ data: value });
        },
      };
      source(subscriberProxy);
    });
    return new PublisherImpl<O>(
      this.socket.requestChannel(payloads).map((payload) => payload.data)
    );
  }

  async connect(): Promise<void> {
    const that = this;
    return new Promise((resolve, reject) => {
      that.client.connect().then(
        (socket) => {
          that.socket = socket;
          resolve();
        },
        (error) => reject(error)
      );
    });
  }
}

class PublisherImpl<T> implements IPublisher<T> {
  private flowable: Flowable<T>;

  constructor(flowable: Flowable<T>) {
    this.flowable = flowable;
  }

  subscribe(subscriber?: Partial<ISubscriber<T>>): void {
    return this.flowable.subscribe(subscriber);
  }

  map<R>(fn: (data: T) => R): IPublisher<R> {
    return new PublisherImpl(this.flowable.map(fn));
  }

  async forEach(
    onNextFn: (next: T) => void,
    requestN: number = 100
  ): Promise<void> {
    return new Promise<void>((resolve, reject) => {
      this.subscribe({
        onComplete() {
          resolve();
        },
        onNext(item) {
          onNextFn(item);
        },
        onError(error) {
          reject(error);
        },
        onSubscribe(sub) {
          sub.request(requestN);
        },
      });
    });
  }
}

export class RSocketStorage implements Store {
  private adapter: IAdapter;

  constructor(adapter: IAdapter) {
    this.adapter = adapter;
  }

  async get(namespace: string, id: string, key: string): Promise<RawItem> {
    const inputArgs: GetStateArgs = {
      namespace,
      id,
      key,
    };
    return this.adapter
      .requestResponse<RawItem>("/nanobus.state/get", inputArgs)
      .then();
  }
}

class GetStateArgs {
  @Expose() namespace: string;
  @Expose() id: string;
  @Expose() key: string;

  constructor(namespace: string, id: string, key: string) {
    this.namespace = namespace;
    this.id = id;
    this.key = key;
  }
}
