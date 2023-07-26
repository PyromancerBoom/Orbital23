### Server Utility Package

This is a package required to be used by Kitex servers that want to be hosted on the API Gateway.

Simply import the package (or copy paste the file in the main package)

An example of usage is as follows :

```
gatewayClient := NewGatewayClient(configuration.Apikey, configuration.ServiceName)

advertisedPort := os.Getenv("PORT")

id, err := gatewayClient.connectServer(configuration.ServiceURL, advertisedPort)
if err != nil {
	log.Fatal(err.Error())
}

go gatewayClient.updateHealthLoop(id, 5) // Health Checks must be in a separate Go routine
```
