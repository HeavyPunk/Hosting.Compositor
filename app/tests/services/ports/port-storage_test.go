package port_storage_test

import (
	"database/sql"
	"fmt"
	"os"
	"simple-hosting/compositor/app/settings"
	ports_storage "simple-hosting/compositor/app/tools/ports-storage"
	file_settings_provider "simple-hosting/go-commons/settings/file-settings-provider"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func getSettings() settings.ServiceSettings {
	config := file_settings_provider.GetSetting[settings.ServiceSettings]("settings.yml")
	return config
}

func prepareDatabase(config settings.ServiceSettings) error {
  dbPath := config.Hypervisor.Services.PortsService.DbPath
  os.Remove(dbPath)
  file, err := os.Create(dbPath)
  if err != nil {
    return err
  }
  file.Close()
  db, err := sql.Open(config.Hypervisor.Services.PortsService.DbDriver, dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`create table ports(
			id integer not null primary key autoincrement,
			port integer not null
		);`)

	return err
}

func disposeDatabase(config settings.ServiceSettings) error {
  err := os.Remove(config.Hypervisor.Services.PortsService.DbPath)
  return err 
}

func TestOccupyPort(t *testing.T) {
	settings := getSettings()
	err := prepareDatabase(settings)
	if err != nil {
		t.Error(err)
	}

  defer func() {
    err := disposeDatabase(settings)
    if err != nil {
      t.Error(err)
    }
  }()

	storage := ports_storage.Init(settings)
	channels := make(map[int]chan error)
	opCount := 5
	for i := 0; i < opCount; i++ {
		channels[i] = make(chan error)
	}

	for i := 0; i < opCount; i++ {
		go func(chanNum int) {
			port, err := storage.GetRandomFreePort()
			if err != nil {
				channels[chanNum] <- err
				return
			}

			if port <= 0 {
				channels[chanNum] <- fmt.Errorf("returned port %d is out of range", port)
				return
			}
			channels[chanNum] <- nil
		}(i)
	}

	for i := 0; i < opCount; i++ {
		if err = <-channels[i]; err != nil {
			t.Errorf("Error from thread %d: %v", i, err)
		}
	}
}
