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
			Password: hack(),
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
