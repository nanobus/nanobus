package io.nanobus.functions;

import java.net.MalformedURLException;
import java.net.URL;

import reactor.core.publisher.Mono;
import reactor.netty.http.client.HttpClient;

public class HTTPInvoker {
    private HttpClient client;

    public HTTPInvoker(String baseURL) {
        try {
            URL url = new URL(baseURL);
            int port = url.getPort();
            this.client = HttpClient.create() //
                    .host(url.getHost()) //
                    .port(port != -1 ? port : url.getDefaultPort()) //
                    .baseUrl(url.getPath());
        } catch (MalformedURLException e) {
            throw new RuntimeException(e);
        }
    }

    public Mono<byte[]> invoke(String namespace, String operation, byte[] input) {
        String uri = "/" + namespace + "/" + operation;
        return this.client.post().uri(uri).send((request, outbound) -> outbound.sendByteArray(Mono.just(input)))
                .responseContent().retain().aggregate().asByteArray();
    }
}
