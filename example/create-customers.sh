#!/bin/sh

wapc new @nanobus/codegen/go customers \
    description="Customers App" \
    version=1.0.0 \
    module="github.com/nanobus/nanobus/example/customers" \
    package="customers"

cd customers

cat > schema.widl <<EOF
namespace "customers.v1"
  @path("/v1")
  @info(
    title: "Customers API",
    description: """
    This API contains operations that an application can perform on customers.
    """,
    version: "1.0",
    contact: {
      name: "MyCompany API Team",
      email: "apiteam@mycompany.io",
      url: "http://mycompany.io"
    },
    license: {
      name: "Apache 2.0",
      url: "https://www.apache.org/licenses/LICENSE-2.0.html"
    }
  )
  @host("mycompany.io")
  @schemes(["https"])
  @consumes(["application/json"])
  @produces(["application/json"])
  @externalDocs(
    url: "http://mycompany.io/docs"
  )

import * from "@widl/restapi"
import * from "@widl/grpc"

"""
Operations that can be performed on a customer.
"""
role Inbound @service @path("/customers") {
  """
  Creates a new customer.
  """
  createCustomer{customer: Customer}: Customer
    @POST
    @response(
      status: CREATED,
      description: "Successful response",
      examples: {
        "application/json": "json"
      }
    )
  """
  Retrieve a customer by id.
  """
  getCustomer(id: u64 @fieldnum(1)): Customer
    @GET
    @path("/{id}")
    @response(
      status: OK,
      description: "Successful response",
      examples: {
        "application/json": "json"
      }
    )
    @response(
      status: NOT_FOUND,
      returns: "Error",
      description: "No customer with that identifier",
      examples: {
        "application/json": "json"
      }
    )
}

role Outbound {
  saveCustomer{customer: Customer}: void
  fetchCustomer(id: u64): Customer
  customerCreated{customer: Customer}: void
}

"""
Customer information.
"""
type Customer {
  "The customer identifer"
  id: u64 @key @fieldnum(1)
  "The customer's first name"
  firstName: string @fieldnum(2)
  "The customer's middle name"
  middleName: string? @fieldnum(3)
  "The customer's last name"
  lastName: string @fieldnum(4)
  "The customer's email address"
  email: string @email @fieldnum(5)
  "The customer's address"
  address: Address @fieldnum(6)
}

type Nested {
  foo: string @fieldnum(1)
  bar: string @fieldnum(2)
}

"""
Address information.
"""
type Address {
  "The address line 1"
  line1: string @fieldnum(1)
  "The address line 2"
  line2: string? @fieldnum(2)
  "The city"
  city: string @fieldnum(3)
  "The state"
  state: string @fieldnum(4) @length(min: 2, max: 2)
  "The zipcode"
  zip: string @fieldnum(5) @length(min: 5)
}

"""
Error response.
"""
type Error {
  "The detailed error message"
  message: string @fieldnum(1)
}
EOF

cat > bus.yaml <<EOF
specs:
  - type: widl
    with:
      filename: schema.widl

compute:
  # type: wapc
  # with:
  #   filename: customers.wasm
  type: mux
  with: {}

services:

  'customers.v1.Inbound':
    createCustomer____ignore:
      summary: Saves the customer to the database
      actions:
        - summary: Set the customer key
          name: '@dapr/set_state'
          with:
            store: statestore
            key: input.id

        - summary: Upsert the customer table
          name: '@dapr/sql_exec'
          with:
            name: postgres
            sql: |
              INSERT INTO customers (id, first_name, last_name)
              VALUES (:input.id, :input.firstName, :input.lastName)
              ON CONFLICT ON CONSTRAINT customers_pkey
              DO UPDATE SET first_name = :input.firstName, last_name = :input.lastName;

        - summary: Publish a message
          name: '@dapr/publish_message'
          with:
            pubsub: pubsub
            topic: test_topic
            format: cloudevents+json
            data: |
              {
                "type": "customer.created",
                "data": input
              }

        - summary: Return the input as the response
          name: assign
          with:
            value: input

outbound:

  'customers.v1.Outbound':
    saveCustomer:
      summary: Saves the customer to the database
      actions:
        - summary: Set the customer key
          name: '@dapr/set_state'
          with:
            store: statestore
            key: input.id

        - summary: Upsert the customer table
          name: '@dapr/sql_exec'
          with:
            name: postgres
            sql: |
              INSERT INTO customers (id, first_name, last_name)
              VALUES (:input.id, :input.firstName, :input.lastName)
              ON CONFLICT ON CONSTRAINT customers_pkey
              DO UPDATE SET first_name = :input.firstName, last_name = :input.lastName;

    fetchCustomer:
      summary: Loads the customer from the database
      actions:
        - summary: Get the state
          name: '@dapr/get_state'
          with:
            store: statestore
            key: input.id

    customerCreated:
      summary: Send a message to the customer
      actions:
        - summary: Publish a message
          name: '@dapr/publish_message'
          with:
            pubsub: pubsub
            topic: test_topic
            format: cloudevents+json
            data: |
              {
                "type": "customer.created",
                "data": input
              }
EOF

cat > dapr.sh <<EOF
# !/bin/sh
INVOKE_BASE_URL=http://localhost:8000 dapr run -d ../components --app-id customers --app-port 32321 --dapr-http-port 3500 -- nanobus --http-listen-addr :8081 --rest-listen-addr :8091 --bus-listen-addr localhost:32321 bus.yaml
EOF

chmod +x dapr.sh

cat > start.sh <<EOF
# !/bin/sh
go run cmd/main.go --http-listen-addr :8000 --outbound-base-uri http://localhost:32321/outbound
EOF

chmod +x start.sh