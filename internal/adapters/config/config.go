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
		DSN                string
		MaxOpenConnections string
		MaxIdleConnections string
		MaxLifeTime        string
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
	}
	return env
}

const envFile string = "../../.env"

func New() (*Container, error) {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load(envFile)
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
		DSN:                getEnv("DB_DSN"),
		MaxOpenConnections: getEnv("DB_MAX_OPEN_CONNECTIONS"),
		MaxIdleConnections: getEnv("DB_MAX_IDLE_CONNECTIONS"),
		MaxLifeTime:        getEnv("DB_MAX_LIFETIME"),
	}

	http := &HTTP{
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
