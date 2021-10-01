import os
from nanobus import UvicornServer, HTTPInvoker, Handlers, Invoker, MsgPackCodec
from serde import serialize, deserialize
from dataclasses import dataclass, field
from interfaces import Inbound, Customer, CustomerPage, CustomerQuery, Outbound

server_host = os.getenv('HOST', "localhost")
server_port = int(os.getenv('PORT', "9000"))
outbound_base_url = os.getenv('OUTBOUND_BASE_URL',
                              "http://localhost:32321/outbound")

codec = MsgPackCodec()
handlers = Handlers(codec)
server = UvicornServer(handlers)
http_invoker = HTTPInvoker(outbound_base_url)
invoker = Invoker(http_invoker.invoke, codec)


@deserialize
@serialize
@dataclass
class _InboundGetCustomerArgs:
    id: int = field(default=int(0), metadata={'serde_rename': 'id'})


def register_inbound_handlers(h: Inbound):
    if not h.create_customer is None:

        async def handler(input: bytes) -> bytes:
            payload: Customer = handlers.codec.decode(input, Customer)
            result = await h.create_customer(payload)
            return handlers.codec.encode(result)

        handlers.register_handler('customers.v1.Inbound', 'createCustomer',
                                  handler)

    if not h.get_customer is None:

        async def handler(input: bytes) -> bytes:
            input_args: _InboundGetCustomerArgs = handlers.codec.decode(
                input, _InboundGetCustomerArgs)
            result = await h.get_customer(input_args.id)
            return handlers.codec.encode(result)

        handlers.register_handler('customers.v1.Inbound', 'getCustomer',
                                  handler)

    if not h.list_customers is None:

        async def handler(input: bytes) -> bytes:
            payload: CustomerQuery = handlers.codec.decode(
                input, CustomerQuery)
            result = await h.list_customers(payload)
            return handlers.codec.encode(result)

        handlers.register_handler('customers.v1.Inbound', 'listCustomers',
                                  handler)


@deserialize
@serialize
@dataclass
class _OutboundFetchCustomerArgs:
    id: int = field(default=int(0), metadata={'serde_rename': 'id'})


class OutboundImpl(Outbound):
    def __init__(self, invoker: Invoker):
        self.invoker = invoker

    async def save_customer(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'saveCustomer',
                                  customer)

    async def fetch_customer(self, id: int) -> Customer:
        input_args = _OutboundFetchCustomerArgs(id)
        return await self.invoker.invoke_with_return('customers.v1.Outbound',
                                                     'fetchCustomer',
                                                     input_args, Customer)

    async def customer_created(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'customerCreated',
                                  customer)


outbound = OutboundImpl(invoker)


def start():
    try:
        server.run(server_host, server_port)
    except KeyboardInterrupt:
        pass

    print("Server stopped.")
