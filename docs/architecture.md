## Architecture

NanoBus closely follows the architecture explained in [DDD, Hexagonal, Onion, Clean, CQRS, … How I put it all together](https://herbertograca.com/2017/11/16/explicit-architecture-01-ddd-hexagonal-onion-clean-cqrs-how-i-put-it-all-together/). Its proposed organization of software components is even more relevant now when considered with sidecar technology. The concepts below are paraphrased and sometimes renamed to better fit NanoBus' terminology.

Now we understand Hexagonal/Onion/Clean architecture. Let's break these concepts down as they relate to NanoBus.

### Ports and Adapters

**Ports** define the entry and exit points (or interfaces) of the application along with the data structures passed in and returned. **Adapters** are components that allow other systems to communicate with your application and vice-versa.

![NanoBus Architecture](images/architecture.svg)

**Primary/Driving Adapters** wrap around ports and instruct the application to perform operations. These can be thought of as transports such as REST/HTTP or gRPC. It can also be events consumed from Pub/Sub brokers such as [Apache Kafka](https://kafka.apache.org). In NanoBus applications, the Ports implemented by Primary Adapters are called **Services**.

Conversely, **Secondary/Driven Adapters** allow your application to interact with components that persist or retrieve data. This could be from a relational database like PostgreSQL or MySQL. A document database like MongoDB, or a key-value store like Redis. In NanoBus applications, the Ports implemented by Secondary Adapters are called a **Providers**.

Both Primary and Secondary Adapters are constructed using composable flows that interact with the application or a component like a Dapr building block.

### Flow-Based Programming

From the article mentioned above.

> You might have noticed that there is no dependency between the Bus and the Command, the Query nor the Handlers. This is because they should, in fact, be unaware of each other in order to provide for good decoupling. The way the Bus will know what Handler should handle what Command, or Query, should be set up with mere configuration.

Between your application and adapters, NanoBus passes data through developer-defined flows: a paradigm called **Flow-Based Programming (FBP)**.

![Flow-Based Programming example](images/fbp.svg)

Each step in the flow can inspect, transform, and augment the data before passing it to the destination component.  Additionally, this mechanism externalizes encoding/decoding, authorization, cryptography and resiliency policies, eliminating the need to update the application internally. The application is only concerned with operations that accept and return strongly typed data structures.

### Onion Architecture

<img align="right" src="images/onion.svg" alt="The Onion Architecture" width="45%" style="margin: 0 0 1rem 1rem;" />

NanoBus applications follow the **Onion Architecture**, comprised of concentric layers that interface with each other towards the center. This design greatly improves a system's ability to evolve over time because each layer addresses a separate concern. Here are the layers NanoBus implements:

* **Core Logic** contains the application's business logic, rules, calculations, algorithms, etc.
* **Provider Services** wrap around ports that invoke flows that act as driving adapters for Dapr or other components.
* **Data Model** (or domain modal) contains data structures for persistent data and events that are communicated through port / flows.

Another goal of this software composition is abstracting away frameworks. Since these layers are plain interfaces and data structures, the core logic can be easily shifted to a different runtime technology or custom integration code.

<br clear="both"/>

### API-First Approach

To tie layers and flows together, the developer can *optionally* use an **API-First Approach**, where the operations of your service are described in an API specification that can be shared with other teams.

![IDL code generation](images/idl-codegen.svg)

NanoBus uses a generic and protocol agnostic Interface Definition Language (IDL). This is input to our code generation tool, which creates RPC-style endpoints that your logic uses to interact with flows. NanoBus also uses the IDL to automatically host your services with REST, gRPC, and other protocols. This eliminates the need to write boilerplate code for the Transports, Endpoints, and Stores layers.
