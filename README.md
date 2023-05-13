### Goals and Objectives

The primary goal of this project is to implement an API Gateway that accepts HTTP requests
encoded in JSON format and uses the Generic-Call feature of Kitex to translate these
requests into Thrift binary format requests. The API Gateway will then forward the request to

one of the backend RPC servers discovered from the registry center. To achieve this goal,
students will learn and apply the following concepts:

1. Golang: Gain proficiency in the Go programming language, which will be the primary
   language used for building the API Gateway, and become familiar with its syntax, data
   structures, and best practices.
2. HTTP: Understand the fundamentals of HTTP, including request and response
   structure, HTTP methods, and status codes.
3. JSON: Learn about JSON data interchange format and its application in encoding
   and decoding data.
4. Thrift: Acquire knowledge about Apache Thrift, an interface definition language (IDL)
   and binary communication protocol.
5. Load Balancing: Acquire knowledge of load balancing strategies for evenly
   distributing requests among available backend RPC servers, and integrate one of the Load
   Balancers provided by Kitex into the project to manage request distribution effectively.
6. Service Register and Discovery: Understand the principles of service registry and
   discovery mechanisms, and utilize one of the registry components offered by Kitex in this
   project, enabling Kitex-based services to register themselves and be discoverable by the API
   Gateway.
7. Building HTTP and RPC servers: Learn to build HTTP servers using Hertz and RPC
   servers using Kitex.

### Reference

• CloudWego
○ https://www.cloudwego.io/
• Kitex
○ https://www.cloudwego.io/docs/kitex
○ https://www.cloudwego.io/docs/kitex/tutorials/advanced-feature/generic-call/
○ https://github.com/cloudwego/kitex-examples
○ https://github.com/kitex-contrib/
▪ For load balance and registry components
• Hertz
○ https://www.cloudwego.io/docs/hertz/
○ https://github.com/cloudwego/hertz-examples
• Golang
○ https://go.dev/
○ https://go.dev/learn/
• HTTP
○ https://www.cloudflare.com/learning/ddos/glossary/hypertext-transfer-protocol-http/
• JSON
○ https://www.json.org/json-en.html
• Thrift
○ https://thrift.apache.org/
• Load Balance
○ https://www.nginx.com/resources/glossary/load-balancing/
• Service Registry and Discovery
○ https://www.nginx.com/blog/service-discovery-in-a-microservices-architecture/
