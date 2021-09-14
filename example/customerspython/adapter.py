import os
from nanobus import AIOHTTPServer, UvicornServer, HTTPInvoker, Handlers, Invoker, MsgPackCodec
from typing import Awaitable, Callable
from serde import serialize, deserialize
from dataclasses import dataclass, field
from interfaces import Customer, Outbound


server_host = os.getenv('HOST', "localhost")
server_port = int(os.getenv('PORT', "9000"))
outbound_base_url = os.getenv(
    'OUTBOUND_BASE_URL', "http://localhost:32321/outbound")

codec = MsgPackCodec()
handlers = Handlers(codec)
#server = AIOHTTPServer(handlers)
server = UvicornServer(handlers)
http_invoker = HTTPInvoker(outbound_base_url)
invoker = Invoker(http_invoker.invoke, codec)


@deserialize
@serialize
@dataclass
class _GetCustomerArgs:
    id: int = field(metadata={'serde_rename': 'id'})


class OutboundImpl(Outbound):
    def __init__(self, invoker: Invoker):
        self.invoker = invoker

    async def save_customer(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'saveCustomer', customer)

    async def customer_created(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'customerCreated', customer)

    async def fetch_customer(self, id: int) -> Customer:
        args = _GetCustomerArgs(
            id,
        )
        return await self.invoker.invoke_with_return('customers.v1.Outbound', 'fetchCustomer', args, Customer)


outbound = OutboundImpl(invoker)


def register_inbound_handlers(
    create_customer: Callable[[Customer], Awaitable[Customer]] = None,
    get_customer: Callable[[int], Awaitable[Customer]] = None,
):
    if not create_customer is None:
        async def handler(input: bytes) -> bytes:
            customer: Customer = handlers.codec.decode(input, Customer)
            result = await create_customer(customer)
            return handlers.codec.encode(result)
        handlers.register_handler(
            'customers.v1.Inbound', 'createCustomer', handler)

    if not get_customer is None:
        async def handler(input: bytes) -> bytes:
            args: _GetCustomerArgs = handlers.codec.decode(
                input, _GetCustomerArgs)
            result = await get_customer(args.id)
            return handlers.codec.encode(result)
        handlers.register_handler(
            'customers.v1.Inbound', 'getCustomer', handler)


def start():
    try:
        server.run(server_host, server_port)
    except KeyboardInterrupt:
        pass

    print("Server stopped.")
