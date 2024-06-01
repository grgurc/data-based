package config

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TODO -> connect using different drivers,
// give the user ability to select them

type dbConfig struct {
	user     string
	password string
	host     string
	port     string
	database string
}

func (cfg *dbConfig) dataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		cfg.user,
		cfg.password,
		cfg.host,
		cfg.port,
		cfg.database,
	)
}

func NewConfig() dbConfig {
	return dbConfig{
		user:     "root",
		password: "test",
		host:     "localhost",
		port:     "3306",
		database: "b2match",
	}
}

func NewDB(config dbConfig) *sqlx.DB {
	db, err := sqlx.Open("mysql", config.dataSourceName())
	if err != nil {
		panic(err)
	}

	return db
}
