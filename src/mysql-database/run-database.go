package runMysqlDatabase

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func RunDatabase() {
	// Capture connection properties.
	cfg := mysql.Config{
		User:   "arian",
		Passwd: "123",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "resort",
	}

	// Get a database handle.
	var err error
	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	log.Println("Connected!")
}
