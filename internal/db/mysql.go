package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var MySQL *sqlx.DB

func SetupMySQL() error {
	var err error
	MySQL, err = sqlx.Open("mysql", viper.GetString("db.mysql.dsn"))
	if err != nil {
		return err
	}

	_ = err
	if idle := viper.GetInt("db.mysql.idle"); idle > 0 {
		MySQL.SetMaxIdleConns(idle)
	}

	if open := viper.GetInt("db.mysql.open"); open > 0 {
		MySQL.SetMaxOpenConns(open)
	}

	return MySQL.Ping()
}
