package config

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type Config struct {
	Postgres Postgres `mapstructure:"postgres"`
}

func New() (Config, error) {
	// return hard-coded config for now
	return Config{
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "factcheck",
		},
	}, nil
}

func NewTest() (Config, error) {
	// config for debugging/tests
	return Config{
		Postgres: Postgres{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			DBName:   "factcheck",
		},
	}, nil
}
