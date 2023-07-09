package ports_service

import (
	"simple-hosting/compositor/app/settings"
	ports_storage "simple-hosting/compositor/app/tools/ports-storage"
)

var config = settings.GetServiceSettings()

func CreatePortRedirect(internalPort int) (PortRedirect, error) {
	storage := ports_storage.Init(config)
	port, err := storage.GetRandomFreePort()
	return PortRedirect{InternalPort: internalPort, ExternalPort: port}, err
}

func ClosePortRedirect(redirect PortRedirect) error {
	port := redirect.ExternalPort
	storage := ports_storage.Init(config)
	return storage.ReleasePort(port)
}
