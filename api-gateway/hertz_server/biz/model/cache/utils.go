package cache

//Note that all these cache must be updated in this specific order, or else there will be bugs, as some cache is dependant on other cache.
func UpdateAllCache() error {
	err := UpdateAdminCache()
	if err != nil {
		return err
	}

	err = UpdateIDLcache()
	if err != nil {
		return err
	}

	err = UpdateGenericClientsCache()
	if err != nil {
		return err
	}

	return nil
}
