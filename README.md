This is a repository for the ByteDance and Tiktok Orbital 2023.

## About

This is the MVP project for our API Gateway based on one Hertz server and multiple RPC servers.

The API Gateway, which is a Hertz server, listens to requests at port 4200 on multiple exposed endpoints "/{serviceName}/{serviceMethod}" [POST] and "/{serviceName}/{serviceMethod}" [GET]. Once it receives an API request, it then forwards the request to the Kitex server (using the internal RPC client built inside the Hertz server). The user service is at port 8888 while the Asset Management service can be initialised on any port from user input from console.

## How to use? [^3]

**Step 1:**

Initialise the any of the Kitex servers using the command: "go run ." from the respoective directory

**Step 2:**

Initialise the Hertz server using the command: "go run ." from the "./hertz_server" directory

To check if the server is running, hit the following endpoint with a curl request

```
curl -X GET "http://localhost:4200/ping"
```

It should reply with the message :

```
{
    "message": "pong"
}
```

**Step 3:**

Send a POST or GET request to the "/{serviceName}/{serviceMethod}" endpoint, for example:

```
curl -X POST -H "Content-Type: application/json"
-d '{
	"ID":"123",
	"Name":"John Doe",
	"Email":"john.doe@example.com",
	"Age":30
	}'
"http://localhost:4200/UserService/insertUser"
```

Currently the project is using a makeshift Hashmap to keep track of IDL mappings. But that would be integrated into the service registry later on.

Further, dyanmic routing is NOT supported as of now and will be coming in the future updates.
