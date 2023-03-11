package settings

type ServiceSettings struct {
	Socket struct {
		Port uint `yaml:"port"`
	} `yaml:"socket"`
	Hypervisor struct {
		Services struct {
			ScriptsDir   string `yaml:"scripts-dir"`
			PortsStorage struct {
				DbPath string `yaml:"db-path"`
			} `yaml:"ports-storage"`
		} `yaml:"services"`
	} `yaml:"hypervisor"`
}
