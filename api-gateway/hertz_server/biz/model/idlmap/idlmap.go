package idlmap

import "fmt"

// Represents the relationship between services, method, and IDLs
type IdlMapping struct {
	Service string
	Method  string
	Path    string
	IDL     string
}

// TODO: Change to hashmap, its currently an array
var IdlHashMap = []IdlMapping{
	{Service: "AssetManagement", Path: "queryAsset", Method: "queryAsset", IDL: "../idl/asset_management.thrift"},
	{Service: "AssetManagement", Path: "insertAsset", Method: "insertAsset", IDL: "../idl/asset_management.thrift"},
	{Service: "UserService", Path: "queryUser", Method: "queryUser", IDL: "../idl/user_service.thrift"},
	//{Service: "UserService", Method: "insertUser", IDL: "../../../../idl/user_service.thrift"},
	{Service: "UserService", Path: "insertUser", Method: "insertUser", IDL: "../idl/user_service.thrift"},
	{Service: "UserService", Path: "insertUserNew", Method: "insertUser", IDL: "../idl/user_service.thrift"},
	// Can add more mappings similarly using service registry
}

func GetIdlFile(service, path string) (IdlMapping, error) {
	for _, value := range IdlHashMap {
		if value.Service == service && value.Path == path {
			return value, nil
		}
	}

	return IdlMapping{}, fmt.Errorf("404 : IDL not found\n")
}
