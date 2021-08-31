from typing import Optional
from dataclasses import dataclass


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
    async def save_customer(self, customer: Customer):
        pass

    async def customer_created(self, customer: Customer):
        pass

    async def fetch_customer(self, id: int) -> Customer:
        pass
