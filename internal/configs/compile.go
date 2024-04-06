package configs

import (
	"time"

	"github.com/dustin/go-humanize"
)

type Compile struct {
	Image               string   `yaml:"image"`
	Timeout             string   `yaml:"timeout"`
	Memory              string   `yaml:"memory"`
	CPUs                float32  `yaml:"cpus"`
	CommandTemplate     []string `yaml:"command_template"`
	SourceFileExtension string   `yaml:"source_file_extension"`
}

func (c *Compile) GetTimeoutInTimeDuration() (time.Duration, error) {
	timeout, err := time.ParseDuration(c.Timeout)
	if err != nil {
		return 0, err
	}

	return timeout, nil
}

func (c *Compile) GetMemoryInBytes() (uint64, error) {
	memoryLimit, err := humanize.ParseBytes(c.Memory)
	if err != nil {
		return 0, err
	}

	return memoryLimit, nil
}
