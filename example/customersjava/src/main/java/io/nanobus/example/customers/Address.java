package io.nanobus.example.customers;

import java.util.*;

public class Address {
  private String line1;
  private Optional<String> line2;
  private String city;
  private String state;
  private String zip;

  public String getLine1() {
    return line1;
  }

  public void setLine1(String line1) {
    this.line1 = line1;
  }

  public Optional<String> getLine2() {
    return line2;
  }

  public void setLine2(Optional<String> line2) {
    this.line2 = line2;
  }

  public String getCity() {
    return city;
  }

  public void setCity(String city) {
    this.city = city;
  }

  public String getState() {
    return state;
  }

  public void setState(String state) {
    this.state = state;
  }

  public String getZip() {
    return zip;
  }

  public void setZip(String zip) {
    this.zip = zip;
  }
}
