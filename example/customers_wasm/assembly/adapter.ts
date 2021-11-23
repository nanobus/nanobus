import { hostCall, register } from "@wapc/as-guest";
import { Decoder, Writer, Encoder, Sizer } from "@wapc/as-msgpack";
import {
  Inbound,
  CustomerActor,
  Outbound,
  Customer,
  CustomerQuery,
  CustomerPage,
  Address,
  Error,
  Value,
} from "./interfaces";

export function registerInbound(handler: Inbound): void {
  InboundHandler = handler;
  register(
    "customers.v1.Inbound/createCustomer",
    Inbound_createCustomerWrapper
  );
  register("customers.v1.Inbound/getCustomer", Inbound_getCustomerWrapper);
  register("customers.v1.Inbound/listCustomers", Inbound_listCustomersWrapper);
}

var InboundHandler: Inbound;

function Inbound_createCustomerWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const request = CustomerCodec.decode(decoder);
  const response = InboundHandler.createCustomer(request);
  return CustomerCodec.toBuffer(response);
}

function Inbound_getCustomerWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const inputArgs = InboundGetCustomerArgsCodec.decode(decoder);
  const response = InboundHandler.getCustomer(inputArgs.id);
  return CustomerCodec.toBuffer(response);
}

function Inbound_listCustomersWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const request = CustomerQueryCodec.decode(decoder);
  const response = InboundHandler.listCustomers(request);
  return CustomerPageCodec.toBuffer(response);
}

class InboundGetCustomerArgs {
  id: u64 = 0;
}

class InboundGetCustomerArgsCodec {
  static decodeNullable(decoder: Decoder): InboundGetCustomerArgs | null {
    if (decoder.isNextNil()) return null;
    return InboundGetCustomerArgsCodec.decode(decoder);
  }

  static decode(decoder: Decoder): InboundGetCustomerArgs {
    const that = new InboundGetCustomerArgs();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        that.id = decoder.readUInt64();
      } else {
        decoder.skip();
      }
    }
    return that;
  }
}

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
      OutboundFetchCustomerArgsCodec.toBuffer(inputArgs)
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

class OutboundFetchCustomerArgs {
  id: u64 = 0;
}

class OutboundFetchCustomerArgsCodec {
  static encode(encoder: Writer, that: OutboundFetchCustomerArgs): void {
    encoder.writeMapSize(1);
    encoder.writeString("id");
    encoder.writeUInt64(that.id);
  }

  static toBuffer(that: OutboundFetchCustomerArgs): ArrayBuffer {
    let sizer = new Sizer();
    OutboundFetchCustomerArgsCodec.encode(sizer, that);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    OutboundFetchCustomerArgsCodec.encode(encoder, that);
    return buffer;
  }
}

class CustomerCodec {
  static decodeNullable(decoder: Decoder): Customer | null {
    if (decoder.isNextNil()) return null;
    return CustomerCodec.decode(decoder);
  }

  static decode(decoder: Decoder): Customer {
    const that = new Customer();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        that.id = decoder.readUInt64();
      } else if (field == "firstName") {
        that.firstName = decoder.readString();
      } else if (field == "middleName") {
        if (decoder.isNextNil()) {
          that.middleName = null;
        } else {
          that.middleName = decoder.readString();
        }
      } else if (field == "lastName") {
        that.lastName = decoder.readString();
      } else if (field == "email") {
        that.email = decoder.readString();
      } else if (field == "address") {
        that.address = AddressCodec.decode(decoder);
      } else {
        decoder.skip();
      }
    }
    return that;
  }

  static encode(encoder: Writer, that: Customer): void {
    encoder.writeMapSize(6);
    encoder.writeString("id");
    encoder.writeUInt64(that.id);
    encoder.writeString("firstName");
    encoder.writeString(that.firstName);
    encoder.writeString("middleName");
    if (that.middleName === null) {
      encoder.writeNil();
    } else {
      encoder.writeString(that.middleName!);
    }
    encoder.writeString("lastName");
    encoder.writeString(that.lastName);
    encoder.writeString("email");
    encoder.writeString(that.email);
    encoder.writeString("address");
    AddressCodec.encode(encoder, that.address);
  }

  static toBuffer(that: Customer): ArrayBuffer {
    let sizer = new Sizer();
    CustomerCodec.encode(sizer, that);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    CustomerCodec.encode(encoder, that);
    return buffer;
  }
}

class CustomerQueryCodec {
  static decodeNullable(decoder: Decoder): CustomerQuery | null {
    if (decoder.isNextNil()) return null;
    return CustomerQueryCodec.decode(decoder);
  }

  static decode(decoder: Decoder): CustomerQuery {
    const that = new CustomerQuery();
    var numFields = decoder.readMapSize();

    while (numFields > 0) {
      numFields--;
      const field = decoder.readString();

      if (field == "id") {
        if (decoder.isNextNil()) {
          that.id = null;
        } else {
          that.id = new Value(decoder.readUInt64());
        }
      } else if (field == "firstName") {
        if (decoder.isNextNil()) {
          that.firstName = null;
        } else {
          that.firstName = decoder.readString();
        }
      } else if (field == "middleName") {
        if (decoder.isNextNil()) {
          that.middleName = null;
        } else {
          that.middleName = decoder.readString();
        }
      } else if (field == "lastName") {
        if (decoder.isNextNil()) {
          that.lastName = null;
        } else {
          that.lastName = decoder.readString();
        }
      } else if (field == "email") {
        if (decoder.isNextNil()) {
          that.email = null;
        } else {
          that.email = decoder.readString();
        }
      } else if (field == "offset") {
        that.offset = decoder.readUInt64();
      } else if (field == "limit") {
        that.limit = decoder.readUInt64();
      } else {
        decoder.skip();
      }
    }
    return that;
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

  static toBuffer(that: CustomerQuery): ArrayBuffer {
    let sizer = new Sizer();
    CustomerQueryCodec.encode(sizer, that);
    let buffer = new ArrayBuffer(sizer.length);
    let encoder = new Encoder(buffer);
    CustomerQueryCodec.encode(encoder, that);
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
        that.items = decoder.readArray(
          (decoder: Decoder): Customer => {
            return CustomerCodec.decode(decoder);
          }
        );
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
