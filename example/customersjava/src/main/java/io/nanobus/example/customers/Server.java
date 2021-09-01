package io.nanobus.example.customers;

import io.nanobus.example.customers.adapter.Adapter;
import io.nanobus.example.customers.service.*;

public final class Server {
    public static void main(String[] args) throws Exception {
        Adapter adapter = new Adapter();

        Outbound outbound = adapter.newOutbound();
        Inbound inbound = new InboundImpl(outbound);
        adapter.registerInboundHandlers(inbound);

        adapter.start();
    }
}