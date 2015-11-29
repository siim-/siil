package entity

import (
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func CreateConnection(url string) error {
	if conn, err := sqlx.Connect("mysql", url); err != nil {
		return err
	} else {
		DB = conn
		return nil
	}
}
