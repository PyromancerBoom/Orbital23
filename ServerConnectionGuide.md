# Server Connection Guide:

Once your service is registered under our system, you will receive an api-key which you can use to register your servers in our system. There is currently no limit to the number of servers you can enlist to our system. However, our system only supports RPC servers as the gateway merely acts as a proxy between HTTP network requests to RPC calls through thrift binary encoding.

Note: Since service registration is not supported right now in the MVP, only AssetManagement and UserService services can be registered into the system. However, after the implementation of service registration through the `:/register` endpoint, servers of any registered services may be connected.

## Requirements to connect to our system:

1. Service Owner must ensure his servers are able to perform all the methods he indicated in his interface definition during the service registration in the `"/register` endpoint. We shall only forward the requests, it is his/her duty to ensure his/her servers are capable of performing the defined functions.

2. Service Owner must ensure that every time a new server is booted up, it registers itself via the `:/connect` endpoint

3. Service Owner must ensure that his servers declare themselves healthy to our system so we may forward requests to his/her RPC server(s) by making requests to `:/health` endpoint. Our system will no longer consider the server healthy if it does not declare itself healthy (at least) every 10 seconds. Server will be delisted after 1 minute of the last health check.

## Guide To Connect your server:

### Step 1:

Register your server in our system via `:/register` endpoint. You will have to provide a Service Name, which you will need to add servers to our system. After registering your service, you will get an API Key. That key will let you register servers to our system.

### Step 2:

On (RPC) server bootup, send request to `:/connect` [HTTP POST] endpoint with json body containting:

* api-key : The API key you received when you registered your service.
* serviceName: The exact name you used when you registered your service.
* serverAddress: The address of the server you wish to connect to our system.
* serverPort: The port of the server which will be used for communication with our system.

**Example Request body:**
```json
{
"api-key":"36e991d3-646d-414a-ac66-0c0e8a310ced",
"serverAddress":"127.0.0.01",
"serverPort":"9999",
"serviceName":"UserService"
}
```

If the request is successful, you will receive a serverID which you must use to continually declare your server's health.

**Example Response:**
```json
{
"message": "Server Connection Request Accepted.",
"serverID": "b7b5e972-9aa7-4e82-95d7-57876ac9b69f",
"status": "status OK"
}
```
### Step 3

Make the (RPC) server declare it is online and healthy by making http requests to the `:/health` [HTTP POST] endpoint at least every 10 seconds (recommended 5 seconds). The request body must contain:

* api-key: The API key you received when you registered your service.
* serverID: The Server ID received when registering the server to the system via `:/connect` endpoint (STEP 2).

**Example Request body:**
```json
{
"api-key":"36e991d3-646d-414a-ac66-0c0e8a310ced",
"serverID":"b7b5e972-9aa7-4e82-95d7-57876ac9b69f"
}
```
**Example Response:**
```json
{
"status": "status OK",
"message": "Successfully Updated the healh of server"
}
```

**Note:**
Time before server declared offline: 10 seconds
Time before server is delisted from sstem: 60 seconds

If server is delisted/is not listed, you will get an Internal Server Error response with the message "Unable to process health update request." when trying to update the server's health.

If server is delisted, you must reconnect it to the system via follwing Step 2 and Step 3.

### Limitations
* Currently, the MVP version does not support service owner to query a list of all of his/her servers that are registered/passing-health-check in our system. It will be included in the final version.
* Currently, a round-robin load balancing is being used. It will be updraged to weighted round-robin in future updates.
