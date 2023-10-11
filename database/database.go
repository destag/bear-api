package database

import (
	"log"

	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS bears (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);
`

type Database struct {
	conn *sqlx.DB
}

func Connect(name string) *Database {
	db, err := sqlx.Connect("sqlite3", name)
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	return &Database{db}
}

func (db *Database) Migrate() {
	db.conn.MustExec(schema)
}
