package mysql

import (
	"fmt"
	"github.com/ferch5003/go-fiber-tutorial/config"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(config *config.EnvVars) (*sqlx.DB, error) {
	db := sqlx.MustConnect("mysql", config.MySQLDSN)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("DB Connected!")

	return db, nil
}
