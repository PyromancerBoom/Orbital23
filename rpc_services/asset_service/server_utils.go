package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type UpdateHealthRequest struct {
	ApiKey   string `json:"ApiKey"`
	ServerID string `json:"ServerID"`
}

type GatewayClient struct {
	ApiKey         string
	ServiceName    string
	GatewayAddress string
}

type ConnectRequest struct {
	ApiKey      string `json:"ApiKey"`
	ServiceName string `json:"ServiceName"`
	Address     string `json:"ServerAddress"`
	Port        string `json:"ServerPort"`
}

func NewGatewayClient(apikey string, serviceName string, gatewayAddress string) *GatewayClient {
	return &GatewayClient{
		ApiKey:         apikey,
		ServiceName:    serviceName,
		GatewayAddress: gatewayAddress,
	}
}

// Connect server to gateway
func (client *GatewayClient) ConnectServer(serverAddress string, serverPort string) (string, error) {
	url := client.GatewayAddress + "/connect"

	req := &ConnectRequest{ApiKey: client.ApiKey, ServiceName: client.ServiceName, Address: serverAddress, Port: serverPort}

	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")

	httpCli := &http.Client{}
	res, err := httpCli.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	j := make(map[string]json.RawMessage)

	// unmarschal JSON
	e := json.Unmarshal(body, &j)
	if e != nil {
		return "", err
	}

	if string(j["ServerID"]) == "" || string(j["Status"]) == "failed" {
		return "", errors.New("Error connecting to gateway. Message from gateway: " + string(j["Message"]))
	}

	return strings.Trim(string(j["ServerID"]), "\""), nil
}

// Declares that server instance is healthy
func (client *GatewayClient) updateHealth(serverID string) error {

	url := client.GatewayAddress + "/health"

	req := &UpdateHealthRequest{ApiKey: client.ApiKey, ServerID: serverID}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	r.Header.Add("Content-Type", "application/json")

	httpCli := &http.Client{}
	res, err := httpCli.Do(r)
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

// Keeps declaring server instance is healthy
func (client *GatewayClient) UpdateHealthLoop(id string, timeBetweenLoops int) {
	ticker := time.NewTicker(time.Duration(timeBetweenLoops) * time.Second)
	for {
		select {
		case <-ticker.C:
			err := client.updateHealth(id)
			if err != nil {
				// Log the error and continue with the health check loop
				log.Println("Error updating health:", err)
			}
		}
	}
}
