package internalsql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Connect(datasource string) (*sql.DB, error) {
	if cfg, err := mysqlDriver.ParseDSN(datasource); err == nil {
		if !cfg.MultiStatements {
			cfg.MultiStatements = true
			datasource = cfg.FormatDSN()
		}
	}

	db, err := sql.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(15 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verifikasi koneksi ke database saat startup
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}

func RunMigration(db *sql.DB, migrationPath string) {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatal("Gagal membuat driver migration: ", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal("Gagal inisialisasi migration: ", err)
	}

	version, dirty, err := m.Version()
	if err == nil && dirty {
		log.Printf("Detected dirty migration state at version %d. Auto-resolving and forcing rollback to version %d...", version, version-1)
		if forceErr := m.Force(int(version - 1)); forceErr != nil {
			log.Printf("Warning: failed to force clean migration version: %v", forceErr)
		}
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Migration: tidak ada perubahan, semua sudah up-to-date")
			return
		}
		log.Fatal("Gagal menjalankan migration: ", err)
	}

	log.Println("Migration berhasil dijalankan")
}
