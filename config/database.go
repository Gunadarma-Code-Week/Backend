package config

import (
	"fmt"
	"gcw/entity"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() *gorm.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPass, dbName, dbPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to create a connection to database")
	}

	if os.Getenv("ENVIRONMENT") != "production" {
		if err := db.AutoMigrate(
			&entity.User{},
			&entity.UserRole{},
			&entity.Team{},
			&entity.Seminar{},
			&entity.HackathonTeam{},
			&entity.CPTeam{},
			&entity.CTFTeam{},
			&entity.NewsLetter{},
			&entity.AuditLog{},
		); err != nil {
			fmt.Println("AutoMigrate error (ignored):", err)
		}

		// Temporary sequence synchronization fix
		tables := map[string]string{
			"teams":           "id_team",
			"users":           "id",
			"hackathon_teams": "id_hackathon_team",
			"cp_teams":        "id_cp_team",
			"ctf_teams":       "id_ctf_team",
			"seminars":        "id_seminar",
			"audit_logs":      "id",
			"user_roles":      "id_role",
		}
		for table, idCol := range tables {
			var seqName string
			db.Raw(fmt.Sprintf("SELECT pg_get_serial_sequence('\"%s\"', '%s')", table, idCol)).Scan(&seqName)

			if seqName == "" {
				// Fallback to standard naming convention
				seqName = fmt.Sprintf("\"%s_%s_seq\"", table, idCol)
			}

			query := fmt.Sprintf("SELECT setval('%s', COALESCE(MAX(\"%s\"), 0) + 1, false) FROM \"%s\"", seqName, idCol, table)
			if err := db.Exec(query).Error; err != nil {
				fmt.Printf("Failed to sync sequence for %s (seq: %s): %v\n", table, seqName, err)
			} else {
				fmt.Printf("Successfully synced sequence for %s (seq: %s)\n", table, seqName)
			}
		}
	}
	return db
}

func CloseDatabaseConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from database")
	}
	dbSQL.Close()
}
