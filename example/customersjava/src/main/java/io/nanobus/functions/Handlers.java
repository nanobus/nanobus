package io.nanobus.functions;

import java.util.function.Function;
import reactor.core.publisher.Mono;

public interface Handlers {
    void registerHandler(String namespace, String operation, Function<byte[], Mono<byte[]>> handler);
}
