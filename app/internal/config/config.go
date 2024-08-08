package config

import (
	"sync"

	"github.com/Polyrom/houses_api/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
)

// TODO: add env-default
type Config struct {
	Debug  *bool `yaml:"debug"`
	Listen struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"listen"`
	Storage StorageConfig `yaml:"storage"`
}

type StorageConfig struct {
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Database    string `yaml:"database"`
	MaxAttempts int    `yaml:"max_attempts"`
}

var instance *Config
var once sync.Once

func Get(l logging.Logger) *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			l.Error(help)
			l.Fatal(err)
		}
	})
	return instance
}
