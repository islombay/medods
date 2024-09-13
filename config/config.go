package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Server ConfigServer
	DB     ConfigDB
	Redis  ConfigRedis
	Mail   ConfigMail
}

type ConfigDB struct {
	Host           string
	Port           string
	DBName         string
	SSLMode        string
	MigrationsPath string
}

type ConfigMail struct {
	Host   string
	Port   string
	Sender string
}

type ConfigServer struct {
	Host string
	Port string
	URL  string
}

type ConfigRedis struct {
	Host string
	Port string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file found")
	}

	cfg := Config{
		Server: ConfigServer{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
			URL:  os.Getenv("SERVER_URL"),
		},
		DB: ConfigDB{
			Host:           os.Getenv("DB_HOST"),
			Port:           os.Getenv("DB_PORT"),
			DBName:         os.Getenv("DB_NAME"),
			SSLMode:        os.Getenv("DB_SSL_MODE"),
			MigrationsPath: os.Getenv("DB_MIGRATIONS_PATH"),
		},
		Redis: ConfigRedis{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
		Mail: ConfigMail{
			Host:   os.Getenv("MAIL_HOST"),
			Port:   os.Getenv("MAIL_PORT"),
			Sender: os.Getenv("MAIL_SENDER"),
		},
	}

	return cfg
}
