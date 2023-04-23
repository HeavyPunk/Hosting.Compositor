package settings

import file_settings_provider "simple-hosting/go-commons/settings/file-settings-provider"

func GetServiceSettings() ServiceSettings {
	var config = file_settings_provider.GetSetting[ServiceSettings]("settings.yml")
	return config
}
