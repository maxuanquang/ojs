package configs

type SystemAccount struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type SystemAccounts struct {
	Admin  SystemAccount `yaml:"admin"`
	Worker SystemAccount `yaml:"worker"`
}
