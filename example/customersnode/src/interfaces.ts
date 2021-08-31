export interface Customer {
  id: number;
  firstName: string;
  middleName?: string;
  lastName: string;
  email: string;
  address?: Address;
}

export interface Address {
  line1: string;
  line2?: string;
  city: string;
  state: string;
  zip: string;
}

export type createCustomerHandler = (customer: Customer) => Promise<Customer>;
export type getCustomerHandler = (id: number) => Promise<Customer>;

export interface InboundHanders {
  createCustomer?: createCustomerHandler;
  getCustomer?: getCustomerHandler;
}

export interface Outbound {
  saveCustomer(customer: Customer): Promise<void>;
  fetchCustomer(id: number): Promise<Customer>;
  customerCreated(customer: Customer): Promise<void>;
}
