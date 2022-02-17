package config

type Store struct {
	User     string `yaml:"user"`
	Pwd      string `yaml:"pwd"`
	Dsn      string `yaml:"dsn"`
	Port     string `yaml:"port"`
	Database string `yaml:"db"`
}
