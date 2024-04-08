package configs

type MQ struct {
	Addresses       []string `yaml:"addresses"`
	ClientID        string   `yaml:"client_id"`
	ConsumerGroupID string   `yaml:"consumer_group_id"`
	Topic           string   `yaml:"topic"`
	NumPartitions   int      `yaml:"num_partitions"`
}
