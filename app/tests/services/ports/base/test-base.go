package ports_test_base

import (
	"database/sql"
	"os"
	"simple-hosting/compositor/app/settings"
  file_settings_provider "simple-hosting/go-commons/settings/file-settings-provider"
)


func GetSettings(settingsPath string) settings.ServiceSettings {
	config := file_settings_provider.GetSetting[settings.ServiceSettings](settingsPath)
	return config
}

func PrepareDatabase(config settings.ServiceSettings) error {
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

func DisposeDatabase(config settings.ServiceSettings) error {
  err := os.Remove(config.Hypervisor.Services.PortsService.DbPath)
  return err 
}

