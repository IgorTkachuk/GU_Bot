package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	Token string      `yaml:"token" env:"BOT_TOKEN" env-description: "Telegram bot token" env-required: "true"`
	Redis RedisConfig `yaml:"redis"`
}

type RedisConfig struct {
	Host     string `yaml:"host" env:"REDIS_HOST" env-description: "redis host" env-required: "true"`
	Port     int32  `yaml:"port" env:"REDIS_PORT" env-description: "redis port" env-required: "true"`
	Password string `yaml:"password" env:"REDIS_PASSWORD" env-description: "redis password" env-required: "true"`
	Db       int    `yaml:"db" env:"REDIS_DB" env-description: "redis database" env-required: "true"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {

		instance = &Config{}

		if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			log.Println(help)
			log.Fatal(err)
		}

	})
	return instance
}
