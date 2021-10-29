#!/usr/bin/env python3
from adapter import start, register_inbound, register_customer_actor, outbound
from interfaces import Context, Inbound, Customer, CustomerPage, CustomerQuery, CustomerActor


class InboundImpl(Inbound):
    async def create_customer(self, customer: Customer) -> Customer:
        await outbound.save_customer(customer)
        await outbound.customer_created(customer)
        return customer

    async def get_customer(self, id: int) -> Customer:
        return await outbound.fetch_customer(id)

    async def list_customers(self, query: CustomerQuery) -> CustomerPage:
        return CustomerPage(
            offset=query.offset,
            limit=query.limit,
        )


class CustomerActorImpl(CustomerActor):
    async def activate(self, ctx: Context):
        print("ACTIVATED")

    async def create_customer(self, ctx: Context,
                              customer: Customer) -> Customer:
        ctx.set("customer", customer)
        return customer

    async def get_customer(self, ctx: Context) -> Customer:
        return await ctx.get("customer", Customer)

    async def deactivate(self, ctx: Context):
        print("DEACTIVATED")


def main():
    register_inbound(InboundImpl())
    register_customer_actor(CustomerActorImpl())

    start()


if __name__ == "__main__":
    main()
