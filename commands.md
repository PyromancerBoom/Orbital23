#### Some common commands for reference

### Hertz

### Kitex

Code generation :
kitex -service a.b.c hello.thrift

# If the current directory is not under $GOPATH/src, you need to add the -module parameter which usually is same as the module name in go.mod

kitex -module "your_module_name" -service a.b.c hello.thrift
