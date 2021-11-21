// Operations that can be performed on a customer.
export interface Inbound {
  // Creates a new customer.
  createCustomer(customer: Customer): Customer;
  // Retrieve a customer by id.
  getCustomer(id: u64): Customer;
  // Return a page of customers using optional search filters.
  listCustomers(query: CustomerQuery): CustomerPage;
}

// Stateful operations that can be performed on a customer.
export interface CustomerActor {
  // Creates the customer state.
  createCustomer(customer: Customer): Customer;
  // Retrieve the customer state.
  getCustomer(): Customer;
}

export interface Outbound {
  saveCustomer(customer: Customer): void;
  fetchCustomer(id: u64): Customer;
  customerCreated(customer: Customer): void;
}

// Customer information.
export class Customer {
  // The customer identifer
  id: u64 = 0;
  // The customer's first name
  firstName: string = "";
  // The customer's middle name
  middleName: string | null = null;
  // The customer's last name
  lastName: string = "";
  // The customer's email address
  email: string = "";
  // The customer's address
  address: Address = new Address();

  static newBuilder(): CustomerBuilder {
    return new CustomerBuilder();
  }
}

class CustomerBuilder {
  instance: Customer = new Customer();

  withId(id: u64): CustomerBuilder {
    this.instance.id = id;
    return this;
  }

  withFirstName(firstName: string): CustomerBuilder {
    this.instance.firstName = firstName;
    return this;
  }

  withMiddleName(middleName: string | null): CustomerBuilder {
    this.instance.middleName = middleName;
    return this;
  }

  withLastName(lastName: string): CustomerBuilder {
    this.instance.lastName = lastName;
    return this;
  }

  withEmail(email: string): CustomerBuilder {
    this.instance.email = email;
    return this;
  }

  withAddress(address: Address): CustomerBuilder {
    this.instance.address = address;
    return this;
  }

  build(): Customer {
    return this.instance;
  }
}

export class CustomerQuery {
  // The customer identifer
  id: Value<u64> | null = null;
  // The customer's first name
  firstName: string | null = null;
  // The customer's middle name
  middleName: string | null = null;
  // The customer's last name
  lastName: string | null = null;
  // The customer's email address
  email: string | null = null;
  offset: u64 = 0;
  limit: u64 = 100;

  static newBuilder(): CustomerQueryBuilder {
    return new CustomerQueryBuilder();
  }
}

class CustomerQueryBuilder {
  instance: CustomerQuery = new CustomerQuery();

  withId(id: Value<u64> | null): CustomerQueryBuilder {
    this.instance.id = id;
    return this;
  }

  withFirstName(firstName: string | null): CustomerQueryBuilder {
    this.instance.firstName = firstName;
    return this;
  }

  withMiddleName(middleName: string | null): CustomerQueryBuilder {
    this.instance.middleName = middleName;
    return this;
  }

  withLastName(lastName: string | null): CustomerQueryBuilder {
    this.instance.lastName = lastName;
    return this;
  }

  withEmail(email: string | null): CustomerQueryBuilder {
    this.instance.email = email;
    return this;
  }

  withOffset(offset: u64): CustomerQueryBuilder {
    this.instance.offset = offset;
    return this;
  }

  withLimit(limit: u64): CustomerQueryBuilder {
    this.instance.limit = limit;
    return this;
  }

  build(): CustomerQuery {
    return this.instance;
  }
}

export class CustomerPage {
  offset: u64 = 0;
  limit: u64 = 0;
  items: Array<Customer> = new Array<Customer>();

  static newBuilder(): CustomerPageBuilder {
    return new CustomerPageBuilder();
  }
}

class CustomerPageBuilder {
  instance: CustomerPage = new CustomerPage();

  withOffset(offset: u64): CustomerPageBuilder {
    this.instance.offset = offset;
    return this;
  }

  withLimit(limit: u64): CustomerPageBuilder {
    this.instance.limit = limit;
    return this;
  }

  withItems(items: Array<Customer>): CustomerPageBuilder {
    this.instance.items = items;
    return this;
  }

  build(): CustomerPage {
    return this.instance;
  }
}

// Address information.
export class Address {
  // The address line 1
  line1: string = "";
  // The address line 2
  line2: string | null = null;
  // The city
  city: string = "";
  // The state
  state: string = "";
  // The zipcode
  zip: string = "";

  static newBuilder(): AddressBuilder {
    return new AddressBuilder();
  }
}

class AddressBuilder {
  instance: Address = new Address();

  withLine1(line1: string): AddressBuilder {
    this.instance.line1 = line1;
    return this;
  }

  withLine2(line2: string | null): AddressBuilder {
    this.instance.line2 = line2;
    return this;
  }

  withCity(city: string): AddressBuilder {
    this.instance.city = city;
    return this;
  }

  withState(state: string): AddressBuilder {
    this.instance.state = state;
    return this;
  }

  withZip(zip: string): AddressBuilder {
    this.instance.zip = zip;
    return this;
  }

  build(): Address {
    return this.instance;
  }
}

// Error response.
export class Error {
  // The detailed error message
  message: string = "";

  static newBuilder(): ErrorBuilder {
    return new ErrorBuilder();
  }
}

class ErrorBuilder {
  instance: Error = new Error();

  withMessage(message: string): ErrorBuilder {
    this.instance.message = message;
    return this;
  }

  build(): Error {
    return this.instance;
  }
}

export class Value<T> {
  value: T;

  constructor(value: T) {
    this.value = value;
  }
}
