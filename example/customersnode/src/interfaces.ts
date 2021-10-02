// Operations that can be performed on a customer.
export interface Inbound {
  // Creates a new customer.
  createCustomer?: (customer: Customer) => Promise<Customer>;
  // Retrieve a customer by id.
  getCustomer?: (id: number) => Promise<Customer>;
  // Return a page of customers using optional search filters.
  listCustomers?: (query: CustomerQuery) => Promise<CustomerPage>;
}

export interface Outbound {
  saveCustomer(customer: Customer): Promise<void>;
  fetchCustomer(id: number): Promise<Customer>;
  customerCreated(customer: Customer): Promise<void>;
}

// Customer information.
export class Customer {
  // The customer identifer
  id: number;
  // The customer's first name
  firstName: string;
  // The customer's middle name
  middleName: string | undefined;
  // The customer's last name
  lastName: string;
  // The customer's email address
  email: string;
  // The customer's address
  address: Address;

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
  id: number | undefined;
  // The customer's first name
  firstName: string | undefined;
  // The customer's middle name
  middleName: string | undefined;
  // The customer's last name
  lastName: string | undefined;
  // The customer's email address
  email: string | undefined;
  offset: number;
  limit: number;

  constructor({
    id = null,
    firstName = null,
    middleName = null,
    lastName = null,
    email = null,
    offset = 0,
    limit = 0,
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
  offset: number;
  limit: number;
  items: Array<Customer>;

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
  foo: string;
  bar: string;

  constructor({ foo = "", bar = "" }: { foo?: string; bar?: string } = {}) {
    this.foo = foo;
    this.bar = bar;
  }
}

// Address information.
export class Address {
  // The address line 1
  line1: string;
  // The address line 2
  line2: string | undefined;
  // The city
  city: string;
  // The state
  state: string;
  // The zipcode
  zip: string;

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
  message: string;

  constructor({ message = "" }: { message?: string } = {}) {
    this.message = message;
  }
}
