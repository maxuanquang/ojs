package configs

type Judge struct {
	Languages []Language `yaml:"languages"`
}

type Language struct {
	Value   string   `yaml:"value"`
	Name    string   `yaml:"name"`
	Compile *Compile `yaml:"compile"`
	Execute *Execute `yaml:"execute"`
}
