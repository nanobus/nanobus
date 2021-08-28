from nanobus import http_initialize
from interfaces import Customer, InboundHandlers, Outbound, OutboundImpl


class Inbound:
    def __init__(self, outbound: Outbound):
        self.outbound = outbound

    async def create_customer(self, customer: Customer) -> Customer:
        await self.outbound.save_customer(customer)
        await self.outbound.customer_created(customer)
        return customer

    async def get_customer(self, id: int) -> Customer:
        return await self.outbound.fetch_customer(id)


def main():
    (handlers, invoker, start) = http_initialize()

    outbound = OutboundImpl(invoker)
    inbound = Inbound(outbound)

    InboundHandlers(
        create_customer=inbound.create_customer,
        get_customer=inbound.get_customer,
    ).register(handlers)

    start()


if __name__ == "__main__":
    main()
