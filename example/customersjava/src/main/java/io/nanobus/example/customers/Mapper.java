package io.nanobus.example.customers;

import com.fasterxml.jackson.core.*;
import com.fasterxml.jackson.databind.*;

import java.io.*;

public class Mapper {
  private ObjectMapper objectMapper;

  public Mapper(ObjectMapper objectMapper) {
    this.objectMapper = objectMapper;
  }

  public <T> T decode(byte[] src, Class<T> valueType) {
    try {
      return objectMapper.readValue(src, valueType);
    } catch (IOException e) {
      throw new RuntimeException(e);
    }
  }

  public byte[] encode(Object value) {
    try {
      return objectMapper.writeValueAsBytes(value);
    } catch (JsonProcessingException e) {
      throw new RuntimeException(e);
    }
  }
}
