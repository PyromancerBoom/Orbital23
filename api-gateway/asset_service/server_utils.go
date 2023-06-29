package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type ConnectRequest struct {
	APIKEY      string `json:"api-key"`
	ServiceName string `json:"serviceName"`
	Address     string `json:"serverAddress"`
	Port        string `json:"serverPort"`
}

// gateway address example : "http://localhost:4200"
// connectsServer to system and gets the server ID back.
func connectServer(gatewayAddress string, apikey string, serviceName string, address string, port string) (string, error) {
	url := gatewayAddress + "/connect"

	req := &ConnectRequest{APIKEY: apikey, ServiceName: serviceName, Address: address, Port: port}

	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	c := make(map[string]json.RawMessage)

	// unmarschal JSON
	e := json.Unmarshal(body, &c)
	if e != nil {
		return "", err
	}

	return strings.Trim(string(c["serverID"]), "\""), nil
}

type UpdateHealthRequest struct {
	APIKEY   string `json:"api-key"`
	ServerID string `json:"serverID"`
}

func updateHealth(gatewayAddress string, apikey string, serverID string) error {

	url := gatewayAddress + "/health"

	req := &UpdateHealthRequest{APIKEY: apikey, ServerID: serverID}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		println(string(body))
	}

	return nil
}
