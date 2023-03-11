package settings

type ServiceSettings struct {
	App struct {
		Port          uint   `yaml:"port"`
		Configuration string `yaml:"configuration"`
	} `yaml:"app"`
	Hypervisor struct {
		Services struct {
			ScriptsDir   string `yaml:"scripts-dir"`
			PortsStorage struct {
				DbPath string `yaml:"db-path"`
			} `yaml:"ports-storage"`
		} `yaml:"services"`
	} `yaml:"hypervisor"`
}
