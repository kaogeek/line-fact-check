package postgres

import (
	"database/sql"
	"fmt"

	"github.com/kaogeek/line-fact-check/factcheck/cmd/api/config"
)

func NewConn(c config.Postgres) (*sql.DB, error) {
	conn, err := sql.Open("pgx", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.DBName))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
