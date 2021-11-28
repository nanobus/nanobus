![NanoBus Logo](docs/images/nanobus-logo.svg)

NanoBus is a lightweight microservice runtime layer that simplifies your application's core logic by moving infrastructure concerns to composable flows. The primary goal of NanoBus is to codify best practices so developers can **focus on business outcomes, not boilerplate code**.

## Key Features

### Virtually no boilerplate code

In conjunction with Dapr, NanoBus allows the developer to focus on what matters most, the application's logic. All the distributed system "glue" is handled automatically.

### Data-aware middleware / flows

Communicating with other services and Dapr building block is simplified in declarative, composable flows. Secure API endpoints, transform data, support multiple serialization formats, and apply resiliency policies using configuration.

### Automatic API endpoints with documentation

Share your service through multiple protocols, including REST, gRPC and NATS, without having to write additional code. OpenAPI/Swagger UI, AsyncAPI and Protobuf provide documentation for your partner teams.

### Consistent polyglot programming model

Using NanoBus and Dapr as a sidecar greatly simplifies distributed application development. Regardless of the chosen programming language, the developer experience feels like local development with plain interfaces and data structures.

### Clean Architecture

NanoBus applications are structured with design principles that allow your application to scale as requirements evolve. Newly created projects use an intuitive layout that follow best practices like [separation of concerns](https://en.wikipedia.org/wiki/Separation_of_concerns).

## How It Works

![NanoBus Architecture](docs/images/architecture.svg)

NanoBus runs jointly with [Dapr](https://dapr.io) in a [sidecar process](https://docs.microsoft.com/en-us/azure/architecture/patterns/sidecar). Conceptually, your application is plugged into the center and NanoBus handles bi-directional communication with Dapr and service/transport protocols.

Dapr provides developers with powerful [building blocks](https://docs.dapr.io/developing-applications/building-blocks/) such as service invocation, state management, publish and subscribe, secret stores, bindings, and actors. These building blocks are integrated with NanoBus flows. Flows are like middleware or data pipelines with configurable actions that perform operations like decoding, transform and routing data from the application to Dapr's components and visa-versa. No SDKs required.

To create services, NanoBus uses succinct yet flexible interface definitions to automatically produce API endpoints, like REST, [gRPC](https://grpc.io), and [NATS](https://nats.io). These transports are pluggable, allowing developers to expose services using multiple protocols without boilerplate code. Additionally, API documentation is auto-generated for customers.

Finally, NanoBus supports pluggable "compute" types: from Docker containers to emerging technologies
like [WebAssembly](https://webassembly.org). In the future, embedded language runtimes like JavaScript/TypeScript, Python or Lua could be supported.

To learn more, see the [architecture page](/docs/architecture.md).

## Getting Started

### Install the [nanogen CLI](https://github.com/nanobus/cli)

Windows

```
powershell -Command "iwr -useb https://raw.githubusercontent.com/nanobus/cli/master/install/install.ps1 | iex"
```

MacOS

```
curl -fsSL https://raw.githubusercontent.com/nanobus/cli/master/install/install.sh | /bin/bash
```

Linux

```
wget -q https://raw.githubusercontent.com/nanobus/cli/master/install/install.sh -O - | /bin/bash
```

Homebrew

```
brew install nanobus/tap/nanogen
```

### Create a NanoBus Application

Choose a supported language:

* Node.js (typescript)
* C# / .NET (csharp)
* Python (python)
* Golang (go)
* WASM AssemblyScript (assemblyscript)
* WASM TinyGo (tinygo)

Coming soon...

* Java (Reactor)
* Rust (Binary & WASM)

```shell
nanogen new typescript hello_world
cd hello_world
make
make run
```

In NanoBus, the developer only needs to follow these steps:

1. Create a new service using the `nanogen` CLI
2. Define the services interfaces (IDL)
3. Create flows that tie operations to Dapr building blocks
4. Implement the service's core logic code
5. Run `make docker`
6. Deploy to your favorite container orchestrator

Be sure to check out the [tutorial](example/README.md)!

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on the code of conduct and the process for submitting pull requests.

## License

This project is licensed under the [Apache License 2.0](https://choosealicense.com/licenses/apache-2.0/) - see the [LICENSE.txt](LICENSE.txt) file for details
