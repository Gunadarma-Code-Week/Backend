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

	// Always run AutoMigrate and Seed in all environments to ensure schema is updated
	if true {
		models := []interface{}{
			&entity.User{},
			&entity.UserRole{},
			&entity.Team{},
			&entity.Seminar{},
			&entity.HackathonTeam{},
			&entity.CPTeam{},
			&entity.CTFTeam{},
			&entity.NewsLetter{},
			&entity.AuditLog{},
			&entity.SystemSetting{},
			&entity.Timeline{},
		}
		for _, model := range models {
			if err := db.AutoMigrate(model); err != nil {
				fmt.Printf("AutoMigrate error for %T (ignored): %v\n", model, err)
			}
		}

		// Seed initial global settings row (ID: 1) if settings table is empty
		var count int64
		db.Model(&entity.SystemSetting{}).Count(&count)
		if count == 0 {
			proposalDeadline := os.Getenv("HACKATHON_PROPOSAL_DEADLINE")
			if proposalDeadline == "" {
				proposalDeadline = "2026-05-24T23:59:59"
			}
			videoDeadline := os.Getenv("HACKATHON_VIDEO_DEADLINE")
			if videoDeadline == "" {
				videoDeadline = "2026-06-06T23:59:59"
			}
			finalDeadline := os.Getenv("HACKATHON_FINAL_DEADLINE")
			if finalDeadline == "" {
				finalDeadline = "2026-06-20T23:59:59"
			}
			
			initialSettings := entity.SystemSetting{
				ID:                            1,
				HackathonRegistrationDisabled: false,
				CPRegistrationDisabled:        false,
				CTFRegistrationDisabled:       false,
				HackathonProposalDisabled:     false,
				HackathonVideoDisabled:        false,
				HackathonFinalDisabled:        false,
				HackathonProposalDeadline:     &proposalDeadline,
				HackathonVideoDeadline:        &videoDeadline,
				HackathonFinalDeadline:        &finalDeadline,
				SeminarRegistrationDisabled:   false,
				SeminarRequireVerification:    false,
			}
			if err := db.Create(&initialSettings).Error; err != nil {
				fmt.Println("Failed to seed initial system settings:", err)
			} else {
				fmt.Println("Successfully seeded initial global system settings record")
			}
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
