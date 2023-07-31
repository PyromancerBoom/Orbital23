// Package to handle server settings from serverconfig.json file
package settings

import (
	"encoding/json"
	"io/ioutil"
)

// Struct to Unmarshal json in
type Settings struct {
	ServerPort   string `json:"serverPort"`
	MasterApiKey string `json:"maskerKey"`

	ConsulAddress string `json:"consulAddress"`
	TTL           int    `json:"timeToLift"`
	TTD           int    `json:"timeToDie"`

	DbUrl                 string `json:"dbUrl"`
	DbName                string `json:"dbName"`
	DbColletionName       string `json:"dbCollectionName"`
	DbPingInterval        int    `json:"dbPingInterval"`
	DbFailedPingInterval  int    `json:"dbFailedPingInterval"`
	DbMaxFailPingDuration int    `json:"dbMaxPingFailDuration"`
}

var setting Settings

// Intialise settings
// Returns error, if any
func InitialiseSettings(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &setting)
	if err != nil {
		return err
	}
	return nil
}

// Returns reference to a settings object
func GetSettings() Settings {
	return setting
}
