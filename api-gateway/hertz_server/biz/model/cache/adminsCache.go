package cache

import (
	"api-gateway/hertz_server/biz/model/repository"

	"go.uber.org/zap"
)

var (
	// adminsCache is used to store admin data in memory.
	adminsCache []repository.AdminConfig
)

// updateAdminCache updates the admins cache in an array.
// @Returns:
// - error: An error if any
func updateAdminCache() error {
	a, err := repository.GetAllAdmins()
	if err != nil {
		return err
	}

	adminsCache = a
	zap.L().Debug("Cached Admins.")
	return nil
}
