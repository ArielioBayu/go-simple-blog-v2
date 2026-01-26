package memberships

import (
	"database/sql"
	"log"
)

type repository struct {
	Db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	rows, err := db.Query("SELECT id, email FROM users")
	if err != nil {
		log.Println("Error query rows", rows)
	}

	for rows.Next() {
		var id int
		var email string
		err := rows.Scan(&id, &email)
		if err != nil {
			log.Println("error rows scan", err)
		}
		log.Printf("id: %d, email : %s\n", id, email)
	}
	defer rows.Close()

	return &repository{
		Db: db,
	}
}
