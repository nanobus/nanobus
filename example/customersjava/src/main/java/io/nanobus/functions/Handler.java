package io.nanobus.functions;

import reactor.core.publisher.Mono;

@FunctionalInterface
public interface Handler {
    Mono<byte[]> apply(String namespace, String operation, byte[] input);
}