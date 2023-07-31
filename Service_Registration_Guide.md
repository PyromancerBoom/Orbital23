# Service Registration Guide

To register a service, you can simply make a POST request to the `:/register` endpoint. The json payload should look like this: (Here we are registering our RegistryProxy service) 
```json
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

## Structure

### Outer Structure
* OwnerName : Name of owner. Any *string*.
* OwnerId: ID of owner. Any *string*.
* Services: Services to be registered. *Array* of Services.

### Services Structure
* ServiceId: ID of service. Any *string*. For future features.
* ServiceName: Name of Service. Has to be same as IDL. *string*.
* IdlContent: `.thrift` IDL content in *string* format.
* Version: Version of the service. Any *int*. For future features.
* ServiceDescription: A brief description. Any *string*. For future features.
* ServerCount: Number of servers. Any *int*. For future features.
* Paths: The paths for the methods of this service. *Array*.
* Path -> ExposedMethod: The actual method name (method name in the idl). *string*.
* Path -> MethodPath: The method you want users to request to. Any valid GoLang method name. *string*.
* RegisteredServers: Info for registered servers. For future features.
* RegisteredServers -> ServerUrl: Url of a server for this service. For future features. *string*.
* RegisteredServers -> Port: Port of a server for this service. For future features. *int*.

### Notes:
* Paths: 
	* You must make sure that you mention all the paths of your functions you mentioned in the idl. Only the paths mentioned here will be able to handle requests. If the path is not mentioned here but the method exists in the idl, no requests will be made to that method. 
	* If you do not want path masking, keep the `ExposedMethod` and `MethodPath` identical.
	* There should be at least one `ExposedMethod` -`MethodPath` pair for each method. More than one is allowed as well. However, if you don't specify a `MethodPath` for a method, that method will not receive requests.
	* User is advised not to use the same `ExposedMethod` name for different methods. It is uncertain where the requests will be routed then.
