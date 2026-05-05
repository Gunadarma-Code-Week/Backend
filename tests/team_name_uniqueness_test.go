package tests

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
	"gcw/service"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestTeamNameUniqueness(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	repo := repository.GateRegistrationRepository(db)
	regSvc := service.NewRegistrationService(repo, &service.DomJudgeService{DomJudgeUrl: ""}, nil)

	// --- 1. Test Duplicate CP Team Name ---
	t.Run("Duplicate CP Team Name", func(t *testing.T) {
		teamName := "Unique CP Team"
		
		lead1 := createTestUser(t, db, "lead1-cp")
		request1 := createCPRequest(teamName)
		_, err := regSvc.CPTeamRegistration(&request1, &lead1)
		if err != nil {
			t.Fatalf("first CP registration failed: %v", err)
		}

		lead2 := createTestUser(t, db, "lead2-cp")
		_, err = regSvc.CPTeamRegistration(&request1, &lead2)
		if err == nil || err.Error() != "Nama Tim Sudah Digunakan" {
			t.Fatalf("expected Nama Tim Sudah Digunakan error, got %v", err)
		}
	})

	// --- 2. Test Duplicate Hackathon Team Name ---
	t.Run("Duplicate Hackathon Team Name", func(t *testing.T) {
		teamName := "Unique Hackathon Team"
		
		lead1 := createTestUser(t, db, "lead1-hack")
		request1 := createHackathonRequest(teamName)
		_, err := regSvc.HackathonTeamRegistration(&request1, &lead1)
		if err != nil {
			t.Fatalf("first Hackathon registration failed: %v", err)
		}

		lead2 := createTestUser(t, db, "lead2-hack")
		_, err = regSvc.HackathonTeamRegistration(&request1, &lead2)
		if err == nil || err.Error() != "Nama Tim Sudah Digunakan" {
			t.Fatalf("expected Nama Tim Sudah Digunakan error, got %v", err)
		}
	})

	// --- 3. Test Same Name Lintas Event (Global Uniqueness) ---
	t.Run("Same Name Different Events (Global Unique)", func(t *testing.T) {
		teamName := "Shared Team Name"

		// Register CP
		leadCP := createTestUser(t, db, "lead-cp-shared")
		requestCP := createCPRequest(teamName)
		_, err := regSvc.CPTeamRegistration(&requestCP, &leadCP)
		if err != nil {
			t.Fatalf("CP registration with shared name failed: %v", err)
		}

		// Register Hackathon dengan nama yang sama — harus DITOLAK (global uniqueness)
		leadHack := createTestUser(t, db, "lead-hack-shared")
		requestHack := createHackathonRequest(teamName)
		_, err = regSvc.HackathonTeamRegistration(&requestHack, &leadHack)
		if err == nil || err.Error() != "Nama Tim Sudah Digunakan" {
			t.Fatalf("expected Nama Tim Sudah Digunakan error when name used across events, got %v", err)
		}
	})


	t.Run("Case Insensitive Duplicate Name", func(t *testing.T) {
		teamName1 := "Ctf Team Name"
		teamName2 := "ctF tEAM nAme"
		
		lead1 := createTestUser(t, db, "lead1-ctf")
		request1 := createCTFRequest(teamName1)
		_, err := regSvc.CTFTeamRegistration(&request1, &lead1)
		if err != nil {
			t.Fatalf("first CTF registration failed: %v", err)
		}

		lead2 := createTestUser(t, db, "lead2-ctf")
		request2 := createCTFRequest(teamName2)
		_, err = regSvc.CTFTeamRegistration(&request2, &lead2)
		if err == nil || err.Error() != "Nama Tim Sudah Digunakan" {
			t.Fatalf("expected Nama Tim Sudah Digunakan error, got %v", err)
		}
	})

	// --- 5. Test Ignored Soft Deleted Team Name ---
	t.Run("Ignored Soft Deleted Team Name", func(t *testing.T) {
		teamName := "Deleted Hackathon Team"
		
		lead1 := createTestUser(t, db, "lead-hack-del1")
		request1 := createHackathonRequest(teamName)
		resp1, err := regSvc.HackathonTeamRegistration(&request1, &lead1)
		if err != nil {
			t.Fatalf("first Hackathon registration failed: %v", err)
		}

		// Soft delete the first team
		err = db.Model(&entity.HackathonTeam{}).Where("id_team = ?", resp1.Team.ID_Team).Update("is_deleted", true).Error
		if err != nil {
			t.Fatalf("failed to soft delete hackathon team: %v", err)
		}

		lead2 := createTestUser(t, db, "lead-hack-del2")
		request2 := createHackathonRequest(teamName)
		_, err = regSvc.HackathonTeamRegistration(&request2, &lead2)
		if err != nil {
			t.Fatalf("second Hackathon registration after soft delete failed: %v", err)
		}
	})
}

func createTestUser(t *testing.T, db *gorm.DB, prefix string) entity.User {
	user := entity.User{
		Name:  prefix,
		Email: fmt.Sprintf("%s-%d@example.com", prefix, time.Now().UnixNano()),
		Role:  "user",
	}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create test user failed: %v", err)
	}
	return user
}

func createCPRequest(name string) dto.RegistrationCPTeamRequest {
	return dto.RegistrationCPTeamRequest{
		RegistraionTeamRequest: dto.RegistraionTeamRequest{
			TeamName:       name,
			Supervisor:     "Supervisor",
			SupervisorNIDN: "1234567890",
		},
		RegistrationCPRequest: dto.RegistrationCPRequest{
			BuktiPembayaran: "-",
		},
	}
}

func createHackathonRequest(name string) dto.RegistrationHackathonTeamRequest {
	return dto.RegistrationHackathonTeamRequest{
		RegistraionTeamRequest: dto.RegistraionTeamRequest{
			TeamName:       name,
			Supervisor:     "Supervisor",
			SupervisorNIDN: "1234567890",
		},
	}
}

func createCTFRequest(name string) dto.RegistrationCTFTeamRequest {
	return dto.RegistrationCTFTeamRequest{
		RegistraionTeamRequest: dto.RegistraionTeamRequest{
			TeamName:       name,
			Supervisor:     "Supervisor",
			SupervisorNIDN: "1234567890",
		},
		RegistrationCTFRequest: dto.RegistrationCTFRequest{
			BuktiPembayaran: "-",
		},
	}
}
