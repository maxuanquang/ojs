package configs

import (
	"time"

	"github.com/dustin/go-humanize"
)

type Compile struct {
	Image           string   `yaml:"image"`
	TimeoutDuration string   `yaml:"timeout_duration"`
	MemoryLimit     string   `yaml:"memory_limit"`
	CPUQuota        int      `yaml:"cpu_quota"`
	CommandTemplate []string `yaml:"command_template"`
	SourceFileName  string   `yaml:"source_file_name"`
	ProgramFileName string   `yaml:"program_file_name"`
}

func (c *Compile) GetTimeoutDuration() (time.Duration, error) {
	timeout, err := time.ParseDuration(c.TimeoutDuration)
	if err != nil {
		return 0, err
	}

	return timeout, nil
}

func (c *Compile) GetMemoryLimitInBytes() (uint64, error) {
	memoryLimit, err := humanize.ParseBytes(c.MemoryLimit)
	if err != nil {
		return 0, err
	}

	return memoryLimit, nil
}
