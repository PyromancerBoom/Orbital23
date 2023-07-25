package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
)

type Configuration struct {
	URL        string `json:"url"`
	Port       string `json:"port"`
	Env        string `json:"env"`
	ServiceURL string `json:"serviceurl"`
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