from nanobus import start
from spec import Customer, registerInboundHandlers, outbound


async def create_customer(customer: Customer) -> Customer:
    await outbound.save_customer(customer)
    await outbound.customer_created(customer)
    return customer


async def get_customer(id: int) -> Customer:
    return await outbound.fetch_customer(id)


def main():
    registerInboundHandlers(
        create_customer=create_customer,
        get_customer=get_customer,
    )
    start()


if __name__ == "__main__":
    main()
