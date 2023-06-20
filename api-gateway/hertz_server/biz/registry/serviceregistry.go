package registry

import (
	"github.com/hashicorp/consul/api"
)

type Service struct {
	consulClient *api.Client
}
