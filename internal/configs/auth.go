package configs

import "time"

type Hash struct {
	Cost int `yaml:"cost"`
}

type Token struct {
	Duration            string `yaml:"duration"`
	RS512KeyPairBitSize uint16 `yaml:"rs512_key_pair_bit_size"`
}

func (t *Token) GetTokenDuration() time.Duration {
	duration, _ := time.ParseDuration(t.Duration)
	return duration
}

type Auth struct {
	Hash  Hash  `yaml:"hash"`
	Token Token `yaml:"token"`
}
