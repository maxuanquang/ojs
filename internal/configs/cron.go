package configs

type Account struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type CreateSystemAccounts struct {
	Schedule string  `yaml:"schedule"`
	Admin    Account `yaml:"admin"`
	Worker   Account `yaml:"worker"`
}

type Cron struct {
	CreateSystemAccounts CreateSystemAccounts `yaml:"create_system_accounts"`
}
