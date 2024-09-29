package config

type Config struct {
	Jwt   JwtConf   `yaml:"jwt"`
	Mysql mysqlConf `yaml:"mysql"`
	Redis redisConf `yaml:"redis"`
}

type JwtConf struct {
	Issuer       string `yaml:"issuer"`
	LongExpires  int64  `yaml:"long_expires"`
	PrivateKey   string `yaml:"private_key"`
	PublicKey    string `yaml:"public_key"`
	ShortExpires int64  `yaml:"short_expires"`
}

type mysqlConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Table    string `yaml:"table"`
	DB       string `yaml:"db"`
}

type redisConf struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
