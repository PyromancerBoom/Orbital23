package idlmap

import "fmt"

// type IdlMap struct{}

// Represents the relationship between services, method, and IDLs
type IdlMapping struct {
	Service string
	Method  string
	IDL     string
}

// func NewIdlMap() *IdlMap {
// 	fmt.Println("Instantiated")
// 	return &IdlMap{}
// }

var IdlHashMap = []IdlMapping{
	{Service: "AssetManagement", Method: "queryAsset", IDL: "../idl/asset_management.thrift"},
	{Service: "AssetManagement", Method: "insertAsset", IDL: "../idl/asset_management.thrift"},
	{Service: "UserService", Method: "queryUser", IDL: "../idl/user_service.thrift"},
	//{Service: "UserService", Method: "insertUser", IDL: "../../../../idl/user_service.thrift"},
	{Service: "UserService", Method: "insertUser", IDL: "../idl/user_service.thrift"},
	// Can add more mappings similarly using service registry
}

func GetIdlFile(service, method string) (string, error) {
	for _, value := range IdlHashMap {
		if value.Service == service && value.Method == method {
			return value.IDL, nil
		}
	}

	return "", fmt.Errorf("404 : IDL not found\n")
}
