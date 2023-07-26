package cache

import (
	repository "api-gateway/hertz_server/biz/model/repository"

	"go.uber.org/zap"
)

var (
	adminsCache []repository.AdminConfig
)

func UpdateAdminCache() error {
	a, err := repository.GetAllAdmins()
	if err != nil {
		return err
	}

	adminsCache = a
	zap.L().Debug("Successfully cached Admin files to memory.")
	return nil
}
