### About the RPC Services

For the testing of this project, the RPC services return a static output.

This is a design decision so that the actual efficiency of the API Gateway could be calculated

### Running the services

If you run the services on a windows environment there are chances you may encounter some errors :

1. Default codec read failed: i/o timeout:
   This error occurs when there is a timeout while reading from the TCP connection during an RPC call. It could happen due to network issues or if the remote service does not respond within the specified timeout.

2. Default codec read failed: EOF:
   The "EOF" error indicates that the end of the file or stream was reached unexpectedly during an attempt to read. This error can occur due to issues with the network or the connection being closed unexpectedly.

This occurs due to some settings in the windows environment. When the services are run on a Linux system these errors are resolved. In this project, the kitex servers are ocnfigured to run on Docker by default.

#### Using Docker to run the services

1. Make an image : `docker build -t <service_name> .`
   For ease of use, set the service_name to the same name as it's folder.

2. Modify `docker-compose.yml` if needed for more instances or different ports.

3. Run docker compose with `docker-compose up` or `docker-compose up -d` (Latter is if you don't want each service's cli in your terminal)

#### Important note before getting started

In this project, the kitex servers are hosted on docker and each instance of the service needs to be specified separately in the docker compose file.
This is so because, there may be a need to modify an instance separately and to give more flexibilty in ports.

Inside each folder three files are to be taken note of :

1. `config.json`, which looks something like this :

```
{
    "url": "0.0.0.0",
    "port": "8080",
    "env": "Dev",
    "serviceurl": "localhost"
}

```

Change serviceurl to same network as Gateway's network.
The fields above imply various conmfigurations to run the Kitex server. `url` is the host url of each kitex instance, `port` is the port, `env` is just some info for logs, `servicurl` is the url sent to Docker to communicate. When hosting locally, do NOT change the `"serviceurl": "localhost"` ` field.

2. `Dockerfile`: Dockerfile for the service

3. `docker-compose.yml`: This is where the details of the instances are specified. This file also fetches information from `config.json` to make it easy to initialise a Kitex server without having to build the code (and the docker image) again and again. Therefore, for example, we can simply change the url in the `config.json` and run the docker container with `docker-compose up` without having to re-build the image or the project.
