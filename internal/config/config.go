package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	IsDebug bool `yaml:"isDebug" env:"IS_DEBUG" env-default:"false"`
	IsProd  bool `yaml:"isProd" env:"IS_PROD" env-default:"false"`
	JWT     struct {
		Secret   string `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
		TokenTTL int    `yaml:"tokenTTL" env:"JWT_TOKEN_TTL" env-required:"true"`
	} `yaml:"JWT"`

	Listen struct {
		BindIP string `yaml:"bindIP" env:"BIND_IP" env-default:"0.0.0.0"`
		Port   string `yaml:"port" env:"PORT" env-default:"8080"`
	} `yaml:"Listen"`

	AppConfig struct {
		LogLevel  string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"info"`
		AdminUser struct {
			Login    string `yaml:"login" env:"ADMIN_LOGIN" env-required:"true"`
			Password string `yaml:"password" env:"ADMIN_PASS" env-required:"true"`
		} `yaml:"AdminUser"`
	} `yaml:"App"`

	DBConfig struct {
		Host     string `yaml:"host" env:"DB_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"DB_PORT" env-required:"true"`
		Username string `yaml:"username" env:"DB_USERNAME"`
		Password string `yaml:"password" env:"DB_PASS"`
		DBName   string `yaml:"name" env:"DB_NAME" env-required:"true"`
	} `yaml:"DB"`
}

func GetEnvConfig() *Config {
	log.Println("collecting config")

	config := &Config{}
	if err := cleanenv.ReadEnv(config); err != nil {
		configHeaderText := "Note System"
		helpText, _ := cleanenv.GetDescription(config, &configHeaderText)
		log.Println(helpText)
		log.Fatal(err)
	}

	return config
}

func GetYAMLConfig(configPath string) *Config {
	log.Println("collecting config")

	config := &Config{}
	if err := cleanenv.ReadConfig(configPath, config); err != nil {
		configHeaderText := "Note System"
		helpText, _ := cleanenv.GetDescription(config, &configHeaderText)
		log.Println(helpText)
		log.Fatal(err)
	}

	return config
}
