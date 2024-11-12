package global

type ConfigDb struct {
	HostPostgres     string `yaml:"POSTGRES_HOST"`
	UserPostgres     string `yaml:"POSTGRES_USER"`
	PasswordPostgres string `yaml:"POSTGRES_PASSWORD"`
	Port             string `yaml:"PORT"`
}

type Registration struct {
	Login    string
	Password string
}
