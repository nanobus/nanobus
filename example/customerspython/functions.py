# Python 3 server example
from interfaces import Handlers, Invoker
from requests.api import request
from typing import Callable, Dict, Optional
from dataclasses import asdict, dataclass
from dacite import from_dict


@dataclass
class GetCustomerArgs:
    id: int = 0


@dataclass
class Address:
    line1: str = ''
    line2: Optional[str] = None
    city: str = ''
    state: str = ''
    zip: str = ''


@dataclass
class Customer:
    id: int = 0
    firstName: str = ''
    middleName: Optional[str] = None
    lastName: str = ''
    email: str = ''
    address: Optional[Address] = None


class Outbound:
    def save_customer(self, customer: Customer):
        pass

    def customer_created(self, customer: Customer):
        pass

    def fetch_customer(self, id: int) -> Customer:
        pass


class OutboundImpl(Outbound):
    def __init__(self, invoker: Invoker):
        self.invoker = invoker

    def save_customer(self, customer: Customer):
        self.invoker.invoke("/customers.v1.Outbound/saveCustomer", customer)

    def customer_created(self, customer: Customer):
        self.invoker.invoke("/customers.v1.Outbound/saveCustomer", customer)

    def fetch_customer(self, id: int) -> Customer:
        myargs = GetCustomerArgs(id=id)
        return self.invoker.invokeWithReturn("/customers.v1.Outbound/fetchCustomer", myargs, Customer)


@dataclass
class InboundHandlers:
    create_customer: Callable[[Customer], Customer] = None
    get_customer: Callable[[int], Customer] = None

    def register(self, handlers: Handlers):
        codec = handlers.codec

        if self.create_customer != None:
            def handler(input: bytes) -> bytes:
                customer: Customer = codec.decode(input, Customer)
                result = self.create_customer(customer)
                return codec.encode(result)
            handlers.register_handler(
                "/customers.v1.Inbound/createCustomer", handler)

        if self.get_customer != None:
            def handler(input: bytes) -> bytes:
                args: GetCustomerArgs = codec.decode(input, GetCustomerArgs)
                result = self.get_customer(args.id)
                return codec.encode(result)
            handlers.register_handler(
                "/customers.v1.Inbound/getCustomer", handler)
