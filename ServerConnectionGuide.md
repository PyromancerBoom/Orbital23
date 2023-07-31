# Server Connection Guide:

Once your service is registered under our system, you will receive an api-key which you can use to register your servers in our system (or use master api-key). There is currently no limit to the number of servers you can enlist to our system. However, our system only supports kitex RPC servers as the gateway merely acts as a proxy between HTTP network requests to RPC calls through thrift binary encoding.

## Requirements to connect to our system:

1. Service Owner must ensure his servers are able to perform all the methods he indicated in his interface definition during the service registration in the `"/register` endpoint. We shall only forward the requests, it is his/her duty to ensure his/her servers are capable of performing the defined functions.

2. Service Owner must ensure that every time a new server is booted up, it registers itself via the `:/connect` endpoint

3. Service Owner must ensure that his servers declare themselves healthy to our system so we may forward requests to his/her RPC server(s) by making requests to `:/health` endpoint. Our system will no longer consider the server healthy if it does not declare itself healthy after a period of time. Server will be delisted after some time (mentioned in the serverconfig.json file in the gateway) if the server fails to make health check requests.

## Guide To Connect your server:

### Step 1:

Register your server in our system via `:/register` endpoint. Guide can be found ![here](/Service_Registration_Guide.md)

### Step 2:

On (RPC) server bootup, send request to `:/connect` [HTTP POST] endpoint with json body containting:

* ApiKey : The API key you received when you registered your service.
* ServiceName: The exact name you used when you registered your service.
* ServerAddress: The address of the server you wish to connect to our system.
* ServerPort: The port of the server which will be used for communication with our system.

**Example Request body:**
```json
{
"ApiKey":"36e991d3-646d-414a-ac66-0c0e8a310ced",
"ServerAddress":"127.0.0.01",
"ServerPort":"9999",
"ServiceName":"UserService"
}
```

If the request is successful, you will receive a serverID which you must use to continually declare your server's health.

**Example Response:**
```json
{
"Message": "Successfully connected server to gateway.",
"ServerID": "b7b5e972-9aa7-4e82-95d7-57876ac9b69f",
"Status": "ok"
}
```
### Step 3

Make the (RPC) server declare it is online and healthy by making http requests to the `:/health` [HTTP POST] endpoint periodically. The request body must contain:

* ApiKey: The API key you received when you registered your service.
* ServerID: The Server ID received when registering the server to the system via `:/connect` endpoint (STEP 2).

**Example Request body:**
```json
{
"ApiKey":"36e991d3-646d-414a-ac66-0c0e8a310ced",
"ServerID":"b7b5e972-9aa7-4e82-95d7-57876ac9b69f"
}
```
**Example Response:**
```json
{
"Status": "ok",
"Message": "Successfully updated server health."
}
```

**Note:**

- Time before server declared offline: specified as `TTL` in ![serverconfig.json](/api-gateway/hertz_server/serverconfig.json) 
- Time before server is delisted from sstem: specified as `TTD` in ![serverconfig.json](/api-gateway/hertz_server/serverconfig.json) 

If server is delisted/is not listed, you will get an Internal Server Error response with the message "Unable to process health update request." when trying to update the server's health.

If server is delisted, you must reconnect it to the system via follwing Step 2 and Step 3.

*We are providing a `server_utils.go` utils file that you may use to conveniently achieve this functionality just by making a few fucntion calls and avoiding all this hassle. Refer to ![README.md](README.md)*
