package module

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DBconn ...
func DBconn(dbURL string) *sqlx.DB {
	db, err := sqlx.Connect("mysql", dbURL)
	if err != nil {
		log.Println("sql.Open() Error ...")
	}
	if err = db.Ping(); err != nil {
		log.Println("db.Ping() Error ...")
	}
	db.SetMaxOpenConns(60)
	return db
}

// DBClose ...
func DBClose(db *sqlx.DB) {
	if db != nil {
		db.Close()
	}
}
