package cache

import (
	"api-gateway/hertz_server/biz/model/repository"

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
	zap.L().Debug("Cached Admins.")
	return nil
}
