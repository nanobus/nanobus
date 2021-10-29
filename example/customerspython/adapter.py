import os
from nanobus import UvicornServer, HTTPInvoker, Handlers, Invoker, MsgPackCodec, JsonCodec
from stateful import Manager, Storage, LRUCache
from serde import serialize, deserialize
from dataclasses import dataclass, field
from interfaces import Inbound, Customer, CustomerPage, CustomerQuery, CustomerActor, Outbound

server_host = os.getenv('HOST', "127.0.0.1")
server_port = int(os.getenv('PORT', "9000"))
bus_url = os.getenv('BUS_URL', "http://127.0.0.1:32321")

codec = MsgPackCodec()
handlers = Handlers(codec)
server = UvicornServer(handlers)
http_invoker = HTTPInvoker(bus_url + "/providers")
invoker = Invoker(http_invoker.invoke, codec)
json_codec = JsonCodec()
cache = LRUCache()
storage = Storage(bus_url, json_codec)
state_manager = Manager(cache, storage, json_codec)


@deserialize
@serialize
@dataclass
class _InboundGetCustomerArgs:
    id: int = field(default=int(0), metadata={'serde_rename': 'id'})


def register_inbound(h: Inbound):
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


def register_customer_actor(h: CustomerActor):
    if not h.create_customer is None:

        async def handler(id: str, input: bytes) -> bytes:
            sctx = await state_manager.to_context("customers.v1.CustomerActor",
                                                  id, h)
            payload: Customer = handlers.codec.decode(input, Customer)
            result = await h.create_customer(sctx, payload)
            response = sctx.response(result)
            return handlers.codec.encode(response)

        handlers.register_stateful_handler('customers.v1.CustomerActor',
                                           'createCustomer', handler)

    if not h.get_customer is None:

        async def handler(id: str, input: bytes) -> bytes:
            sctx = await state_manager.to_context("customers.v1.CustomerActor",
                                                  id, h)
            result = await h.get_customer(sctx)
            response = sctx.response(result)
            return handlers.codec.encode(response)

        handlers.register_stateful_handler('customers.v1.CustomerActor',
                                           'getCustomer', handler)


@deserialize
@serialize
@dataclass
class _OutboundFetchCustomerArgs:
    id: int = field(default=int(0), metadata={'serde_rename': 'id'})


class OutboundImpl(Outbound):
    def __init__(self, invoker: Invoker):
        self.invoker = invoker

    # Saves a customer to the backend database
    async def save_customer(self, customer: Customer):
        await self.invoker.invoke('customers.v1.Outbound', 'saveCustomer',
                                  customer)

    # Fetches a customer from the backend database
    async def fetch_customer(self, id: int) -> Customer:
        input_args = _OutboundFetchCustomerArgs(id)
        return await self.invoker.invoke_with_return('customers.v1.Outbound',
                                                     'fetchCustomer',
                                                     input_args, Customer)

    # Sends a customer creation event
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
