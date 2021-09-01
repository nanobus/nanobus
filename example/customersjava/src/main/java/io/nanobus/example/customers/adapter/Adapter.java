package io.nanobus.example.customers.adapter;

import io.nanobus.functions.Codec;
import io.nanobus.functions.HTTPHandlers;
import io.nanobus.functions.HTTPInvoker;
import io.nanobus.functions.Invoker;
import io.nanobus.functions.MsgPackCodec;

import io.nanobus.example.customers.models.*;
import io.nanobus.example.customers.service.*;

public class Adapter {
    static final int PORT = Integer.parseInt(System.getProperty("port", "9000"));
    static final String HOST = System.getProperty("host", "localhost");
    static final String OUTBOUND_BASE_URL = System.getProperty("outbound.base.url", "http://localhost:32321/outbound");

    private HTTPHandlers handlers;
    private Codec codec;
    private Invoker invoker;

    public Adapter() {
        this.codec = new MsgPackCodec();
        HTTPInvoker invoke = new HTTPInvoker(OUTBOUND_BASE_URL);
        this.invoker = new Invoker((namespace, operation, input) -> //
        invoke.invoke(namespace, operation, input), this.codec);
        this.handlers = new HTTPHandlers();
    }

    public void start() throws Exception {
        this.handlers.listen(PORT, HOST);
    }

    public void registerInboundHandlers(Inbound inbound) {
        handlers.registerHandler("customers.v1.Inbound", "createCustomer", input -> {
            Customer customer = codec.decode(input, Customer.class);
            return inbound.createCustomer(customer)
                .map(codec::encode);
        });
        handlers.registerHandler("customers.v1.Inbound", "getCustomer", input -> {
            GetCustomerArgs args = codec.decode(input, GetCustomerArgs.class);
            return inbound.getCustomer(args.getId())
                .map(codec::encode);
        });
    }

    public Outbound newOutbound() {
        return new OutboundImpl(invoker);
    }
}
