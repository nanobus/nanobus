import { hostCall, register } from "@wapc/as-guest";
import {
  Decoder,
  Writer,
  Encoder,
  Sizer,
  Codec,
  Value,
} from "@wapc/as-msgpack";
import { Customer, CustomerPage, CustomerQuery } from "./interfaces";

export interface Outbound {
  saveCustomer(customer: Customer): void;
  fetchCustomer(id: u64): Customer;
}

export class OutboundImpl implements Outbound {
  saveCustomer(customer: Customer): void {
    hostCall("", "customers.v1.Outbound", "saveCustomer", customer.toBuffer());
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
    return Customer.decode(decoder);
  }

  customerCreated(customer: Customer): void {
    hostCall(
      "",
      "customers.v1.Outbound",
      "customerCreated",
      customer.toBuffer()
    );
  }
}

export var outbound = new OutboundImpl();

export interface Inbound {
  createCustomer(customer: Customer): Customer;
  getCustomer(id: u64): Customer;
  listCustomers(query: CustomerQuery): CustomerPage;
}

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
  const request = new Customer();
  request.decode(decoder);
  const response = inboundInstance.createCustomer(request);
  return response.toBuffer();
}

function Inbound_getCustomerWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const inputArgs = new InboundGetCustomerArgs();
  inputArgs.decode(decoder);
  const response = inboundInstance.getCustomer(inputArgs.id);
  return response.toBuffer();
}

function Inbound_listCustomersWrapper(payload: ArrayBuffer): ArrayBuffer {
  const decoder = new Decoder(payload);
  const request = new CustomerQuery();
  request.decode(decoder);
  const response = inboundInstance.listCustomers(request);
  return response.toBuffer();
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
        this.customer = Customer.decode(decoder);
      } else {
        decoder.skip();
      }
    }
  }

  encode(encoder: Writer): void {
    encoder.writeMapSize(1);
    encoder.writeString("customer");
    this.customer.encode(encoder);
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
