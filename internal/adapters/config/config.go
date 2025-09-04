package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App   *App
		Redis *Redis
		DB    *DB
		HTTP  *HTTP
	}

	App struct {
		Name string
		Env  string
	}

	Redis struct {
		Addr     string
		Password string
	}

	DB struct {
		Connection string
		Host       string
		Port       string
		User       string
		Password   string
		Name       string
	}

	HTTP struct {
		Env            string
		URL            string
		Port           string
		AllowedOrigins string
	}
)

func getEnv(value string) string {
	env := os.Getenv(value)
	if env == "" {
		slog.Error("couldn't get enviroment variable", "error", value)
		panic("enviroment variable not found")
	}
	return env
}

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
	}

	app := &App{
		Name: getEnv("APP_NAME"),
		Env:  getEnv("APP_ENV"),
	}

	redis := &Redis{
		Addr:     getEnv("REDIS_ADDR"),
		Password: getEnv("REDIS_PASSWORD"),
	}

	db := &DB{
		Connection: getEnv("DB_CONNECTION"),
		Host:       getEnv("DB_HOST"),
		Port:       getEnv("DB_PORT"),
		User:       getEnv("DB_USER"),
		Password:   getEnv("DB_PASSWORD"),
		Name:       getEnv("DB_NAME"),
	}

	http := &HTTP{
		Env:            getEnv("APP_ENV"),
		URL:            getEnv("APP_URL"),
		Port:           getEnv("APP_PORT"),
		AllowedOrigins: getEnv("APP_ALLOWED_ORIGINS"),
	}

	return &Container{
		app,
		redis,
		db,
		http,
	}, nil

}
