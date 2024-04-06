package configs

type Execute struct {
	Image           string   `yaml:"image"`
	CommandTemplate []string `yaml:"command_template"`
	CPUs            float32  `yaml:"cpus"`
}
