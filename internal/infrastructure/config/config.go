package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server   Server   `env-prefix:"SERVER_"`
	Database Database `env-prefix:"DB_"`
	Logger   Logger   `env-prefix:"LOGGER_"`
	JWT      JWT      `env-prefix:"JWT_"`
}

type Server struct {
	Port int `env:"PORT" env-required:"true"`
}

type JWT struct {
	JWTSecret string `env:"SECRET" env-required:"true"`
}

type Database struct {
	Host     string `env:"HOST" env-required:"true"`
	Port     int    `env:"PORT" env-required:"true"`
	User     string `env:"USER" env-required:"true"`
	Password string `env:"PASSWORD" env-required:"true"`
	Name     string `env:"NAME" env-required:"true"`
}

type Logger struct {
	FileName string `env:"FILE" env-required:"true"`
}

func MustLoad() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Fatal("Error loading .env file")
	}

	var config Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		logrus.Fatal(err)
		return nil, err
	}

	logrus.Info("Loaded app configuration")

	return &config, nil
}
