package config

type Config struct {
	Mysql mysqlConf `yaml:"mysql"`
}

type mysqlConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Table    string `yaml:"table"`
	DB       string `yaml:"db"`
}
