package database

import (
	"github.com/jmoiron/sqlx"
)

var schema = `
CREATE TABLE IF NOT EXISTS bears (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);
`

func Migrate(db *sqlx.DB) {

	db.MustExec(schema)

	// or, you can use MustExec, which panics on error
	// cityState := `INSERT INTO place (country, telcode) VALUES (?, ?)`
	// countryCity := `INSERT INTO place (country, city, telcode) VALUES (?, ?, ?)`
	// db.MustExec(cityState, "Hong Kong", 852)
	// db.MustExec(cityState, "Singapore", 65)
	// db.MustExec(countryCity, "South Africa", "Johannesburg", 27)
}
