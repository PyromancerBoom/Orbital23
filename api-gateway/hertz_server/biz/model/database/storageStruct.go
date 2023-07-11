package database

// Structure for storage is defined here

type ClientData struct {
	ApiKey    string    `json:"ApiKey"`
	OwnerName string    `json:"OwnerName"`
	OwnerId   string    `json:"OwnerId"`
	Services  []Service `json:"Services"`
}

type Service struct {
	ServiceId          string             `json:"ServiceId"`
	ServiceName        string             `json:"ServiceName"`
	ExposedMethod      string             `json:"ExposedMethod"`
	Path               string             `json:"Path"`
	IdlContent         string             `json:"IdlContent"`
	Version            string             `json:"version"`
	ServiceDescription string             `json:"ServiceDescription"`
	ServerCount        int                `json:"ServerCount"`
	RegisteredServers  []RegisteredServer `json:"RegisteredServers"`
}

type RegisteredServer struct {
	ServerUrl string `json:"ServerUrl"`
	Port      int    `json:"Port"`
}
