## Implemented features:
* Generic JSON Mapping
* Service Registry

## Errors
* Kitex throwing EOF error, as shown in the picture
* Erros in communication with Registry over the internet (clow priority for resolving)
* Kitex throws codec errors during generic call communication, however, the correct output is retrieved

## ToDo
* Finaliase Implementation of Service Discovery
* Try out a simplr Random Load Balancer
* Perform benchmarking on single Hertz server
* Build .config for service registry and ship information with all kitex instances.
* Look into cluster designs
