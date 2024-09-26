package config

type Config struct {
	Jwt   JwtConf   `yaml:"jwt"`
	Mysql mysqlConf `yaml:"mysql"`
}

type JwtConf struct {
	Issuer       string `yaml:"issuer"`
	LongExpires  int    `yaml:"long_expires"`
	Secret       string `yaml:"secret"`
	ShortExpires int    `yaml:"short_expires"`
}

type mysqlConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Table    string `yaml:"table"`
	DB       string `yaml:"db"`
}
