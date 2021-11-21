import {
  Decoder,
  Writer,
  Encoder,
  Sizer,
  Codec,
  Value,
} from "@wapc/as-msgpack";

export interface Inbound {
  createCustomer(customer: Customer): Customer;
  getCustomer(id: u64): Customer;
  listCustomers(query: CustomerQuery): CustomerPage;
}

// Customer information.
export class Customer implements Codec {
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

  static decodeNullable(decoder: Decoder): Customer | null {
    if (decoder.isNextNil()) return null;
    return Customer.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): Customer {
    const o = new Customer();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        this.id = decoder.readUInt64();
      } else if (field == "firstName") {
        this.firstName = decoder.readString();
      } else if (field == "middleName") {
        if (decoder.isNextNil()) {
          this.middleName = null;
        } else {
          this.middleName = decoder.readString();
        }
      } else if (field == "lastName") {
        this.lastName = decoder.readString();
      } else if (field == "email") {
        this.email = decoder.readString();
      } else if (field == "address") {
        this.address = Address.decode(decoder);
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(6);
    encoder.writeString("id");
    encoder.writeUInt64(this.id);
    encoder.writeString("firstName");
    encoder.writeString(this.firstName);
    encoder.writeString("middleName");
    if (this.middleName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(this.middleName!);
    }
    encoder.writeString("lastName");
    encoder.writeString(this.lastName);
    encoder.writeString("email");
    encoder.writeString(this.email);
    encoder.writeString("address");
    this.address.encode(encoder);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }

  static newBuilder(): CustomerBuilder {
    return new CustomerBuilder();
  }
}

export class CustomerBuilder {
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

export class CustomerQuery implements Codec {
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

  static decodeNullable(decoder: Decoder): CustomerQuery | null {
    if (decoder.isNextNil()) return null;
    return CustomerQuery.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): CustomerQuery {
    const o = new CustomerQuery();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        if (decoder.isNextNil()) {
          this.id = null;
        } else {
          this.id = new Value(decoder.readUInt64());
        }
      } else if (field == "firstName") {
        if (decoder.isNextNil()) {
          this.firstName = null;
        } else {
          this.firstName = decoder.readString();
        }
      } else if (field == "middleName") {
        if (decoder.isNextNil()) {
          this.middleName = null;
        } else {
          this.middleName = decoder.readString();
        }
      } else if (field == "lastName") {
        if (decoder.isNextNil()) {
          this.lastName = null;
        } else {
          this.lastName = decoder.readString();
        }
      } else if (field == "email") {
        if (decoder.isNextNil()) {
          this.email = null;
        } else {
          this.email = decoder.readString();
        }
      } else if (field == "offset") {
        this.offset = decoder.readUInt64();
      } else if (field == "limit") {
        this.limit = decoder.readUInt64();
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(7);
    encoder.writeString("id");
    if (this.id === null) {
      encoder.writeNil();
    } else {
      encoder.writeUInt64(this.id!.value);
    }
    encoder.writeString("firstName");
    if (this.firstName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(this.firstName!);
    }
    encoder.writeString("middleName");
    if (this.middleName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(this.middleName!);
    }
    encoder.writeString("lastName");
    if (this.lastName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(this.lastName!);
    }
    encoder.writeString("email");
    if (this.email === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(this.email!);
    }
    encoder.writeString("offset");
    encoder.writeUInt64(this.offset);
    encoder.writeString("limit");
    encoder.writeUInt64(this.limit);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }

  static newBuilder(): CustomerQueryBuilder {
    return new CustomerQueryBuilder();
  }
}

export class CustomerQueryBuilder {
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

export class CustomerPage implements Codec {
  offset: u64 = 0;
  limit: u64 = 0;
  items: Array<Customer> = new Array<Customer>();

  static decodeNullable(decoder: Decoder): CustomerPage | null {
    if (decoder.isNextNil()) return null;
    return CustomerPage.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): CustomerPage {
    const o = new CustomerPage();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "offset") {
        this.offset = decoder.readUInt64();
      } else if (field == "limit") {
        this.limit = decoder.readUInt64();
      } else if (field == "items") {
        this.items = decoder.readArray((decoder: Decoder): Customer => {
          return Customer.decode(decoder);
        });
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(3);
    encoder.writeString("offset");
    encoder.writeUInt64(this.offset);
    encoder.writeString("limit");
    encoder.writeUInt64(this.limit);
    encoder.writeString("items");
    encoder.writeArray(this.items, (encoder: Writer, item: Customer): void => {
      item.encode(encoder);
    });
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }

  static newBuilder(): CustomerPageBuilder {
    return new CustomerPageBuilder();
  }
}

export class CustomerPageBuilder {
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

export class Nested implements Codec {
  foo: string = "";
  bar: string = "";

  static decodeNullable(decoder: Decoder): Nested | null {
    if (decoder.isNextNil()) return null;
    return Nested.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): Nested {
    const o = new Nested();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "foo") {
        this.foo = decoder.readString();
      } else if (field == "bar") {
        this.bar = decoder.readString();
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(2);
    encoder.writeString("foo");
    encoder.writeString(this.foo);
    encoder.writeString("bar");
    encoder.writeString(this.bar);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }

  static newBuilder(): NestedBuilder {
    return new NestedBuilder();
  }
}

export class NestedBuilder {
  instance: Nested = new Nested();

  withFoo(foo: string): NestedBuilder {
    this.instance.foo = foo;
    return this;
  }

  withBar(bar: string): NestedBuilder {
    this.instance.bar = bar;
    return this;
  }

  build(): Nested {
    return this.instance;
  }
}

// Address information.
export class Address implements Codec {
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

  static decodeNullable(decoder: Decoder): Address | null {
    if (decoder.isNextNil()) return null;
    return Address.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): Address {
    const o = new Address();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "line1") {
        this.line1 = decoder.readString();
      } else if (field == "line2") {
        if (decoder.isNextNil()) {
          this.line2 = null;
        } else {
          this.line2 = decoder.readString();
        }
      } else if (field == "city") {
        this.city = decoder.readString();
      } else if (field == "state") {
        this.state = decoder.readString();
      } else if (field == "zip") {
        this.zip = decoder.readString();
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(5);
    encoder.writeString("line1");
    encoder.writeString(this.line1);
    encoder.writeString("line2");
    if (this.line2 === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(this.line2!);
    }
    encoder.writeString("city");
    encoder.writeString(this.city);
    encoder.writeString("state");
    encoder.writeString(this.state);
    encoder.writeString("zip");
    encoder.writeString(this.zip);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }

  static newBuilder(): AddressBuilder {
    return new AddressBuilder();
  }
}

export class AddressBuilder {
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
export class Error implements Codec {
  // The detailed error message
  message: string = "";

  static decodeNullable(decoder: Decoder): Error | null {
    if (decoder.isNextNil()) return null;
    return Error.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): Error {
    const o = new Error();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "message") {
        this.message = decoder.readString();
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(1);
    encoder.writeString("message");
    encoder.writeString(this.message);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }

  static newBuilder(): ErrorBuilder {
    return new ErrorBuilder();
  }
}

export class ErrorBuilder {
  instance: Error = new Error();

  withMessage(message: string): ErrorBuilder {
    this.instance.message = message;
    return this;
  }

  build(): Error {
    return this.instance;
  }
}
