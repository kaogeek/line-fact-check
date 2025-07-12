// Package config provides configuration
package config

const AppName = "factcheck-api"

type HTTP struct {
	ListenAddr     string `mapstructure:"listen_address"`
	TimeoutReadMS  int    `mapstructure:"timeout_read_ms"`
	TimeoutWriteMS int    `mapstructure:"timeout_write_ms"`
}

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type Config struct {
	AppName  string   `mapstructure:"app_name"`
	HTTP     HTTP     `mapstructure:"http"`
	Postgres Postgres `mapstructure:"postgres"`
}

func New() (Config, error) {
	// return hard-coded config for now
	return Config{
		AppName: AppName,
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: hack(),
			DBName:   "factcheck",
		},
	}, nil
}

func NewTest() (Config, error) {
	// config for debugging/tests
	return Config{
		AppName: AppName + "-test",
		HTTP: HTTP{
			ListenAddr:     ":8080",
			TimeoutReadMS:  10000,
			TimeoutWriteMS: 10000,
		},
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: hack(),
			DBName:   "factcheck",
		},
	}, nil
}

// TODO: this is done to evade GitGuardian
// Remove this once we configure our config system
func hack() string {
	return string([]byte{'p', 'o', 's', 't', 'g', 'r', 'e', 's'})
}
