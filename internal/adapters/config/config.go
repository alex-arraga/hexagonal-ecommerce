package config

import (
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type (
	Container struct {
		App             *App
		Redis           *Redis
		DB              *DB
		HTTP            *HTTP
		PaymentProvider *PaymentProvider
	}

	App struct {
		Name string
		Env  string
	}

	MercadoPago struct {
		PublicKey   string
		AccessToken string
	}

	PaymentProvider struct {
		MercadoPago MercadoPago
	}

	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	DB struct {
		DSN                string
		MaxOpenConnections string
		MaxIdleConnections string
		MaxLifeTime        string
	}

	HTTP struct {
		Env            string
		Domain         string
		Port           string
		AllowedOrigins []string
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
		DB:       0,
	}

	pp := &PaymentProvider{
		MercadoPago: MercadoPago{
			PublicKey:   getEnv("MERCADO_PAGO_PUBLIC_KEY"),
			AccessToken: getEnv("MERCADO_PAGO_ACCESS_TOKEN"),
		},
	}

	db := &DB{
		DSN:                getEnv("DB_DSN"),
		MaxOpenConnections: getEnv("DB_MAX_OPEN_CONNECTIONS"),
		MaxIdleConnections: getEnv("DB_MAX_IDLE_CONNECTIONS"),
		MaxLifeTime:        getEnv("DB_MAX_LIFETIME"),
	}

	allowedOriginsStr := getEnv("APP_ALLOWED_ORIGINS")
	allowedOriginsOpts := strings.Split(allowedOriginsStr, ",")

	http := &HTTP{
		Domain:         getEnv("APP_DOMAIN"),
		Env:            getEnv("APP_ENV"),
		Port:           getEnv("APP_PORT"),
		AllowedOrigins: allowedOriginsOpts,
	}

	return &Container{
		app,
		redis,
		db,
		http,
		pp,
	}, nil

}
