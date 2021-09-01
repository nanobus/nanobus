package io.nanobus.functions;

public interface Codec {
    byte[] encode(Object value);
    <T> T decode(byte[] data, Class<T> clazz);
}