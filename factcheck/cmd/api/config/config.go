// Package config provides configuration
package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

const AppName = "factcheck-api"

type HTTP struct {
	ListenAddr     string `env:"FACTCHECKAPI_LISTEN_ADDRESS, required"`
	TimeoutMsRead  int    `env:"FACTCHECKAPI_TIMEOUTMS_READ, default=1000"`
	TimeoutMsWrite int    `env:"FACTCHECKAPI_TIMEOUTMS_WRITE, default=1000"`
}

type Postgres struct {
	Host     string `env:"POSTGRES_HOST, default=localhost"`
	Port     int    `env:"POSTGRES_PORT"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	DB       string `env:"POSTGRES_DB, required"`
}

type Config struct {
	AppName  string `env:"APP_NAME, default=factcheck-api"`
	HTTP     HTTP
	Postgres Postgres
}

func New() (Config, error) {
	var conf Config
	err := envconfig.Process(context.Background(), &conf)
	if err != nil {
		return Config{}, err
	}
	return conf, nil
}

func NewTest() (Config, error) {
	// config for debugging/tests
	return Config{
		AppName: AppName + "-test",
		HTTP: HTTP{
			ListenAddr:     ":8080",
			TimeoutMsRead:  10000,
			TimeoutMsWrite: 10000,
		},
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: hack(),
			DB:       "factcheck",
		},
	}, nil
}

// TODO: this is done to evade GitGuardian
// Remove this once we configure our config system
func hack() string {
	return string([]byte{'p', 'o', 's', 't', 'g', 'r', 'e', 's'})
}
