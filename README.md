This is a repository for the ByteDance and Tiktok Orbital 2023.

## About

This is the MVP project for our API Gateway based on one Hertz server and multiple RPC servers.

The API Gateway, which is a Hertz server, listens to requests at port 4200 on multiple exposed endpoints "/{serviceName}/{serviceMethod}" [POST] and "/{serviceName}/{serviceMethod}" [GET]. Once it receives an API request, it then forwards the request to the Kitex server (using the internal RPC client built inside the Hertz server). The user service is at port 8888 while the Asset Management service can be initialised on any port from user input from console.

#### Behind the scenes :

For now, the registration functionality is NOT integrated with the service registry. However, services can be registered at the (POST) `:/register` endpoint with the following json format which is accepted with our service registry Consul as well :

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

Moreover, to see all registered service we can hit the endpoint `:/show`. (This will be removed later and was made just for testing purposes)

## Performance

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

**Step 4:**

Send a POST or GET requests to the "/{serviceName}/{path}" endpoint, for example:

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
curl -X GET http://localhost:4200/AssetManagement/queryAsset?ID=2
```

- The "ID" should be capital as it's case sensitive.
