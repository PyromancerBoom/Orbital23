package repository

// Structure for storage is defined here
// While the format will be same, the fields would be in smallcase in MongoDB.
// So for example, ApiKey in the dataabase would be stored as "apikey"
// "OwnerName" would be stored as "ownername" by default
// and so on

type AdminConfig struct {
	ApiKey    string    `json:"ApiKey"`
	OwnerName string    `json:"OwnerName"`
	OwnerId   string    `json:"OwnerId"`
	Services  []Service `json:"Services"`
}

type Service struct {
	ServiceId          string             `json:"ServiceId"`
	ServiceName        string             `json:"ServiceName"`
	IdlContent         string             `json:"IdlContent"`
	Version            string             `json:"version"`
	ServiceDescription string             `json:"ServiceDescription"`
	ServerCount        int                `json:"ServerCount"`
	Paths              []Path             `json:"Paths"`
	RegisteredServers  []RegisteredServer `json:"RegisteredServers"`
}

type Path struct {
	MethodPath    string `json:"MethodPath"`
	ExposedMethod string `json:"ExposedMethod"`
}

type RegisteredServer struct {
	ServerUrl string `json:"ServerUrl"`
	Port      int    `json:"Port"`
}
