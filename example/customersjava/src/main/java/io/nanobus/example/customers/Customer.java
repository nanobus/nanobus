package io.nanobus.example.customers;

import java.util.*;

public class Customer {
  private long id;
  private String firstName;
  private Optional<String> middleName;
  private String lastName;
  private String email;
  private Address address;

  public long getId() {
    return id;
  }

  public void setId(long id) {
    this.id = id;
  }

  public String getFirstName() {
    return firstName;
  }

  public void setFirstName(String firstName) {
    this.firstName = firstName;
  }

  public Optional<String> getMiddleName() {
    return middleName;
  }

  public void setMiddleName(Optional<String> middleName) {
    this.middleName = middleName;
  }

  public String getLastName() {
    return lastName;
  }

  public void setLastName(String lastName) {
    this.lastName = lastName;
  }

  public String getEmail() {
    return email;
  }

  public void setEmail(String email) {
    this.email = email;
  }

  public Address getAddress() {
    return address;
  }

  public void setAddress(Address address) {
    this.address = address;
  }
}
