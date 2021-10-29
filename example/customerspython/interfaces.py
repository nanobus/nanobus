from typing import Optional, Awaitable, Type, TypeVar
from serde import serialize, deserialize
from dataclasses import dataclass, field

T = TypeVar("T")


class LogicalAddress:
    type: str
    id: str


class Context:
    def self() -> LogicalAddress:
        pass

    async def get(self, key: str, type: Type[T]) -> Awaitable[Optional[T]]:
        pass

    def set(self, key: str, data: T):
        pass

    def remove(self, key: str):
        pass


@deserialize
@serialize
@dataclass
# Address information.
class Address:
    # The address line 1
    line1: str = field(default='', metadata={'serde_rename': 'line1'})
    # The address line 2
    line2: Optional[str] = field(default=None,
                                 metadata={'serde_rename': 'line2'})
    # The city
    city: str = field(default='', metadata={'serde_rename': 'city'})
    # The state
    state: str = field(default='', metadata={'serde_rename': 'state'})
    # The zipcode
    zip: str = field(default='', metadata={'serde_rename': 'zip'})


@deserialize
@serialize
@dataclass
class CustomerQuery:
    # The customer identifer
    id: Optional[int] = field(default=None, metadata={'serde_rename': 'id'})
    # The customer's first name
    first_name: Optional[str] = field(default=None,
                                      metadata={'serde_rename': 'firstName'})
    # The customer's middle name
    middle_name: Optional[str] = field(default=None,
                                       metadata={'serde_rename': 'middleName'})
    # The customer's last name
    last_name: Optional[str] = field(default=None,
                                     metadata={'serde_rename': 'lastName'})
    # The customer's email address
    email: Optional[str] = field(default=None,
                                 metadata={'serde_rename': 'email'})
    offset: int = field(default=0, metadata={'serde_rename': 'offset'})
    limit: int = field(default=100, metadata={'serde_rename': 'limit'})


@deserialize
@serialize
@dataclass
# Customer information.
class Customer:
    # The customer identifer
    id: int = field(default=int(0), metadata={'serde_rename': 'id'})
    # The customer's first name
    first_name: str = field(default='', metadata={'serde_rename': 'firstName'})
    # The customer's middle name
    middle_name: Optional[str] = field(default=None,
                                       metadata={'serde_rename': 'middleName'})
    # The customer's last name
    last_name: str = field(default='', metadata={'serde_rename': 'lastName'})
    # The customer's email address
    email: str = field(default='', metadata={'serde_rename': 'email'})
    # The customer's address
    address: Address = field(default=Address(),
                             metadata={'serde_rename': 'address'})


@deserialize
@serialize
@dataclass
class CustomerPage:
    offset: int = field(default=int(0), metadata={'serde_rename': 'offset'})
    limit: int = field(default=int(0), metadata={'serde_rename': 'limit'})
    items: list[Customer] = field(default_factory=list,
                                  metadata={'serde_rename': 'items'})


# Operations that can be performed on a customer.
class Inbound:
    # Creates a new customer.
    async def create_customer(self, customer: Customer) -> Customer:
        pass

    # Retrieve a customer by id.
    async def get_customer(self, id: int) -> Customer:
        pass

    # Return a page of customers using optional search filters.
    async def list_customers(self, query: CustomerQuery) -> CustomerPage:
        pass


# Stateful operations that can be performed on a customer.
class CustomerActor:
    # Creates the customer state.
    async def create_customer(self, ctx: Context,
                              customer: Customer) -> Customer:
        pass

    # Retrieve the customer state.
    async def get_customer(self, ctx: Context) -> Customer:
        pass


class Outbound:
    # Saves a customer to the backend database
    async def save_customer(self, customer: Customer):
        pass

    # Fetches a customer from the backend database
    async def fetch_customer(self, id: int) -> Customer:
        pass

    # Sends a customer creation event
    async def customer_created(self, customer: Customer):
        pass
