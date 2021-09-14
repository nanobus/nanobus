from adapter import register_inbound_handlers, start, outbound
from interfaces import Customer


async def create_customer(customer: Customer) -> Customer:
    await outbound.save_customer(customer)
    await outbound.customer_created(customer)
    return customer


async def get_customer(id: int) -> Customer:
    return await outbound.fetch_customer(id)


def main():
    register_inbound_handlers(
        create_customer=create_customer,
        get_customer=get_customer,
    )
    start()


if __name__ == "__main__":
    main()
