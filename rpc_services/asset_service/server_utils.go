package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type GatewayClient struct {
	API_KEY        string
	Service_Name   string
	GatewayAddress string
}

func NewGatewayClient(apikey string, serviceName string, gatewayAddress string) *GatewayClient {
	return &GatewayClient{apikey, serviceName, gatewayAddress}
}

type ConnectRequest struct {
	APIKEY      string `json:"api-key"`
	ServiceName string `json:"serviceName"`
	Address     string `json:"serverAddress"`
	Port        string `json:"serverPort"`
}

func connectServerWithRetry(client *GatewayClient, serverAddress string, serverPort string) (string, error) {
	for {
		log.Println("Connecting to gateway...")
		id, err := client.connectServer(serverAddress, serverPort)
		if err == nil {
			return id, nil
		}

		log.Println("Error connecting to gateway:", err.Error())
		log.Println("Retrying connection in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}

// gateway address example : "http://localhost:4200"
// connectsServer to system and gets the server ID back.
func (client *GatewayClient) connectServer(serverAddress string, serverPort string) (string, error) {
	url := client.GatewayAddress + "/connect"

	req := &ConnectRequest{APIKEY: client.API_KEY, ServiceName: client.Service_Name, Address: serverAddress, Port: serverPort}

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

	return strings.Trim(string(j["serverID"]), "\""), nil
}

type UpdateHealthRequest struct {
	APIKEY   string `json:"api-key"`
	ServerID string `json:"serverID"`
}

// declares that server is healthy
func (client *GatewayClient) updateHealth(serverID string) error {

	url := client.GatewayAddress + "/health"

	req := &UpdateHealthRequest{APIKEY: client.API_KEY, ServerID: serverID}

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

// keeps declaring server is healthy continuously
func (client *GatewayClient) updateHealthLoop(id string, timeBetweenLoops int) {
	go client.helper(client.GatewayAddress, client.API_KEY, id, timeBetweenLoops)
}

func (client *GatewayClient) helper(gateway string, api string, id string, timeBetweenLoops int) error {
	ticker := time.NewTicker(time.Duration(timeBetweenLoops) * time.Second)
	for {
		err := client.updateHealth(id)
		if err != nil {
			log.Fatal(err.Error())
			return err
		}
		<-ticker.C
	}
}
