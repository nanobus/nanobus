# Python 3 server example
from interfaces import http_initialize
from functions import Customer, InboundHandlers, Outbound, OutboundImpl


class Inbound:
    def __init__(self, outbound: Outbound):
        self.outbound = outbound

    def create_customer(self, customer: Customer) -> Customer:
        self.outbound.save_customer(customer)
        self.outbound.customer_created(customer)
        return customer

    def get_customer(self, id: int) -> Customer:
        return self.outbound.fetch_customer(id)


if __name__ == "__main__":
    (handlers, invoker, start) = http_initialize()

    outbound = OutboundImpl(invoker)
    inbound = Inbound(outbound)

    InboundHandlers(
        create_customer=inbound.create_customer,
        get_customer=inbound.get_customer,
    ).register(handlers)

    start()
