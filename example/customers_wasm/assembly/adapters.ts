import { hostCall, register } from "@wapc/as-guest";
import { Decoder, Writer, Encoder, Sizer, Codec } from "@wapc/as-msgpack";
import {
  Address,
  Customer,
  CustomerPage,
  CustomerQuery,
  Inbound,
  Outbound,
  Value,
} from "./interfaces";

export class OutboundImpl implements Outbound {
  saveCustomer(customer: Customer): void {
    hostCall(
      "",
      "customers.v1.Outbound",
      "saveCustomer",
      CustomerCodec.toBuffer(customer)
    );
  }

  fetchCustomer(id: u64): Customer {
    const inputArgs = new OutboundFetchCustomerArgs();
    inputArgs.id = id;
    const payload = hostCall(
      "",
      "customers.v1.Outbound",
      "fetchCustomer",
      inputArgs.toBuffer()
    );
    const decoder = new Decoder(payload);
    return CustomerCodec.decode(decoder);
  }

  customerCreated(customer: Customer): void {
    hostCall(
      "",
      "customers.v1.Outbound",
      "customerCreated",
      CustomerCodec.toBuffer(customer)
    );
  }
}

export var outbound = new OutboundImpl();

var inboundInstance: Inbound;
export function registerInbound(h: Inbound): void {
  inboundInstance = h;
  register(
    "customers.v1.Inbound/createCustomer",
    Inbound_createCustomerWrapper
  );
  register("customers.v1.Inbound/getCustomer", Inbound_getCustomerWrapper);
  register("customers.v1.Inbound/listCustomers", Inbound_listCustomersWrapper);
}

function Inbound_createCustomerWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const request = CustomerCodec.decode(decoder);
  const response = inboundInstance.createCustomer(request);
  return CustomerCodec.toBuffer(response);
}

function Inbound_getCustomerWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const inputArgs = new InboundGetCustomerArgs();
  inputArgs.decode(decoder);
  const response = inboundInstance.getCustomer(inputArgs.id);
  return CustomerCodec.toBuffer(response);
}

function Inbound_listCustomersWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const request = CustomerQueryCodec.decode(decoder);
  const response = inboundInstance.listCustomers(request);
  return CustomerPageCodec.toBuffer(response);
}

class InboundGetCustomerArgs implements Codec {
  id: u64 = 0;

  static decodeNullable(decoder: Decoder): InboundGetCustomerArgs | null {
    if (decoder.isNextNil()) return null;
    return InboundGetCustomerArgs.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): InboundGetCustomerArgs {
    const o = new InboundGetCustomerArgs();
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
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(1);
    encoder.writeString("id");
    encoder.writeUInt64(this.id);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }
}

class InboundSomethingSimpleArgs implements Codec {
  customer: Customer = new Customer();

  static decodeNullable(decoder: Decoder): InboundSomethingSimpleArgs | null {
    if (decoder.isNextNil()) return null;
    return InboundSomethingSimpleArgs.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): InboundSomethingSimpleArgs {
    const o = new InboundSomethingSimpleArgs();
    o.decode(decoder);
    return o;
  }

  decode(decoder: Decoder): void {
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "customer") {
        this.customer = CustomerCodec.decode(decoder);
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(1);
    encoder.writeString("customer");
    CustomerCodec.encode(encoder, this.customer);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }
}

class OutboundFetchCustomerArgs implements Codec {
  id: u64 = 0;

  static decodeNullable(decoder: Decoder): OutboundFetchCustomerArgs | null {
    if (decoder.isNextNil()) return null;
    return OutboundFetchCustomerArgs.decode(decoder);
  }

  // decode
  static decode(decoder: Decoder): OutboundFetchCustomerArgs {
    const o = new OutboundFetchCustomerArgs();
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
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(1);
    encoder.writeString("id");
    encoder.writeUInt64(this.id);
  }

  toBuffer(): ArrayBuffer {
    let sizer = new Sizer();
    this.encode(sizer);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    this.encode(encoder);
    return buffer;
  }
}

///////////////////////////

class CustomerCodec {
  static decodeNullable(decoder: Decoder): Customer | null {
    if (decoder.isNextNil()) return null;
    return CustomerCodec.decode(decoder);
  }

  static decode(decoder: Decoder): Customer {
    const o = new Customer();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        o.id = decoder.readUInt64();
      } else if (field == "firstName") {
        o.firstName = decoder.readString();
      } else if (field == "middleName") {
        if (decoder.isNextNil()) {
          o.middleName = null;
        } else {
          o.middleName = decoder.readString();
        }
      } else if (field == "lastName") {
        o.lastName = decoder.readString();
      } else if (field == "email") {
        o.email = decoder.readString();
      } else if (field == "address") {
        o.address = AddressCodec.decode(decoder);
      } else {
        decoder.skip();
      }
    }

    return o;
  }

  static encode(encoder: Writer, o: Customer): void {
    encoder.writeMapSize(6);
    encoder.writeString("id");
    encoder.writeUInt64(o.id);
    encoder.writeString("firstName");
    encoder.writeString(o.firstName);
    encoder.writeString("middleName");
    if (o.middleName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(o.middleName!);
    }
    encoder.writeString("lastName");
    encoder.writeString(o.lastName);
    encoder.writeString("email");
    encoder.writeString(o.email);
    encoder.writeString("address");
    AddressCodec.encode(encoder, o.address);
  }

  static toBuffer(o: Customer): ArrayBuffer {
    let sizer = new Sizer();
    CustomerCodec.encode(sizer, o);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    CustomerCodec.encode(encoder, o);
    return buffer;
  }
}

class CustomerQueryCodec {
  static decodeNullable(decoder: Decoder): CustomerQuery | null {
    if (decoder.isNextNil()) return null;
    return CustomerQueryCodec.decode(decoder);
  }

  static decode(decoder: Decoder): CustomerQuery {
    const o = new CustomerQuery();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        if (decoder.isNextNil()) {
          o.id = null;
        } else {
          o.id = new Value(decoder.readUInt64());
        }
      } else if (field == "firstName") {
        if (decoder.isNextNil()) {
          o.firstName = null;
        } else {
          o.firstName = decoder.readString();
        }
      } else if (field == "middleName") {
        if (decoder.isNextNil()) {
          o.middleName = null;
        } else {
          o.middleName = decoder.readString();
        }
      } else if (field == "lastName") {
        if (decoder.isNextNil()) {
          o.lastName = null;
        } else {
          o.lastName = decoder.readString();
        }
      } else if (field == "email") {
        if (decoder.isNextNil()) {
          o.email = null;
        } else {
          o.email = decoder.readString();
        }
      } else if (field == "offset") {
        o.offset = decoder.readUInt64();
      } else if (field == "limit") {
        o.limit = decoder.readUInt64();
      } else {
        decoder.skip();
      }
    }

    return o;
  }

  static encode(encoder: Writer, that: CustomerQuery): void {
    encoder.writeMapSize(7);
    encoder.writeString("id");
    if (that.id === null) {
      encoder.writeNil();
    } else {
      encoder.writeUInt64(that.id!.value);
    }
    encoder.writeString("firstName");
    if (that.firstName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(that.firstName!);
    }
    encoder.writeString("middleName");
    if (that.middleName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(that.middleName!);
    }
    encoder.writeString("lastName");
    if (that.lastName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(that.lastName!);
    }
    encoder.writeString("email");
    if (that.email === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(that.email!);
    }
    encoder.writeString("offset");
    encoder.writeUInt64(that.offset);
    encoder.writeString("limit");
    encoder.writeUInt64(that.limit);
  }

  static toBuffer(o: CustomerQuery): ArrayBuffer {
    let sizer = new Sizer();
    CustomerQueryCodec.encode(sizer, o);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    CustomerQueryCodec.encode(encoder, o);
    return buffer;
  }
}

class CustomerPageCodec {
  static decodeNullable(decoder: Decoder): CustomerPage | null {
    if (decoder.isNextNil()) return null;
    return CustomerPageCodec.decode(decoder);
  }

  static decode(decoder: Decoder): CustomerPage {
    const that = new CustomerPage();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "offset") {
        that.offset = decoder.readUInt64();
      } else if (field == "limit") {
        that.limit = decoder.readUInt64();
      } else if (field == "items") {
        that.items = decoder.readArray((decoder: Decoder): Customer => {
          return CustomerCodec.decode(decoder);
        });
      } else {
        decoder.skip();
      }
    }
    return that;
  }

  static encode(encoder: Writer, that: CustomerPage): void {
    encoder.writeMapSize(3);
    encoder.writeString("offset");
    encoder.writeUInt64(that.offset);
    encoder.writeString("limit");
    encoder.writeUInt64(that.limit);
    encoder.writeString("items");
    encoder.writeArray(that.items, (encoder: Writer, item: Customer): void => {
      CustomerCodec.encode(encoder, item);
    });
  }

  static toBuffer(that: CustomerPage): ArrayBuffer {
    let sizer = new Sizer();
    CustomerPageCodec.encode(sizer, that);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    CustomerPageCodec.encode(encoder, that);
    return buffer;
  }
}

class AddressCodec {
  static decodeNullable(decoder: Decoder): Address | null {
    if (decoder.isNextNil()) return null;
    return AddressCodec.decode(decoder);
  }

  static decode(decoder: Decoder): Address {
    const that = new Address();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "line1") {
        that.line1 = decoder.readString();
      } else if (field == "line2") {
        if (decoder.isNextNil()) {
          that.line2 = null;
        } else {
          that.line2 = decoder.readString();
        }
      } else if (field == "city") {
        that.city = decoder.readString();
      } else if (field == "state") {
        that.state = decoder.readString();
      } else if (field == "zip") {
        that.zip = decoder.readString();
      } else {
        decoder.skip();
      }
    }
    return that;
  }

  static encode(encoder: Writer, that: Address): void {
    encoder.writeMapSize(5);
    encoder.writeString("line1");
    encoder.writeString(that.line1);
    encoder.writeString("line2");
    if (that.line2 === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(that.line2!);
    }
    encoder.writeString("city");
    encoder.writeString(that.city);
    encoder.writeString("state");
    encoder.writeString(that.state);
    encoder.writeString("zip");
    encoder.writeString(that.zip);
  }

  static toBuffer(that: Address): ArrayBuffer {
    let sizer = new Sizer();
    AddressCodec.encode(sizer, that);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    AddressCodec.encode(encoder, that);
    return buffer;
  }
}

class ErrorCodec {
  static decodeNullable(decoder: Decoder): Error | null {
    if (decoder.isNextNil()) return null;
    return ErrorCodec.decode(decoder);
  }

  static decode(decoder: Decoder): Error {
    const that = new Error();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "message") {
        that.message = decoder.readString();
      } else {
        decoder.skip();
      }
    }
    return that;
  }

  static encode(encoder: Writer, that: Error): void {
    encoder.writeMapSize(1);
    encoder.writeString("message");
    encoder.writeString(that.message);
  }

  static toBuffer(that: Error): ArrayBuffer {
    let sizer = new Sizer();
    ErrorCodec.encode(sizer, that);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    ErrorCodec.encode(encoder, that);
    return buffer;
  }
}
