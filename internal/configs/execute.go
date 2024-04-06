package configs

type Execute struct {
	Image           string   `yaml:"image"`
	CommandTemplate []string `yaml:"command_template"`
	CPUQuota        int      `yaml:"cpu_quota"`
}
