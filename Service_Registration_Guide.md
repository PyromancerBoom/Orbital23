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
		"RegisteredServers": []
		}
	]
	}
]
```

### Structure
