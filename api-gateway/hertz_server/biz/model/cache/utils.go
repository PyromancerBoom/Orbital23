// Package used for caching Admin details, IDL files, Kitex Clients for improved performance.
package cache

// UpdateCache updates the Admins Cache, IDL Cache and Kitex Generic Clients cache.
// @Returns:
// - error: An error if any
func UpdateCache() error {

	err := updateAdminCache()
	if err != nil {
		return err
	}

	err = updateIDLcache()
	if err != nil {
		return err
	}

	err = updateKitexClientsCache()
	if err != nil {
		return err
	}

	err = updatePathMaskCache()
	if err != nil {
		return err
	}

	return nil
}
