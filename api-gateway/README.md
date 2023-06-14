## About

***Note : As of Milestone-1 and for the POC, we have used a dedicated client to communicate with kitex server. However, for the main project, we plan to use the Generic call feature inside the hertz API Gateway server. This server would forward requests various GET, POST, DELETE, etc requests based on the service Id and other details fetched from the service registry to the target backend service.***

This is the starter project for our API Gateway based on the POC we made and consists of communication[^1] between one Hertz server and one RPC server. The Hertz server is generated using the asset_api.thrift IDL file, while the Kitex server and client are generated using the asset_management.thrift IDL file.

The Hertz server listens to requests at port 4200 on two exposed endpoints at "/asset/insert" [POST] and "/asset/query" [GET]. Once it receives an API request, it then forwards the request to the Kitex server (using the internal Kitex client built inside the Hertz server). The Kitex server sits on port 8888 and responds to the RPC calls made to it.

## Endpoints

| Endpoint      | Method | Description                                                                                           |
| ------------- | ------ | ----------------------------------------------------------------------------------------------------- |
| /asset/query  | GET    | Used to query about an asset, with its `id` speficied in the url query section                        |
| /asset/insert | POST   | Used to insert an new asset into the RPC database[^2] . Usage can be inferred from the tutorial below |

## How to use? [^3]

**Step 1:**

Initialise the Kitex server using the command: "go run ." from the "./kitex_server" directory. Enter host port as "8888"

**Step 2:**

Initialise the Hertz server using the command: "go run ." from the "./hertz_server" directory

**Step 3:**

Send a POST request to the "asset/insert" endpoint by:

```
curl --location --request POST 'http://127.0.0.1:4200/asset/insert' \
--header 'Content-Type: application/json' \
--data-raw '{
"Name": "APPLE",
	"ID": "1",
	"Market": "US"
}'
```

**Step 4:**

Send a GET request to the "asset/query" endpoint by:

```
curl --location --request GET 'http://127.0.0.1:4200/asset/query?id=1'
```

[^1]: Kitex server is using the port "8888" and Hertz server using the port "4200", so please keep these ports free for the demo servers to run.
[^2]: We are not current using an actual database in the demo. The data structure used is a go splice which acts as a makeshift database.
[^3]: It is assumed that go is already installed in your system

by : Tahsin Hasem and Priyansh Bimbisariye