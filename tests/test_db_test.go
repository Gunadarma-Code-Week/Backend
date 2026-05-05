package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"gcw/entity"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupIsolatedPostgresDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()

	_ = godotenv.Load("../.env")

	dbHost := getEnvOrDefault("DB_HOST", "127.0.0.1")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")

	if dbName == "" || dbUser == "" {
		t.Skip("DB_NAME/DB_USER belum diset; skip integration-style unit test")
	}

	baseDSN := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost,
		dbUser,
		dbPass,
		dbName,
		dbPort,
	)

	baseDB, err := gorm.Open(postgres.Open(baseDSN), &gorm.Config{})
	if err != nil {
		t.Skipf("postgres tidak bisa diakses: %v", err)
	}

	schema := fmt.Sprintf("ut_%d", time.Now().UnixNano())
	if err := baseDB.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schema)).Error; err != nil {
		t.Skipf("gagal create schema isolasi test: %v", err)
	}

	testDSN := fmt.Sprintf("%s search_path=%s", baseDSN, schema)
	db, err := gorm.Open(postgres.Open(testDSN), &gorm.Config{})
	if err != nil {
		_ = baseDB.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE`, schema)).Error
		t.Fatalf("gagal buka koneksi schema test: %v", err)
	}

	if err := db.AutoMigrate(
		&entity.Team{},
		&entity.User{},
		&entity.HackathonTeam{},
		&entity.CPTeam{},
	); err != nil {
		_ = baseDB.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE`, schema)).Error
		t.Fatalf("gagal migrate schema test: %v", err)
	}

	cleanup := func() {
		_ = baseDB.Exec(fmt.Sprintf(`DROP SCHEMA IF EXISTS "%s" CASCADE`, schema)).Error

		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
		if sqlBase, err := baseDB.DB(); err == nil {
			_ = sqlBase.Close()
		}
	}

	return db, cleanup
}

func getEnvOrDefault(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}
