#### Some common commands for reference

### Hertz

Execute is not under GOPATH

option 1,you do not have a go.mod, add go mod name after "-module"

hz new -module example.com/m -idl idl/hello.thrift

Tidy & get dependencies
go mod tidy

option 2,you already have a go.mod
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0

Tidy & get dependencies
go mod tidy

https://www.cloudwego.io/docs/hertz/tutorials/toolkit/usage/usage-thrift/

### Kitex

Code generation :
kitex -service a.b.c hello.thrift

# If the current directory is not under $GOPATH/src, you need to add the -module parameter which usually is same as the module name in go.mod

kitex -module "your_module_name" -service a.b.c hello.thrift
