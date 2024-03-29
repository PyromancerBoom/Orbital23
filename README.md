This is the repository for the ByteDance and Tiktok Project for 2023 NUS Orbital.

This README markdown file aims to provide a comprehensive documentation for this project.

# Table of Contents

1. [About](#about)
2. [Features and Design](#features)
   - [Architecture Diagram](#diagram)
   - [Components](#components)
   - [IDL Management](#idlmanagement)
   - [Server Utility Package](#serverutils)
3. [Data Management](#data)
4. [Getting Started with an example](#gettingstarted)
   1. [Set up Hertz](#step1)
      - [Settings for Gateway](#hertzsettings)
   2. [Set up Consul](#step2)
   3. [Start MongoDB](#step3)
   4. [Setting up kitex](#step4)
      - [Settings for Kitex](#kitexsettings)
   5. [Registering a service](#step5)
   6. [Updating Data](#step6)
   7. [Send requests](#step7)
   8. [Proxy to Ramp Up Performance](#step8)
5. [Performance](#perf)
6. [Testing](#tests)
7. [Limitations/Issues](#limit)
8. [What Else?](#misc)

## About <a name="about"></a>

This is the project for our API Gateway based on one scalable Hertz server which serves multiple RPC servers.

The API Gateway has the following endpoints :

- _GET and POST /:serviceName/:path_ - Every service has a registered URL on which its users send requests.

- _GET /ping_ - This endpoint checks the availability of the API gateway. Run it to verify if the server is running or not

- _POST /register_ - Handles admin registration requests

- _PUT /update_ - Server to Update details of Registered services

- _POST /connect_ - Once services are registered, their servers are connected through this endpoint automatically with the help of the server utility package.

- _POST /health_ - This endpoint handles health checks for a service. Services ping to /health to update their health checks.

#### About the directory:

- The api-gateway folder contains the api gateway server and a server_utils package which is required for a smooth setup of Kitex servers.

- The rpc_services contains the services used throughout the development of the project. It has 3 services - Asset Service, User service, and Registry Proxy service. Details have been provided throughout this documentation. But TLDR:

  - Asset Service is used to store some info about Assets and fetch that info.
  - The User service is similar to the asset service used for registering and returning info about people using that service.
  - As for Registry Proxy service, it is a special service which performs server health checks and connection requests made to the gateway. By simply booting up the Registry proxy service, the gateway will start forwarding connection and health check requests to the RegistryProxy service. (More on this in the Getting Started Guide).

- The idl folder just contains the IDLs for the RPC Services and Gateway we made when developign this project. **_Please note that this is NOT the IDL Management system. These files are only provided for reference_**\_ Thus, the gateway will function even if the idl folder is removed.

_Please note that certain details in this project have been "mocked" during development to simplify testing and expedite the process keeping in mind the architecture. However, the API Gateway is designed to be functional, scalable, and modular, ensuring it can accommodate future updates and enhancements seamlessly. Despite the mocked data, the implementation follows best practices and adheres to the intended functionality, allowing for efficient communication between services and robust handling of incoming requests._

**For better testing of gateway, the RPC services may return a static response. This was done on purpose. However, the code can be uncommented to return dynamic responses.**

PS :
On Windows, Kitex may throw some errors like :

```
[Error] KITEX: OnRead Error: default codec read failed: EOF
default codec read failed: i/o timeout
```

This happens to be due to the Windows environment. On running Kitex servers in Linux, the errors seem to go away.
The issue has been discussed in [here](https://github.com/cloudwego/kitex/issues/932) and [here as well](https://github.com/cloudwego/kitex/issues/964)

<a href="#top">Back to top</a>

## Features and Design <a name="features"></a>

### Architecture Diagram <a name="diagram"></a>

![gateway_design](gateway_design.png)

### The project has the following features : <a name="components"></a>

1. API Gateway Server: The API Gateway is implemented as a Hertz server that listens to requests on port 4200. It exposes multiple endpoints in the format `/{serviceName}/{path}` for both POST and GET requests. The API Gateway acts as an intermediary between user requests and the Kitex RPC servers performing load balancing, service discovery, health checks, etc. It accepts incoming HTTP requests and decodes on which service to make a generic RPC call and returns a response.

2. Service Registration and Caching: A user can send in POST requests on `/register` to register their services. The JSON payload includes details about the owner, their services, IDLs, and masked service URLs. This data is stored in MongoDB and is also cached inside the gateway server for quick access. Refer to the <a href="#data">Data</a> section below for the exact info on how data is stored.

3. RPC Protocol Translation: The API Gateway forwards incoming API requests to the Kitex servers using the internal RPC clients within the Hertz server. It enables communication between the client and the respective RPC servers responsible for handling specific services. This is done with their Thrift IDL information sent on registration.

4. Automated Server Connection: The API Gateway provides another package called _server_utils_ to make it easy for Services can automate the registration of their servers to our system by making an HTTP POST request at the `/connect` endpoint of the API Gateway using their registered API key. This enables services to scale up or down dynamically according to their needs.

5. Health Checks: Servers connected to the API Gateway need to declare their health by making periodic requests to the `:/health` endpoint at least every 10 seconds. This ensures that the API Gateway considers the servers healthy and forwards requests to them. This is part of the Server Utility Package and can be easily automated. For this, we have used Consul.

6. Discovery and Load Balancing: The gateway uses Consul's resolver for service discovery and load balancing. It uses Consul’s DNS interface to resolve service names. It can be used to discover services registered with Consul and load balance requests between them. The MVP version of the NewConsulResolver implements round-robin load balancing, distributing the requests equally among the connected RPC servers.

7. Service Registry: Consul has been integrated as the service registry as well. It is completely isolated from the RPC servers. The servers have to, therefore, interact with the Gateway for Consul-related business. This was done to prevent malicious attacks on the service registry. The consul service registry provides a beautiful graphical view of all the connected servers to help admins track, manage and troubleshoot connections. This also allows us more freedom in logic handling for tasks related to the Registry.

- Consul Agent is for Service Registry, Health Checks, Discovery, and Load Balancing. It is hosted on `localhost:8500` for this project. The logic for exact functioning however has been coded manually.

8. Registry Proxy Server: An RPC server that can perform health checks, on behalf of the gateway has been included in the project as well. This RPC server is a special server that has direct access to the Consul Service Registry. When servers ping the `/health` or `/connect` endpoint of the gateway, the gateway can proxy handling this request to this RPC service (if one or more of the Registry Proxy servers are online). This frees up space and resources so that the gateway can handle other requests. This is an optional server that the admin may decide to boot up; if the gateway detects it's offline, then it will perform a health check and connection requests itself.

Note:

- This service is kept optional because it may be a bottleneck if only a few servers are making requests to the `/health` or `/connect` endpoints. It is advised to boot this service up only when there are many servers connected to the system.
- There may be multiple instances of this service running at the same time. A round robin load balancing is used to determine connection.
  <a href="#top">Back to top</a>

### IDL Management <a name="idlmanagement"></a>

IDL (Interface Definition Language) management plays a crucial role in facilitating communication between the API Gateway and the backend RPC servers using Kitex. Here's how IDLs are managed in the project:

During service registration with POST /register, service owners include their service's IDL details in the JSON payload.

- IDL Mapping and Translation:
  API Gateway stores and uses the IDL contents to translate incoming JSON API requests to Thrift binary format for Kitex.
  The IDLs are also cached for faster API Calls and are used to make and store all the generic clients so that the parsing need to be done again and again. Thrift's compact binary format facilitates efficient communication between the API Gateway and backend RPC servers.

- IDL Versioning and Compatibility:
  The Data Model allows services to keep track of their IDLs and versions. Although one limitation (due to time constraints) is that only the latest IDL version is stored. IDLs can be updated using the `/update` endpoint.

- Synchronization:
  Ensuring IDL information remains consistent and synchronized between components is crucial, which is why IDLs are cached upon service registration/updatation.

  <a href="#top">Back to top</a>

### Server Utility Package <a name="serverutils"></a>

The Server Utility Package provides convenient functions to interact with the API Gateway for automated server registration, health checks, and communication. Importing this package simplifies the process of connecting backend RPC servers to the API Gateway.

1. Importing the Package:

To use the Server Utility Package, import it into your main package or copy-paste the file contents directly into your main package.

2. Example Usage:

The package exports 3 methods which are needed to start your Kitex server :

```
func NewGatewayClient(apikey string, serviceName string, gatewayAddress string) *GatewayClient
```

The above function makes a new client to talk to the gateway server

```
func (client *GatewayClient) ConnectServer(serverAddress string, serverPort string) (string, error)
```

Use the above method to connect the server to the gateway server. Methods needs to be called on the gateway client.

```
func (client *GatewayClient) UpdateHealthLoop(id string, timeBetweenLoops int)
```

The above methods keeps declaring server instance is healthy. Methods needs to be called on the gateway client.
<br>
The following code snippet demonstrates how to use the Server Utility Package to connect a backend RPC server to the API Gateway and perform health checks:
(The full usage can be found in rpc_services/asset_service)

```

func main() {
  // Other code above

	// Initialize Gateway Client with API key and service name
	gatewayAddress := config.GatewayAddress
	gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName, gatewayAddress)

	// Connect to the API Gateway and get the server ID
	id, err := gatewayClient.ConnectServer(config.ServiceURL, advertisedPort)
		if err != nil {
			log.Fatal(err.Error())
	}

	// Perform health checks in a separate goroutine (every 5 seconds)
	go gatewayClient.UpdateHealthLoop(id, config.HealthCheckFrequency)

	// Your server's main logic here
}
```

<a href="#top">Back to top</a>

## Data Management<a name="data"></a>

Data is stored in MongoDB for this Project for ease of use and flexibility. The MongoDB client is configured to connect to Mongo at `localhost:27107`

For each unique owner, a document is made in MongoDB, with the following data :

(The API Key is provided by the gateway. Rest is provided by the user.)

```
{
  "ApiKey": "your-api-key", // Generated and Appended by the Gateway. Unique for each user. Master api key used for testing is "master_api_key_uuid"
  "OwnerName": "John Doe", // Name of the Admin
  "OwnerId": "user123", // Unique user name for the Admin
  "Services": [ // The various services registered by the admin
    {
      "ServiceId": "service1", // Identifier for each service
      "ServiceName": "Service One", // Name of service. Used for making API calls as well
      "IdlContent": "<Thrift IDL content for Service One>", // IDL Content which helps define the structure of the service's API
      "Version": "1.0.0", // Version of the Service
      "ServiceDescription": "This is Service One description.", // Description about the service
      "ServerCount": 3, // Number of servers. Mocked for this project, but accounted for. We can connect as many as required with the correct API key or master key
      "Paths": [ // Paths provide flexibility in the URL for each service.
        {
          "MethodPath": "/method1", // This is the method which the path should correspond to
          "ExposedMethod": "GET" // This is the path in the URL "gateway/{serviceName}/{path}"
        },
        {
          "MethodPath": "/method2",
          "ExposedMethod": "POST"
        }
      ],
      "RegisteredServers": [ // These are the registered servers which are allowed to connect to the gateway
                            // Again, this has been mocked, but accounted for. All servers are allowed to connect with any URL as long as API key is correct.
        {
          "ServerUrl": "http://server1.example.com",
          "Port": 8080
        },
        {
          "ServerUrl": "http://server2.example.com",
          "Port": 8080
        },
        {
          "ServerUrl": "http://server3.example.com",
          "Port": 8080
        }
      ]
    },
    {  // We can register more services here
      "ServiceId": "service2",
      ... and so on
    }
  ]
}
```

The IDL must be provided by stringifying it. A tool like https://jsonformatter.org/json-stringify-online can be used for this.

A user can register multiple services, and multiple registered servers for their services and along with some flexibility in exposed URLs, the method is masked with Path field.

_For ease of testing, regardless of how many Registered Servers are there, we can connect more, and with different ServerURLs. "Mocking" the authentication of RPC servers this way will save time on testing. This is done with the master api key, which can be used instead of the generated api key. While the gateway can perfectly function with the generated api key as well, the master api key makes the testing simpler._

<a href="#top">Back to top</a>

## Getting Started with an example <a name="gettingstarted"></a>

#### 1. Build and run the Hertz server <a name="step1"></a>

Run `go build; go run .` to build and run the API Gateway server.

To check if the server is running, hit the following GET endpoint
`"http://0.0.0.0:4200/ping"`

It should reply with the message :

```
{
    "message": "pong"
}
```

#### Settings for gateway <a name="hertzsettings"></a>

The gateway has a `serverconfig.json` file which the admin may change if he wants to change some of the settings.

```
{
  "serverPort": "4200", // Host port of the gateway
  "maskerKey": "master_api_key_uuid", // Specify any master api key (can be used if the masterkey gets leaked)

  "consulAddress": "127.0.0.1:8500", // Address of the consul service registry
  "timeToLift": 45, // Time in seconds before a service is declared unhealthy, if it fails to declare itself healthy by pinging the :/health endpoint
  "timeToDie": 90, // Time in seconds before a service is removed from the registry, if it fails to declare itself healthy by pinging the :/health endpoint

  "dbUrl": "mongodb://localhost:27017", // Url of the mongoDb used for IDL and Service Management
  "dbName": "api_gateway_db", // Name of the mongoDb database
  "dbCollectionName": "admin_services", // Name of the collections
  "dbPingInterval": 30, // How frequently gateway should ping the mongoDb database
  "dbFailedPingInterval": 5, // How frequently gateway should ping the mongoDb database if it pinging has failed once or more
  "dbMaxPingFailDuration": 180 // After how many tries to stop the pinging
}
```

- If you want to change the consul service registry address, change the `consulAddress` field here (api-gateway/hertz_server/`serverconfig.json`) **and** the `consulAddress` field in the `serviceConfig.json` file inside the /rpc_services/registry_proxy_service directory.

<a href="#top">Back to top</a>

#### 2. Start Consul <a name="step2"></a>

Assuming Consul is already installed, run `consul agent -dev`. This will start Consul on `localhost:8500` with a beautiful GUI of all Services connected.

![consul](consulservicesconnection.png)

<a href="#top">Back to top</a>

#### 3. Start MongoDB <a name="step3"></a>

Start MongoDB server on `localhost:27017`

<a href="#top">Back to top</a>

#### 4. Setup up Kitex with the Server utility package as mentioned above. <a name="step4"></a>

Here's an example of setting up a Kitex server with Docker :

```
  config, err := LoadConfiguration("serviceConfig.json")
	if err != nil {
		log.Fatal(err)
	}

  gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName)
	advertisedPort := os.Getenv("PORT")

	advertisedPort = GetFreePort()

	id, err := gatewayClient.connectServer(config.ServiceURL, advertisedPort)
    if err != nil {
    log.Fatal(err.Error())
	}

	go gatewayClient.updateHealthLoop(id, 5)

	url := config.URL
	port := config.Port
	addrDocker, _ := net.ResolveTCPAddr("tcp", url+":"+port)

	svr := asset_management.NewServer(new(AssetManagementImpl),
	server.WithServiceAddr(addrDocker),
  server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
	server.WithReadWriteTimeout(100*time.Second))

	kitexerr := svr.Run()

	if kitexerr != nil {
	log.Println(kitexerr.Error())
	}
```

And on LocalHost :

```
  config, err := LoadConfiguration("serviceConfig.json")
    if err != nil {
      log.Fatal(err)
    }

  var addr = getAddr() // Function in utils.go

	gatewayClient := NewGatewayClient(config.Apikey, config.ServiceName)

	id, err := gatewayClient.connectServer(addr.IP.String(), strconv.Itoa(addr.Port))
	if err != nil {
		log.Fatal(err.Error())
	}

	go gatewayClient.updateHealthLoop(id, 5)

	svr := asset_management.NewServer(new(AssetManagementImpl),
		server.WithServiceAddr(addr),
		server.WithLimit(&limit.Option{MaxConnections: 100000, MaxQPS: 100000}),
		server.WithReadWriteTimeout(100*time.Second))

	kitexerr := svr.Run()

	if kitexerr != nil {
		log.Println(kitexerr.Error())
	}
```

### Kitex Settings <a name="kitexsettings"></a>

The RPC services (asset service and user service) we provided have a `serviceConfig.json` file where the server configurations are kept. This is provided in case the admin changes things around, they can edit the `serviceConfig.json` file instead of changing source code. This makes sure that the server can be launched with different settings without having to rebuild the project or the docker image. These settings are used to start up RPC server properly :

```
{
    "apikey": "master_api_key_uuid", //apikey to be used for server connection
    "ServiceName": "AssetManagement", //name of the service
    "dockerUrl": "0.0.0.0", //Address of gateway in docker
    "dockerPort": "8080", //Port of the gateway in docker
    "env": "Stage", //
    "serviceurl": "localhost", //Url of the server.
    "serverPort": "0", //Port to be used by server. Keep 0 for random port.
    "IsDockerised": false, //Change to true if docker is being used.
    "healthCheckFrequency": 15, //How frequently, in seconds, do you want to perform health checks?
    "gatewayAddress": "http://localhost:4200" //address of the gateway.
}
```

To launch the Kitex server in Docker: 1. Change `IsDockerised` to `true` 2. Change `gatewayAddress` to the new address. Eg. "`
http://host.docker.internal:4200`" i.e., if gateway is on the same machine hosted locally.

During production, the gatewayAddress would be updated to the absolute URL of the API Gateway and can be abstracted within the Server utility package to avoid hardcoding.

_This has been implemented in the Asset Service, feel free to use that as an example.(the main.go file inside the Asset Service)_

<a href="#top">Back to top</a>

#### 5. Register a service <a name="step5"></a>

Send a post request to `0.0.0.0:4200`. Let's say we want to register the Asset Service (in the rpc_services) folder. Therefore we'd send a request with the following JSON :

```
[
    {
      "OwnerName": "John Doe",
      "OwnerId": "UserName",
      "Services": [
        {
          "ServiceId": "1",
          "ServiceName": "AssetManagement",
          "IdlContent": "namespace Go asset.management\n\nstruct QueryAssetRequest {\n    1: string ID;\n}\n\nstruct QueryAssetResponse {\n    1: bool   Exist;\n    2: string ID;\n    3: string Name;\n    4: string Market;\n}\n\nstruct InsertAssetRequest {\n    1: string ID;\n    2: string Name;\n    3: string Market;\n}\n\nstruct InsertAssetResponse {\n    1: bool Ok;\n    2: string Msg;\n}\n\nservice AssetManagement {\n    QueryAssetResponse queryAsset(1: QueryAssetRequest req);\n    InsertAssetResponse insertAsset(1: InsertAssetRequest req);\n}\n",
          "Version": "1.0",
          "ServiceDescription": "Service Description",
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

Note:

The above service would have the following public endpoints registered on the gateway:

- (POST) `/AssetManagement/newAsset` which maps to the private "insertAsset" endpoint of the service
- (GET) `/AssetManagement/getAsset` which maps to the private "queryAsset" endpoint of the service

A response with the API Key and Status will be recieved.

### Another example for registration :

To register a service, you can simply make a POST request to the `/register` endpoint. The json payload should look like this: (Here we are registering our RegistryProxy service)

```
[
	{
	"OwnerName": "XXXXX",
	"OwnerId": "XXXXX",
	"Services":
	[
		{

		"ServiceId": "X",
		"ServiceName": "RegistryProxy",
		"IdlContent": "namespace Go registy.proxy\n\nstruct ConnectRequest {\n 1: string ApiKey;\n 2: string ServiceName\n 3: string ServerAddress\n 4: string ServerPort\n 5: i64 TTL\n 6: i64 TTD\n}\n\nstruct ConnectResponse {\n 1: string Status;\n 2: string Message;\n 3: string ServerID;\n}\n\nstruct HealtRequest {\n 1: string ApiKey;\n 2: string ServerID;\n}\n\nstruct HealthResponse {\n 1: string Status;\n 2: string Message;\n}\n\nservice RegistryProxy {\n ConnectResponse connectServer(1: ConnectRequest req);\n HealthResponse healthCheckServer(1: HealtRequest req);\n}\n",

		"Version": "1.0",
		"ServiceDescription": "Service for proxying health checks and coneection requests",

		"ServerCount": 0,
		"Paths": [
			{
			"ExposedMethod": "healthCheckServer",
			"MethodPath": "healthCheckServer"
			},
			{
			"ExposedMethod": "connectServer",
			"MethodPath": "connectServer"
			}
		],
		"RegisteredServers": []
		}
	]
	}
]
```

<a href="#top">Back to top</a>

#### 6. Sending requests <a name="step6"></a>

Send a POST or GET requests to the "/{serviceName}/{path}", in this case `http://localhost:4200/AssetManagement/newAsset` endpoint, for example:

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

<a href="#top">Back to top</a>

#### 7. Update your service <a name="step7"></a>

This step can actually be done anytime after registration but placing it here made sense.

For updating, put in the api key received on registration. OR, for easier testing we made the provision that if a service tries to updated using a master key the update will still be authorised. Further add the correct owner ID in the parameters. So the request URL should look something like this :

`http://localhost:4200/update?ownerid=UserName`

The master api key is "master_api_key_uuid" without the quotes.

_The Master Key is just a temporary provision made for easy testing and WILL be removed in the future. Managing API keys can be quite a hassle._

Try updating with the new request body :

```
[
    {
      "OwnerName": "John Doe Updated",
      "OwnerId": "UserName",
      "Services": [
        {
          "ServiceId": "1",
          "ServiceName": "AssetManagement",
          "IdlContent": " ",
          "Version": "5.0",
          "ServiceDescription": "Service Description Updated",
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

With the above updated blank IDL, we try sending a request now, and but it will not be authorised by the gateway.

In MongoDB we can notice the updated changes.

A provision for getting back information for an Admin has not yet been implemented due to time constraints of the project. But it would certainly be a great feature to have.

<a href="#top">Back to top</a>

#### 8. Setup Registry Proxy Server(s) **[Optional]** <a name="step8"></a>

If the server load is getting too high and many rpc servers are connected, you may decide to connect a special RPC server we made, the Registry Proxy Service. The purpose of this special RPC server is to allow the gateway to proxy all the health check requests/server connection requests from different servers so that the gatway can have resources to handle more service requests. By adding this server, we were able to ramp up performance from ~_2600 req/s_ to ~_3000 req/s_ for 50 users and 3 rpc servers.

You may setup this server by:

1. Register this service in the gateway. You may do so by sending a `POST` request to `/register` endpoint as such:

```
[
    {
        "OwnerName": "XXXXX",
        "OwnerId": "XXXXX",
        "Services": [
            {
                "ServiceId": "X",
                "ServiceName": "RegistryProxy",
                "IdlContent": "namespace Go registy.proxy\n\nstruct ConnectRequest {\n    1: string ApiKey;\n    2: string ServiceName\n    3: string ServerAddress\n    4: string ServerPort\n    5: i64 TTL\n    6: i64 TTD\n}\n\nstruct ConnectResponse {\n    1: string Status;\n    2: string Message;\n    3: string ServerID;\n}\n\nstruct HealtRequest {\n    1: string ApiKey;\n    2: string ServerID;\n}\n\nstruct HealthResponse {\n    1: string Status;\n    2: string Message;\n}\n\nservice RegistryProxy {\n    ConnectResponse connectServer(1: ConnectRequest req);\n    HealthResponse healthCheckServer(1: HealtRequest req);\n}\n",
                "Version": "1.0",
                "ServiceDescription": "Service for proxying health checks and coneection requests",
                "ServerCount": 0,
                "Paths": [
                    {
                        "ExposedMethod": "healthCheckServer",
                        "MethodPath": "healthCheckServer"
                    },
                    {
                        "ExposedMethod": "connectServer",
                        "MethodPath": "connectServer"
                    }
                ],
                "RegisteredServers": [
                ]
            }
        ]
    }
]
```

2. Booting up and instance of registry proxy from the rpc_services provided by running `go run .` in the /rpc_services/registry_proxy_service directory.

Shortly after the server is booted up, the gateway will detect the server and start to proxy health check requests and server connection requests to this registry proxy server, freeing up more resources for the gateway to handle other requests.

<a href="#top">Back to top</a>

## Performance <a name="perf"></a>

#### Current Performance <a name="currentperf"></a>

On Load testing with Apache JMeter, we were abe to get the following benchmarks

![performance1](perf-50users-final.png)
3 instances of AssetManagement, 1 registry proxy server (for proxied health checks)

- Users : 50
- Total time: 3 mins
- Ramp up time : 1sec

#### MVP Performance : <a name="mvpperf"></a>

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

<a href="#top">Back to top</a>

## Testing <a name="tests"></a>

While there are not many unit tests, the project has been tested for functionality and integration.

These testcases can be imported in postman for ease and can be found here : [GatewayTests.json](./GatewayTests.json)

## Limitations/Issues<a name="limit"></a>

On thorough testing we found some limitations such as :

- Kitex Servers cannot reconnect if the Gateway server goes down, even for a second. The server_utility package needs to be updated for this.
- While the Gateway is designed to be Scalable, the only non-scalable aspect as of now is the Data Management with the Database and service registry.

<a href="#top">Back to top</a>

## What Else? <a name="misc"></a>

- Nothing for now :p. Hope you liked this project.

<a href="#top">Back to top</a>
