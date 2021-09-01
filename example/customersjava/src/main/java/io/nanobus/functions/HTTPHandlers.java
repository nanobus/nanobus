package io.nanobus.functions;

import java.util.HashMap;
import java.util.Map;
import java.util.function.Function;

import io.netty.handler.codec.http.HttpMethod;
import io.netty.handler.codec.http.HttpResponseStatus;
import io.netty.handler.ssl.util.SelfSignedCertificate;
import reactor.core.publisher.Mono;
import reactor.netty.http.Http11SslContextSpec;
import reactor.netty.http.HttpProtocol;
import reactor.netty.http.server.HttpServer;

public class HTTPHandlers implements Handlers {
    static final boolean SECURE = System.getProperty("secure") != null;
    static final boolean WIRETAP = System.getProperty("wiretap") != null;
    static final boolean COMPRESS = System.getProperty("compress") != null;
    static final boolean HTTP2 = System.getProperty("http2") != null;

    private Map<String, Function<byte[], Mono<byte[]>>> handlers = new HashMap<>();

    @Override
    public void registerHandler(String namespace, String operation, Function<byte[], Mono<byte[]>> handler) {
        if (handler != null) {
            handlers.put(namespace + "/" + operation, handler);
        }
    }

    public void listen(int port, String host) throws Exception {
        HttpServer server = HttpServer //
                .create() //
                .port(port) //
                .host(host) //
                .wiretap(WIRETAP) //
                .compress(COMPRESS) //
                .handle((req, res) -> {
                    if (req.method() != HttpMethod.POST) {
                        res.status(HttpResponseStatus.METHOD_NOT_ALLOWED)
                            .sendString(Mono.just("Invalid method")).then();
                    }
                    Function<byte[], Mono<byte[]>> handler = handlers.get(req.path());
                    if (handler == null) {
                        return res.status(404) //
                                .sendString(Mono.just("Not Found")).then();
                    }

                    return res.sendByteArray(//
                            req.receive().retain().aggregate().asByteArray() //
                                    .flatMap(handler::apply));
                });

        if (SECURE) {
            SelfSignedCertificate ssc = new SelfSignedCertificate();
            server = server.secure(spec -> spec.sslContext(Http11SslContextSpec. //
                    forServer(ssc.certificate(), ssc.privateKey())));
        }

        if (HTTP2) {
            server = server.protocol(HttpProtocol.H2);
        }

        server.bindNow().onDispose().block();
    }
}
