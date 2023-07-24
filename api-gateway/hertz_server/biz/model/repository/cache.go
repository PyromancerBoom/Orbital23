package repository

var (
	AdminsCache []AdminConfig
)

func UpdateAdminCache() error {
	a, err := GetAllAdmins()
	if err != nil {
		return err
	}

	AdminsCache = a
	return nil
}
