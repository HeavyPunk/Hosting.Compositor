package settings

type ServiceSettings struct {
	Socket struct {
		Port uint `yaml:"port"`
	} `yaml:"socket"`
	Hypervisor struct {
	} `yaml:"hypervisor"`
}
