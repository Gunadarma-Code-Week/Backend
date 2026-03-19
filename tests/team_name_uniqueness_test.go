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
		if err == nil || err.Error() != "TEAM NAME ALREADY TAKEN" {
			t.Fatalf("expected TEAM NAME ALREADY TAKEN error, got %v", err)
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
		if err == nil || err.Error() != "TEAM NAME ALREADY TAKEN" {
			t.Fatalf("expected TEAM NAME ALREADY TAKEN error, got %v", err)
		}
	})

	// --- 3. Test Same Name Different Events ---
	t.Run("Same Name Different Events", func(t *testing.T) {
		teamName := "Shared Team Name"
		
		// Register CP
		leadCP := createTestUser(t, db, "lead-cp-shared")
		requestCP := createCPRequest(teamName)
		_, err := regSvc.CPTeamRegistration(&requestCP, &leadCP)
		if err != nil {
			t.Fatalf("CP registration with shared name failed: %v", err)
		}

		// Register Hackathon (should succeed)
		leadHack := createTestUser(t, db, "lead-hack-shared")
		requestHack := createHackathonRequest(teamName)
		_, err = regSvc.HackathonTeamRegistration(&requestHack, &leadHack)
		if err != nil {
			t.Fatalf("Hackathon registration with shared name failed: %v", err)
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
