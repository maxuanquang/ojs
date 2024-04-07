package configs

import (
	"time"

	"github.com/dustin/go-humanize"
)

type Execute struct {
	Image           string   `yaml:"image"`
	CPUs            float32  `yaml:"cpus"`
	Memory          string   `yaml:"memory"`
	Timeout         string   `yaml:"timeout"`
	CommandTemplate []string `yaml:"command_template"`
}

func (e Execute) GetTimeoutInTimeDuration() (time.Duration, error) {
	return time.ParseDuration(e.Timeout)
}

func (e Execute) GetMemoryInBytes() (uint64, error) {
	return humanize.ParseBytes(e.Memory)
}
