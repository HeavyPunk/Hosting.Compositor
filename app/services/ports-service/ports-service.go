package ports_service

import (
	"errors"
	ports_storage "simple-hosting/compositor/app/tools/ports-storage"
)

func occupyPort(port int) error {
	if ports_storage.IsFreePort(port) {
		ports_storage.OccupyPort(port)
	} else {
		return errors.New("Cannot occupy port " + string(port) + ". Port is busy")
	}
	return nil
}

func releasePort(port int) error {
	if !ports_storage.IsFreePort(port) {
		ports_storage.ReleasePort(port)
	} else {
		return errors.New("Cannot release port " + string(port) + ". Port is free")
	}
	return nil
}

func CreatePortRedirect(internalPort int) (PortRedirect, error) {
	port, err := ports_storage.GetRandomFreePort()
	return PortRedirect{InternalPort: internalPort, ExternalPort: port}, err
}

func ClosePortRedirect(redirect PortRedirect) error {
	port := redirect.ExternalPort
	return ports_storage.ReleasePort(port)
}
