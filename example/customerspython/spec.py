from nanobus import Invoker, invoker, handlers
from typing import Awaitable, Callable, Optional
from dataclasses import dataclass


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
    async def save_customer(self, customer: Customer):
        pass

    async def customer_created(self, customer: Customer):
        pass

    async def fetch_customer(self, id: int) -> Customer:
        pass


class OutboundImpl(Outbound):
    def __init__(self, invoker: Invoker):
        self.invoker = invoker

    async def save_customer(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'saveCustomer', customer)

    async def customer_created(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'customerCreated', customer)

    async def fetch_customer(self, id: int) -> Customer:
        args = GetCustomerArgs(id=id)
        return await self.invoker.invoke_with_return('customers.v1.Outbound', 'fetchCustomer', args, Customer)


def registerInboundHandlers(
    create_customer: Callable[[Customer], Awaitable[Customer]] = None,
    get_customer: Callable[[int], Awaitable[Customer]] = None,
):
    codec = handlers.codec

    if create_customer != None:
        async def handler(input: bytes) -> bytes:
            customer: Customer = codec.decode(input, Customer)
            result = await create_customer(customer)
            return codec.encode(result)
        handlers.register_handler(
            'customers.v1.Inbound', 'createCustomer', handler)

    if get_customer != None:
        async def handler(input: bytes) -> bytes:
            args: GetCustomerArgs = codec.decode(input, GetCustomerArgs)
            result = await get_customer(args.id)
            return codec.encode(result)
        handlers.register_handler(
            'customers.v1.Inbound', 'getCustomer', handler)


outbound = OutboundImpl(invoker)
