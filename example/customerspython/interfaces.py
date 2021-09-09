from typing import Optional
from serde import serialize, deserialize
from dataclasses import dataclass, field


@deserialize
@serialize
@dataclass
class Address:
    line_1: str = field(default='', metadata={'serde_rename': 'line1'})
    line_2: Optional[str] = field(default=None, metadata={'serde_rename': 'line2'})
    city: str = field(default='', metadata={'serde_rename': 'city'})
    state: str = field(default='', metadata={'serde_rename': 'state'})
    zip: str = field(default='', metadata={'serde_rename': 'zip'})


@deserialize
@serialize
@dataclass
class Customer:
    id: int = field(default=0, metadata={'serde_rename': 'id'})
    first_name: str = field(default='', metadata={'serde_rename': 'firstName'})
    middle_name: Optional[str] = field(default=None, metadata={'serde_rename': 'middleName'})
    last_name: str = field(default='', metadata={'serde_rename': 'lastName'})
    email: str = field(default='', metadata={'serde_rename': 'email'})
    address: Optional[Address] = field(default=None, metadata={'serde_rename': 'address'})


class Outbound:
    async def save_customer(self, customer: Customer):
        pass

    async def customer_created(self, customer: Customer):
        pass

    async def fetch_customer(self, id: int) -> Customer:
        pass
