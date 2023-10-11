package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
)

// import (
//     _ "modernc.org/sqlite"
// )

// db, err := otelsqlx.Open("sqlite", "file::memory:?cache=shared",
// 	otelsql.WithAttributes(semconv.DBSystemSqlite))

func Connect(name string) *sqlx.DB {
	db, err := otelsqlx.Connect("sqlite3", name,
		otelsql.WithAttributes(semconv.DBSystemSqlite))
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
