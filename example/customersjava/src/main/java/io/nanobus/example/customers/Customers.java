package io.nanobus.example.customers;

import com.fasterxml.jackson.core.*;
import com.fasterxml.jackson.databind.*;
import io.netty.handler.ssl.util.SelfSignedCertificate;
import reactor.core.publisher.Mono;
import reactor.netty.http.Http11SslContextSpec;
import reactor.netty.http.HttpProtocol;
import reactor.netty.http.client.HttpClient;
import reactor.netty.http.server.HttpServer;
import org.msgpack.jackson.dataformat.MessagePackFactory;

import java.io.*;

import static io.netty.handler.codec.http.HttpHeaderNames.CONTENT_TYPE;

public final class Customers {

  static final boolean SECURE = System.getProperty("secure") != null;
  static final int PORT = Integer.parseInt(System.getProperty("port", SECURE ? "8443" : "8080"));
  static final boolean WIRETAP = System.getProperty("wiretap") != null;
  static final boolean COMPRESS = System.getProperty("compress") != null;
  static final boolean HTTP2 = System.getProperty("http2") != null;

  public static void main(String[] args) throws Exception {
    // Instantiate ObjectMapper for MessagePack
    ObjectMapper objectMapper = new ObjectMapper(new MessagePackFactory())
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);

    Mapper mapper = new Mapper(objectMapper);

    HttpClient client =
            HttpClient.create()
                    .host("localhost")
                    .baseUrl("/outbound")
                    .port(32321);

    Outbound outbound = new OutboundImpl(client, mapper);
    Inbound inbound = new InboundImpl(outbound);


    HttpServer server =
            HttpServer.create()
                    .port(PORT)
                    .wiretap(WIRETAP)
                    .compress(COMPRESS)
                    .route(r -> r
                            .post("/customers.v1.Inbound/createCustomer",
                                    (req, res) ->
                                            res.sendByteArray(
                                                    req.receive().retain().aggregate().asByteArray()
                                                            .map(inputBytes -> mapper.decode(inputBytes, Customer.class))
                                                            .flatMap(customer -> inbound.createCustomer(customer).
                                                                    map(mapper::encode)))
                            )
                            .post("/customers.v1.Inbound/getCustomer",
                                    (req, res) ->
                                            res.sendByteArray(
                                                    req.receive().retain().aggregate().asByteArray()
                                                            .map(inputBytes -> mapper.decode(inputBytes, GetCustomerArgs.class))
                                                            .flatMap(a -> inbound.getCustomer(a.getId()).
                                                                    map(mapper::encode)))
                            )
                    );

    if (SECURE) {
      SelfSignedCertificate ssc = new SelfSignedCertificate();
      server = server.secure(
              spec -> spec.sslContext(Http11SslContextSpec.forServer(ssc.certificate(), ssc.privateKey())));
    }

    if (HTTP2) {
      server = server.protocol(HttpProtocol.H2);
    }

    server.bindNow()
            .onDispose()
            .block();
  }
}