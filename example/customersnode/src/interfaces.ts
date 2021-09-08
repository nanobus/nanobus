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

export interface Inbound {
  createCustomer?: (customer: Customer) => Promise<Customer>;
  getCustomer?: (id: number) => Promise<Customer>;
}

export interface Outbound {
  saveCustomer(customer: Customer): Promise<void>;
  fetchCustomer(id: number): Promise<Customer>;
  customerCreated(customer: Customer): Promise<void>;
}
