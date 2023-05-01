TopSpin is a library for building microservices and micromonoliths in Go based on Domain Driven Design (DDD) practices.

The true advantage of this approach may not be immediately apparent in small service projects or those that do not require long-term maintenance or extension, but it notably simplifies the addition of new business features on a consistent basis. Furthermore, since the implementation is defined using real business terms (ubiquitous language), it is not bound to a REST model, making code navigation and discovery more natural.

Some of the choices made in this library may seem overly ceremonial within the Go ecosystem. However, the intention is to use it as a more orthodox reference starting point and later to move development to a more simplified and ergonomic scheme. Additionally, we plan to construct a code generator that will make the development of services based on this library even simpler.

For use cases in which the entire service or parts of it do not benefit from the use of DDD patterns, we provide the possibility of operating directly in a RESTful scheme (e.g., during sign-in and sign-out).

Notice that a REST adapter still make use of the commands and queries. However, now GET and POST requests will be sufficient, and the only meaningful information in the URL will be the name of the command to be executed or a suitable alias-command mapping. Moreover, the API documentation (OpenAPI) becomes more humanly understandable [TODO: show an example], especially for people with less technical background. The latter in the context of API manual testing can be very convenient.

Some of the advantages of using this library include simplified testing since we are isolated from infrastructure issues (e.g., HTTP, gRPC, etc.), and it is possible to manually send queries and commands from the console if required. This is trivial to achieve since it is easy to create a CLI adapter that operates on the elements of the bounded context. Such a feature can be particularly useful during the development stage.
