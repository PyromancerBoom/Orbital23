package main

import (
	"encoding/json"
	"io/ioutil"
)

type Settings struct {
	DockerUrl      string `json:"dockerUrl"`
	DockerPort     string `json:"dockerPort"`
	Env            string `json:"env"`
	ServiceURL     string `json:"serviceurl"`
	ServerPort     string `json:"serverPort"`
	IsDockerised   bool   `json:"isDockerised"`
	GatewayAddress string `json:"gatewayAddress"`
	ConsulAddress  string `json:"consulAddress"`
}

func LoadSettings(filename string) (Settings, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Settings{}, err
	}
	var c Settings
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Settings{}, err
	}
	return c, nil
}
