package internalsql

import (
	"database/sql"
	"log"
)

func Connect(datasource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal("error connecting database", err)
	}

	return db, nil
}
