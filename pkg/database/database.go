package database

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Init() {
	var err error
	dsn := fmt.Sprintf(`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	DB, err = sqlx.Connect("pgx", dsn)
	if err != nil {
		panic(fmt.Sprintf(`Error while establishing DB Connection: %s`, err.Error()))
	}
	DB.SetMaxOpenConns(5)
	DB.SetMaxIdleConns(2)
}
