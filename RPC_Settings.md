# RPC Server Settings
The RPC services (asset service and user service) we provided have a `serviceConfig.json` file where the server configurations are kept. This is a json file which is read by the RPC server on bootup to set things up easily. This is provided in case the admin changes things around, they can edit the `serviceConfig.json` file instead of changing source code.

### File structure
```json
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
### Notes
-  If you want to dockerise the gateway:
	1. Change `IsDockerised` to `true`
	2. Change `gatewayAddress` to the new address. Eg. "`
http://host.docker.internal:4200`"
