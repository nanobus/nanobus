package io.nanobus.functions;

import reactor.core.publisher.Mono;

public class Invoker {
    private Handler invokeFn;

    private Codec codec;

    public Invoker(Handler invoke, Codec codec) {
        this.invokeFn = invoke;
        this.codec = codec;
    }

    public Mono<Void> invoke(String namespace, String operation, Object value) {
        byte[] data = codec.encode(value);
        return invokeFn.apply(namespace, operation, data).then();
    }

    public <T> Mono<T> invokeWithReturn(String namespace, String operation, Object value, Class<T> clazz) {
        byte[] data = codec.encode(value);
        return invokeFn.apply(namespace, operation, data).map(bytes -> codec.decode(bytes, clazz));
    }
}
