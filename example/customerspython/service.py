#!/usr/bin/env python3
from adapter import start, register_inbound_handlers, outbound
from interfaces import Inbound, Customer, CustomerPage, CustomerQuery


async def create_customer(customer: Customer) -> Customer:
    await outbound.save_customer(customer)
    await outbound.customer_created(customer)
    return customer


async def get_customer(id: int) -> Customer:
    return await outbound.fetch_customer(id)


async def list_customers(query: CustomerQuery) -> CustomerPage:
    return CustomerPage(
        offset=query.offset,
        limit=query.limit,
    )


def main():
    register_inbound_handlers(
        Inbound(
            create_customer=create_customer,
            get_customer=get_customer,
            list_customers=list_customers,
        ))

    start()


if __name__ == "__main__":
    main()
