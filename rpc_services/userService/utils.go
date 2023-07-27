package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type Configuration struct {
	Apikey               string `json:"apikey"`
	ServiceName          string `json:"servicename"`
	DockerUrl            string `json:"dockerUrl"`
	DockerPort           string `json:"dockerPort"`
	Env                  string `json:"env"`
	ServiceURL           string `json:"serviceurl"`
	ServerPort           string `json:"serverPort"`
	HealthCheckFrequency int    `json:"healthCheckFrequency"`
	IsDockerised         bool   `json:"isDockerised"`
	GatewayAddress       string `json:"gatewayAddress"`
}

func LoadConfiguration(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}
	var c Configuration
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}
	return c, nil
}

func getAddr() *net.TCPAddr {

	port, _ := GetFreePort()

	a := "127.0.0.1:" + strconv.Itoa(port)

	addr, err := net.ResolveTCPAddr("tcp", a)
	if err != nil {
		fmt.Println("Error occured." + err.Error() + "Retrying")
		return getAddr()
	}
	return addr
}

func GetFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

// Make an address from address and port
func MakeAddress(address string, port string) (*net.TCPAddr, error) {
	a := strings.TrimSpace(address) + ":" + strings.TrimSpace(port)

	addr, err := net.ResolveTCPAddr("tcp", a)
	if err != nil {
		return nil, err
	}
	return addr, nil
}
