package config

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/yaml.v2"
)

// TODO -> connect using different drivers

type mySql struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Database string `yaml:"database"`
}

func (cfg *mySql) dataSourceName() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
}

func newConfig(fileName string) mySql {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	var c mySql
	if err = yaml.Unmarshal(data, &c); err != nil {
		fmt.Println(err)
	}

	return c
}

func NewDbFromYaml(fileName string) *sqlx.DB {
	config := newConfig(fileName)

	db, err := sqlx.Open("mysql", config.dataSourceName())
	if err != nil {
		panic(err)
	}

	return db
}
