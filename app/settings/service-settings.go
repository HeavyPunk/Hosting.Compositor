package settings

type ServiceSettings struct {
	App struct {
		Port          uint   `yaml:"port"`
		Configuration string `yaml:"configuration"`
	} `yaml:"app"`
	Hypervisor struct {
		Services struct {
			ScriptsDir   string `yaml:"scripts-dir"`
			PortsService struct {
				DbPath string `yaml:"db-path"`
			} `yaml:"ports-service"`
			OutboundIP string `yaml:"outbound-ip"`
		} `yaml:"services"`
	} `yaml:"hypervisor"`
}
