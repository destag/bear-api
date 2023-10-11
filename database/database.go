package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func Connect(name string) *sqlx.DB {
	db, err := sqlx.Connect("sqlite3", name)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
