import { Expose } from "class-transformer";

export interface LogicalAddress {
  readonly type: string;
  readonly id: string;
  toString(): string;
}

export declare type ClassConstructor<T> = {
  new (...args: any[]): T;
};

export interface Context {
  readonly self: LogicalAddress;
  get<T>(key: string, cls: ClassConstructor<T>): Promise<T | undefined>;
  set<T>(key: string, data: T): void;
  remove(key: string): void;
}

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

// Operations that can be performed on a customer.
export interface Inbound {
  // Creates a new customer.
  createCustomer(customer: Customer): Promise<Customer>;
  // Retrieve a customer by id.
  getCustomer(id: number): Promise<Customer>;
  // Return a page of customers using optional search filters.
  listCustomers(query: CustomerQuery): Promise<CustomerPage>;
}

// Stateful operations that can be performed on a customer.
export interface CustomerActor {
  // Creates the customer state.
  createCustomer(ctx: Context, customer: Customer): Promise<Customer>;
  // Retrieve the customer state.
  getCustomer(ctx: Context): Promise<Customer>;
}

export interface Outbound {
  // Saves a customer to the backend database
  saveCustomer(customer: Customer): Promise<void>;
  // Fetches a customer from the backend database
  fetchCustomer(id: number): Promise<Customer>;
  // Sends a customer creation event
  customerCreated(customer: Customer): Promise<void>;
  // Get customers from the database
  getCustomers(): IPublisher<Customer>;
  // Transform customers
  transformCustomers(
    prefix: string,
    source: Source<Customer>
  ): IPublisher<Customer>;
}

// Customer information.
export class Customer {
  // The customer identifer
  @Expose() id: number;
  // The customer's first name
  @Expose() firstName: string;
  // The customer's middle name
  @Expose() middleName: string | undefined;
  // The customer's last name
  @Expose() lastName: string;
  // The customer's email address
  @Expose() email: string;
  // The customer's address
  @Expose() address: Address;

  constructor({
    id = 0,
    firstName = "",
    middleName = null,
    lastName = "",
    email = "",
    address = new Address(),
  }: {
    id?: number;
    firstName?: string;
    middleName?: string | undefined;
    lastName?: string;
    email?: string;
    address?: Address;
  } = {}) {
    this.id = id;
    this.firstName = firstName;
    this.middleName = middleName;
    this.lastName = lastName;
    this.email = email;
    this.address = address;
  }
}

export class CustomerQuery {
  // The customer identifer
  @Expose() id: number | undefined;
  // The customer's first name
  @Expose() firstName: string | undefined;
  // The customer's middle name
  @Expose() middleName: string | undefined;
  // The customer's last name
  @Expose() lastName: string | undefined;
  // The customer's email address
  @Expose() email: string | undefined;
  @Expose() offset: number;
  @Expose() limit: number;

  constructor({
    id = null,
    firstName = null,
    middleName = null,
    lastName = null,
    email = null,
    offset = 0,
    limit = 100,
  }: {
    id?: number | undefined;
    firstName?: string | undefined;
    middleName?: string | undefined;
    lastName?: string | undefined;
    email?: string | undefined;
    offset?: number;
    limit?: number;
  } = {}) {
    this.id = id;
    this.firstName = firstName;
    this.middleName = middleName;
    this.lastName = lastName;
    this.email = email;
    this.offset = offset;
    this.limit = limit;
  }
}

export class CustomerPage {
  @Expose() offset: number;
  @Expose() limit: number;
  @Expose() items: Array<Customer>;

  constructor({
    offset = 0,
    limit = 0,
    items = new Array<Customer>(),
  }: { offset?: number; limit?: number; items?: Array<Customer> } = {}) {
    this.offset = offset;
    this.limit = limit;
    this.items = items;
  }
}

export class Nested {
  @Expose() foo: string;
  @Expose() bar: string;

  constructor({ foo = "", bar = "" }: { foo?: string; bar?: string } = {}) {
    this.foo = foo;
    this.bar = bar;
  }
}

// Address information.
export class Address {
  // The address line 1
  @Expose() line1: string;
  // The address line 2
  @Expose() line2: string | undefined;
  // The city
  @Expose() city: string;
  // The state
  @Expose() state: string;
  // The zipcode
  @Expose() zip: string;

  constructor({
    line1 = "",
    line2 = null,
    city = "",
    state = "",
    zip = "",
  }: {
    line1?: string;
    line2?: string | undefined;
    city?: string;
    state?: string;
    zip?: string;
  } = {}) {
    this.line1 = line1;
    this.line2 = line2;
    this.city = city;
    this.state = state;
    this.zip = zip;
  }
}

// Error response.
export class Error {
  // The detailed error message
  @Expose() message: string;

  constructor({ message = "" }: { message?: string } = {}) {
    this.message = message;
  }
}
