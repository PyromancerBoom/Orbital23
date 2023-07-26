This is the repository for the ByteDance and Tiktok Orbital 2023.

# Table of Contents

1. [About](#about)
2. [Features and Design](#features)
   - [Components](#components)
   - [IDL Management](#idlmanagement)
   - [Server Utility Package](#serverutils)
3. [Getting Started with an example](#gettingstarted)
   - [Setting up kitex](#kitexsetup)
4. [Performance](#performance)
5. [Data](#data)
6. [Fourth Example](#fourth-examplehttpwwwfourthexamplecom)

## About <a name="about"></a>

This is the project for our API Gateway based on one scalable Hertz server which serves multiple RPC servers.

The API Gateway has the following endpoints :

- _GET and POST /:serviceName/:path_
  Every service has a registered URL which is it's s

- _GET /ping_
  This endpoint checks the availability of the API gateway (ping).
  Run it to verify if the server is running or not

- _POST /register_
  Handles user registration requests.

- _PUT /update_
  Server to Update details of Registered services

- _POST /connect_
  Once services are registered, their servers are connected through this endpoint automatically.

- _POST /health_
  This endpoint handles health checks for a service. Services ping to /health to update their healthchecks.

<a href="#top">Back to top</a>

### Features and Design <a name="features"></a>

### The project has the following features : <a name="components"></a>

1. API Gateway Server: The API Gateway is implemented as a Hertz server that listens to requests on port 4200. It exposes multiple endpoints in the format `/{serviceName}/{path}` for both POST and GET requests. The API Gateway acts as an intermediary between user requests and the Kitex RPC servers performing load balancing, service discovery, health checks, etc. It accepts incoming HTTP requests and decodes on which service to make a generic RPC call and returns a response.

2. Service Registration and Caching: A user can send in POST requests on `:/register` to register their services. The JSON payload includes details about the owner, their services, IDLs and masked service URLs. This data is stored in MongoDB and is also cached inside the gateway server for quick access. Refer to the <a href="#data">Data</a>section below for the exact info on how data is stored.

3. RPC Protocal Translation: The API Gateway forwards incoming API requests to the Kitex servers using the internal RPC clients within the Hertz server. It enables communication between the client and the respective RPC servers responsible for handling specific services. This is done with their Thrift IDL information sent on registration.

4. Automated Server Connection: The API Gateway provides another package called _server_utils_ to make it easy for Services can automate the registration of their servers to our system by making an HTTP POST request at the `:/connect` endpoint of the API Gateway using their registered API key. This enables services to scale up or down dynamically according to their needs.

5. Health Checks: Servers connected to the API Gateway need to declare their health by making periodic requests to the `:/health` endpoint at least every 10 seconds. This ensures that the API Gateway considers the servers healthy and forwards requests to them. This is part of the Server Utility Package and can be easily automated. For this, we have used Consul.

6. Discovery and Load Balancing: The gateway uses Consul's resolver for service discovery and load balancing. It uses Consul’s DNS interface to resolve service names. It can be used to discover services registered with Consul and load balance requests between them. The MVP version of the NewConsulResolver implements round-robin load balancing, distributing the requests equally among the connected RPC servers.

7. Service Registry: Consul has been integrated as the service registry as well. It is completely isolated from the RPC servers. The servers have to therefore, interact with the Gateway for Consul related business. This was done to prevent malicious attacks on the service registry. The consul service registry provides a beautiful graphical view of all the connected servers to help admins track, manage and troubleshoot connections. This also allows us for more freedom in logic handling for tasks related to the Registry.

- Consul Agent is for the purpose of Service Registry, Health Checks, Discovery, Load Balancing. It is hosted on `localhost:8500` for this project. The logic for exact functioning however has been coded manually.

### IDL Management <a name="idlmanagement"></a>

IDL (Interface Definition Language) management plays a crucial role in facilitating communication between the API Gateway and the backend RPC servers using Kitex. Here's how IDLs are managed in the project:

During service registration with POST /register, service owners include their service's IDL details in the JSON payload.

- IDL Mapping and Translation:
  API Gateway stores and uses the IDL contents to translate incoming JSON API requests to Thrift binary format for Kitex.
  The IDLs are also cached for faster API Calls and are used to make and store all the generic clients so that the parsing need to be done again and again. Thrift's compact binary format facilitates efficient communication between the API Gateway and backend RPC servers.

- IDL Versioning and Compatibility:
  The Data Model allows services to keep track of their IDLs and versions. Although one limitation (due to time constraints) is that only the latest IDL version is stored. IDLs can be updated using the `/update` endpoint.

- Synchronization:
  Ensuring IDL information remains consistent and synchronized between components is crucial, which is why IDLs are cached regularly.

<a href="#top">Back to top</a>

#### Registration request format

```
[
    {
      "OwnerName": "John Doe 2",
      "OwnerId": "hellowworld",
      "Services": [
        {
          "ServiceId": "1",
          "ServiceName": "AssetManagement",
          "IdlContent": "namespace Go asset.management\n\nstruct QueryAssetRequest {\n    1: string ID;\n}\n\nstruct QueryAssetResponse {\n    1: bool   Exist;\n    2: string ID;\n    3: string Name;\n    4: string Market;\n}\n\nstruct InsertAssetRequest {\n    1: string ID;\n    2: string Name;\n    3: string Market;\n}\n\nstruct InsertAssetResponse {\n    1: bool Ok;\n    2: string Msg;\n}\n\nservice AssetManagement {\n    QueryAssetResponse queryAsset(1: QueryAssetRequest req);\n    InsertAssetResponse insertAsset(1: InsertAssetRequest req);\n}\n",
          "Version": "1.0",
          "ServiceDescription": "Service A Description",
          "ServerCount": 2,
          "Paths": [
            {
              "ExposedMethod": "insertAsset",
              "MethodPath": "newAsset"
            },
            {
              "ExposedMethod": "queryAsset",
              "MethodPath": "getAsset"
            }
          ],
          "RegisteredServers": [
            {
              "ServerUrl": "http://localhost:8000",
              "Port": 8000
            },
            {
              "ServerUrl": "http://localhost:8001",
              "Port": 8001
            }
          ]
        }
      ]
    }
]
```

## Performance

On Load testing with Postman, we were able to have the following benchmarks:
The lower the blue line is, the better.
The red line indicates error rate.

![performance1](perf-25users-mvp.png)
2 instances of User Service and 3 instances of Asset Management Service

- Users : 25
- Total time : 5 mins
- Ramp up time: 1 min

![performance1](perf-50users-mvp.png)
2 instances of User Service and 3 instances of Asset Management Service

- Users : 50
- Total time : 5 mins
- Ramp up time: 1 min

Despite the spike in between the server showed great recovery.

![performance1](perf-30users.png)
3 instances of User Service and 3 instances of Asset Management Service

- Users : 30
- Total time : 5 mins
- Ramp up time: 1 min

Again, after the spike, the gateway showed great recovery.

### Registration :

For now, the registration functionality is NOT integrated with the service registry. However, service information for registration can be sent at the (POST) `:/register` endpoint with the following json format which is accepted with our service registry Consul as well :

```
[
  {
    "Service": {
      "Name": "first-service",
      "Tags": [],
      "Address": "serviceName/path",
      "Port": 80,
      "Meta": {
        "serviceDescription": "Service for managing assets",
        "serviceVersion": "1.0",
        "idl":"idlcontent" (Required)
      },
      "Check": {
        "HTTP": "http://serviceAddress/service/myservice/1/health",
        "Interval": "10s"
      }
    }
  }
]
```

The `"idl":"idlcontent"` field is required as require the IDL content to make the RPC calls from the gateway to the service.

Later, we will connect the registration with our registry and IDL mappings.

Moreover, to see all registered service we can hit the endpoint `:/show`. (This will be removed later and was made just for testing purposes)

## Connecting services to the gateway:

Ensure that your servers can perform all the methods indicated in your interface definition during service registry.

- Register your servers via the `:/connect` endpoint each time a new server is booted up. Use the `server_utils.go` file for the methods which can be used in the service to connect to the gateway. Examples are present in the `main.go` of the services.
- Declare your servers as healthy to our system by making requests to the `:/health` endpoint at least every 10 seconds.

To connect your server to our system, follow these steps:

1. Register your service and receive an API Key.
2. On server bootup, send a request to `:/connect` with the API Key, service details, server address, and port.
3. Get a serverID upon successful connection.
4. Declare your server's health by regularly sending requests to `:/health` with the API Key and serverID.
   Remember:

Note:

- New servers should register themselves via :/connect.
- Servers must declare themselves healthy every 10 seconds.
- Servers are delisted if they don't declare health for 1 minute.
- If delisted, reconnect..

Load balancing currently uses round-robin, but will be upgraded to weighted round-robin.

For the detailed guide on service connection, check out [Server Connection Guide](ServerConnectionGuide.md)

## How to use? [^3]

**Step 1:**

Initialise the Hertz Server using the command: `go run .` from the respective directory

To check if the server is running, hit the following GET endpoint
`"http://localhost:4200/ping"`

It should reply with the message :

```
{
    "message": "pong"
}
```

**Step 2:**

Start consul agent using `consul agent -dev`. (Consul needs to be installed for this)
The consul GUI can be accessed at `http://localhost:8500`

**Step 3:**

Initialise multiple Kitex Services from their respective directories by using the command `go run .`
There can be multiple instances for a service. The service would automatically detect a free port and start the service on that free port locally.

Currently, we have the following functional services :

1. Asset Management

##### Public endpoints:

- (POST) `/AssetManagement/newAsset` which maps to the private "insertAsset" endpoint of the service
- (GET) `/AssetManagement/getAsset` which maps to the private "queryAsset" endpoint of the service

2. User Service

##### Public endpoints:

- (POST) `/UserService/newUser` which maps to the "insertUser" endpoint of the service
- (POST) ` / UserService/insertUser` which also maps to the same "insertUser" endpoint of the service
- (GET) `/AssetManagement/getUser` which maps to the private "queryUser" endpoint of the service

The expected data for the above endpoints is provided below in Step 4.

Once initialised they are automatically connected to consul, for example :

![consul](consulservicesconnection.png)

**Step 4:**

Send a POST or GET requests to the "/{serviceName}/{path}" endpoint, for example:

#### Asset Management

```
curl -X POST -H "Content-Type: application/json"
-d '{
  "ID": "2",
  "Name": "Google",
  "Market": "US"
}'
"http://localhost:4200/AssetManagement/newAsset"
```

Now try quering the info,

```
curl -X GET http://localhost:4200/AssetManagement/getAsset?ID=2
```

#### User service

```
curl -X POST -H "Content-Type: application/json"
-d '{
  "id": "3",
  "name": "Doe",
  "email": "johndoe@example.com",
  "age": 24
}'
"http://localhost:4200/UserService/newUser"
```

Now try quering the info,

```
curl -X GET http://localhost:4200/UserService/getUser?id=2
```
