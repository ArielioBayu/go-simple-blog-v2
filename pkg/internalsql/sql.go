package internalsql

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(datasource string) (*sql.DB, error) {
	db, err := sql.Open("mysql", datasource)
	if err != nil {
		log.Fatal("error connecting database", err)
	}

	return db, nil
}

// RunMigration menjalankan semua file migration yang belum diaplikasikan
// dari folder yang ditentukan oleh migrationPath (contoh: "scripts/migrations").
func RunMigration(db *sql.DB, migrationPath string) {
	// Buat driver migrate dari koneksi *sql.DB yang sudah ada
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("Gagal membuat driver migration: ", err)
	}

	// Inisialisasi migrate dengan source file dan driver database
	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Gagal inisialisasi migration: ", err)
	}

	// Jalankan semua migration yang belum diaplikasikan (Up)
	if err := m.Up(); err != nil {
		// ErrNoChange bukan error nyata — berarti semua migration sudah up-to-date
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Migration: tidak ada perubahan, semua sudah up-to-date")
			return
		}
		log.Fatal("Gagal menjalankan migration: ", err)
	}

	log.Println("Migration berhasil dijalankan")
}
