# NanoBus Demo

## How I want to build applications

- Simple function-based programming model
- Data structures are strongly typed without serialization concerns
- Limited boilerplate / lightweight (or no) SDK
- Abstraction layers: Business logic and Integration
- Simplified integration with other services and data stores
- No worrying about data formats, OpenAPI or gRPC

Opinionated yet flexible!


## The Demo (two ideas that can work in tandem)

1) Function-based programming model
    - Dapr as the infrastructure and building blocks
    - IDL and code generation to create "the glue"

2) Flow-based computing model for Dapr
    - Address cross-cutting concerns
    - Annotation or Configuration driven
    - Acts as the integration layer for common use cases


## Goal: Codify best practices into a sidecar so you can focus on business logic

- Validation
- Versioning
- Resiliency
  - Timeouts
  - Retry / Backoff
  - Circuit breaking / Fallback
- Portability / Multi-cloud
- Mock testing
- A/B testing
- Live reloading
- Wire tap into the traffic
- Data transformation / Filtering
- Content-based routing
- Generated API Docs
---------------------------------------
- OAuth Client Credentials
- Custom TLS Configs
- Schedules / Triggers / Jobs (Locks) K8s?
- Global locks
- Leadership election
- Eventual consistency / Saga pattern
- Database Migrations (DDL) ???
---------------------------------------
- Tracing / Observability
- Service discovery

Related issues:
* https://github.com/dapr/dapr/issues/1820
* https://github.com/dapr/dapr/issues/3000 & https://github.com/dapr/dapr/pull/3074/files
* https://github.com/dapr/dapr/issues/2582
* https://github.com/dapr/dapr/issues/724
* https://github.com/dapr/dapr/issues/1556
* https://github.com/dapr/dapr/issues/777
* https://github.com/dapr/dapr/issues/501 & https://github.com/dapr/dapr/issues/1601
* https://github.com/dapr/dapr/issues/2716 (context is passed through)

## Project initialization

Run PostgreSQL

```shell
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -d postgres

CREATE TABLE customers (
	id int primary key,
	first_name varchar(200) NOT NULL,
	middle_name varchar(200),
	last_name varchar(200) NOT NULL
);
```
