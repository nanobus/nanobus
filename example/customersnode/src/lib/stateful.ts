import url, { UrlWithStringQuery } from "url";
import http from "http";
import { LRUCache as LRU } from "typescript-lru-cache";
import { Codec, StatefulHandler } from "./nanobus";

export class LogicalAddress {
  type: string;
  id: string;

  constructor(type: string, id: string) {
    this.type = type;
    this.id = id;
  }

  toString(): string {
    return `${this.type}/${this.id}`;
  }
}

export class Response {
  mutation: Mutation;
  result: any;

  constructor(mutation: Mutation, result: any) {
    this.mutation = mutation;
    this.result = result;
  }
}

export class Context {
  // context.Context;
  readonly self: LogicalAddress;
  private state: State;

  constructor(self: LogicalAddress, state: State) {
    this.self = self;
    this.state = state;
  }

  async get<T>(key: string): Promise<T | undefined> {
    return this.state.get(key);
  }

  async getItem<T>(key: string): Promise<Item<T> | undefined> {
    return this.state.getItem(key);
  }

  set<T>(key: string, data: T) {
    this.state.set(key, data);
  }

  setItem<T>(key: string, item: Item<T>) {
    this.state.setItem(key, item);
  }

  remove(key: string) {
    this.state.remove(key);
  }

  response(result: any): Response {
    const mutation = this.state.toMutation();
    return new Response(mutation, result);
  }
}

export class Item<T> {
  namespace?: string;
  type?: string;
  version?: string;
  data: T;

  constructor({
    namespace = undefined,
    type = undefined,
    version = undefined,
    data,
  }: { namespace?: string; type?: string; version?: string; data?: any } = {}) {
    this.namespace = namespace;
    this.type = type;
    this.version = version;
    this.data = data;
  }
}

export class RawItem {
  namespace?: string;
  type?: string;
  version?: string;
  data: ArrayBuffer;

  constructor({
    namespace = undefined,
    type = undefined,
    version = undefined,
    data,
  }: {
    namespace?: string;
    type?: string;
    version?: string;
    data?: ArrayBuffer;
  } = {}) {
    this.namespace = namespace;
    this.type = type;
    this.version = version;
    this.data = data;
  }
}

export interface Namespacer {
  getNamespace(): string;
}

export interface Typer {
  getType(): string;
}

export interface Versioner {
  getVersion(): string;
}

export interface Activator {
  activate(ctx: Context): Promise<void>;
}

export interface Deactivator {
  deactivate(ctx: Context): Promise<void>;
}

interface Cache {
  get(address: LogicalAddress): CachedState | undefined;
  set(address: LogicalAddress, state: CachedState): void;
  remove(address: LogicalAddress): void;
}

export class LRUCache implements Cache {
  private lru: LRU<string, CachedState>;

  constructor() {
    this.lru = new LRU();
  }

  get(address: LogicalAddress): CachedState | undefined {
    return this.lru.get(address.toString());
  }

  set(address: LogicalAddress, state: CachedState): void {
    this.lru.set(address.toString(), state);
  }

  remove(address: LogicalAddress): void {
    this.lru.delete(address.toString());
  }
}

export class CachedState {
  revision: number;
  state: State;
  actor: any;

  constructor(revision: number, state: State, actor: any) {
    this.revision = revision;
    this.state = state;
    this.actor = actor;
  }
}

export interface Store {
  get(namespace: string, id: string, key: string): Promise<RawItem | undefined>;
}

export class Storage implements Store {
  private u: UrlWithStringQuery;
  private codec: Codec;

  constructor(baseUrl: string, codec: Codec) {
    this.u = url.parse(baseUrl);
    this.codec = codec;
  }

  async get(namespace: string, id: string, key: string): Promise<RawItem> {
    return new Promise((resolve, reject) => {
      const options = {
        hostname: this.u.hostname,
        port: this.u.port,
        path: this.u.path + "state/" + namespace + "/" + id + "/" + key,
        method: "GET",
      };

      const req = http.request(options, (res) => {
        const buffers: Uint8Array[] = [];
        res.on("data", (chunk) => {
          buffers.push(chunk);
        });

        res.on("end", () => {
          try {
            if (buffers.length === 0) {
              resolve(null);
              return;
            }

            const data = Buffer.concat(buffers);
            const item: any = this.codec.decoder(data);
            if (
              item.data &&
              (typeof item.data === "string" || item.data instanceof String)
            ) {
              item.data = Buffer.from(
                Buffer.from(item.data).toString(),
                "base64"
              );
            }
            resolve(item as RawItem);
          } catch (error) {
            console.error(error);
            reject(error);
          }
        });
      });

      req.on("error", (error) => {
        console.error(error);
        reject(error);
      });

      req.end();
    });
  }
}

export interface Mutation {
  set: { [key: string]: RawItem };
  remove: string[];
}

export class State {
  readonly address: LogicalAddress;
  private store: Store;
  private codec: Codec;
  private state: Map<string, Item<any>>;
  private changed: Map<string, Item<any>>;
  private removed: Set<string>;

  constructor(address: LogicalAddress, store: Store, codec: Codec) {
    this.address = address;
    this.store = store;
    this.codec = codec;
    this.state = new Map();
    this.changed = null;
    this.removed = null;
  }

  clear(): void {
    this.state = new Map();
    this.changed = null;
    this.removed = null;
  }

  async get<T>(key: string): Promise<T | undefined> {
    const item = await this.getItem<T>(key);
    if (item) {
      return item.data;
    }
    return undefined;
  }

  async getItem<T>(key: string): Promise<Item<T> | undefined> {
    if (this.removed && this.removed.has(key)) {
      return undefined;
    }

    var item: Item<T>;
    if (this.changed) {
      item = this.changed.get(key);
    }
    if (!item) {
      item = this.state.get(key);
    }
    if (item) {
      return item;
    }

    const rawItem = await this.store.get(
      this.address.type,
      this.address.id,
      key
    );
    if (!rawItem) {
      return undefined;
    }

    const data = this.codec.decoder(rawItem.data);
    item = new Item({
      namespace: rawItem.namespace,
      type: rawItem.type,
      version: rawItem.version,
      data: data,
    });
    this.state.set(key, item);

    return item;
  }

  set(key: string, data: any): void {
    this.setItem(key, dataToItem(data));
  }

  setItem(key: string, item: Item<any>): void {
    if (!this.changed) {
      this.changed = new Map();
    }
    this.changed.set(key, item);
    if (this.removed) {
      this.removed.delete(key);
    }
  }

  remove(key: string): void {
    if (!this.changed) {
      this.changed = new Map();
    }
    this.changed.set(key, new Item());
    if (!this.removed) {
      this.removed = new Set();
    }
    this.removed.add(key);
  }

  isDirty(): boolean {
    return this.changed != null || this.removed != null;
  }

  reset(): void {
    if (this.changed) {
      this.changed.forEach((item, key) => this.state.set(key, item));
    }
    this.changed = null;
    this.removed = null;
  }

  toMutation(): Mutation {
    var set: { [key: string]: RawItem } = undefined;
    if (this.changed) {
      set = {};
      this.changed.forEach((item, key) => {
        // Skip deleted items;
        if (!item.data) {
          return;
        }
        const data = this.codec.encoder(item.data);
        set[key] = new RawItem({
          namespace: item.namespace,
          type: item.type,
          version: item.version,
          data: data,
        });
      });
    }

    var remove: string[] = undefined;
    if (this.removed) {
      remove = [];
      this.removed.forEach((removed) => remove.push(removed));
    }

    return { set: set, remove: remove };
  }
}

function dataToItem<T>(data: T): Item<T> {
  var namespace: string = undefined;
  var type: string = undefined;
  var version: string = undefined;

  if (isNamespacer(data)) {
    namespace = (data as Namespacer).getNamespace();
  }
  if (isTyper(data)) {
    type = (data as Typer).getType();
  } else {
    type = data.constructor.name;
  }
  if (isVersioner(data)) {
    version = (data as Versioner).getVersion();
  }

  return new Item({
    namespace: namespace,
    type: type,
    version: version,
    data: data,
  });
}

function isNamespacer(object: any): object is Namespacer {
  return (<Namespacer>object).getNamespace !== undefined;
}

function isTyper(object: any): object is Typer {
  return (<Typer>object).getType !== undefined;
}

function isVersioner(object: any): object is Versioner {
  return (<Versioner>object).getVersion !== undefined;
}

export class Manager {
  private cache: Cache;
  private store: Store;
  private codec: Codec;

  constructor(cache: Cache, store: Store, codec: Codec) {
    this.cache = cache;
    this.store = store;
    this.codec = codec;
  }

  toContext(type: string, id: string, actor: any): Context {
    const address = new LogicalAddress(type, id);
    const state = this.getState(address, 0, actor);

    return new Context(address, state);
  }

  getState(address: LogicalAddress, revision: number, actor: any): State {
    if (this.cache) {
      const cached = this.cache.get(address);
      if (cached) {
        if (revision !== 0 && revision !== cached.revision) {
          cached.state.clear();
        } else {
          cached.revision++;
        }

        return cached.state;
      }
    }

    const state = new State(address, this.store, this.codec);
    const cachedState = new CachedState(revision + 1, state, actor);

    if (this.cache) {
      this.cache.set(address, cachedState);
    }

    if (isActivator(actor)) {
      const sctx = new Context(address, state);
      (actor as Activator).activate(sctx);
    }

    return state;
  }

  deactivate(address: LogicalAddress) {
    if (this.cache) {
      this.cache.remove(address);
    }
  }

  deactivateHandler(type: string, actor: any): StatefulHandler {
    return async (id, _) => {
      const sctx = this.toContext(type, id, actor);
      if (isDeactivator(actor)) {
        (actor as Deactivator).deactivate(sctx);
      }
      this.deactivate(new LogicalAddress(type, id));
      return new ArrayBuffer(0);
    };
  }
}

function isActivator(object: any): object is Activator {
  return (<Activator>object).activate !== undefined;
}

function isDeactivator(object: any): object is Deactivator {
  return (<Deactivator>object).deactivate !== undefined;
}
