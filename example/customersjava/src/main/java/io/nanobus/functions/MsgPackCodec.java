package io.nanobus.functions;

import java.io.IOException;

import com.fasterxml.jackson.core.*;
import com.fasterxml.jackson.databind.*;
import org.msgpack.jackson.dataformat.MessagePackFactory;

public class MsgPackCodec implements Codec {
    private ObjectMapper objectMapper = new ObjectMapper(new MessagePackFactory())
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
            .configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false);

    public byte[] encode(Object value) {
        try {
            return objectMapper.writeValueAsBytes(value);
        } catch (JsonProcessingException e) {
            throw new RuntimeException(e);
        }
    }

    public <T> T decode(byte[] src, Class<T> clazz) {
        try {
            return objectMapper.readValue(src, clazz);
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}