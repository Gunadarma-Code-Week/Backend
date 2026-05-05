package tests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
	"gcw/service"
)

func TestUTSUB003_CPRegistrationProvisioning(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	lead := entity.User{
		Name:  "Unit Lead",
		Email: fmt.Sprintf("lead-%d@example.com", time.Now().UnixNano()),
		Role:  "user",
	}
	if err := db.Create(&lead).Error; err != nil {
		t.Fatalf("create lead user failed: %v", err)
	}

	repo := repository.GateRegistrationRepository(db)
	regSvc := service.NewRegistrationService(repo, &service.DomJudgeService{DomJudgeUrl: ""}, nil)

	request := dto.RegistrationCPTeamRequest{
		RegistraionTeamRequest: dto.RegistraionTeamRequest{
			TeamName:       "Tim CP Unit",
			Supervisor:     "Dosen Unit",
			SupervisorNIDN: "1234567890",
		},
		RegistrationCPRequest: dto.RegistrationCPRequest{
			JoinCode:        "654321",
			BuktiPembayaran: "-",
		},
	}

	result, err := regSvc.CPTeamRegistration(&request, &lead)
	if err != nil {
		t.Fatalf("cp registration failed: %v", err)
	}

	if result.Team.ID_Team == 0 {
		t.Fatal("expected team id to be generated")
	}
	if len(result.CPTeam.JoinCode) != 6 {
		t.Fatalf("expected generated join code length 6, got %q", result.CPTeam.JoinCode)
	}
	if strings.TrimSpace(result.CPTeam.DomjudgeUsername) == "" || strings.TrimSpace(result.CPTeam.DomjudgePassword) == "" {
		t.Fatal("expected domjudge credentials to be generated")
	}

	var refreshedLead entity.User
	if err := db.First(&refreshedLead, lead.ID).Error; err != nil {
		t.Fatalf("reload lead failed: %v", err)
	}
	if refreshedLead.IDTeam == nil {
		t.Fatal("expected lead user to be attached to created team")
	}
}

func TestUTSUB004_CPCredentialRetrieval(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	joinCode := fmt.Sprintf("CP%06d", time.Now().UnixNano()%1000000)
	team := entity.Team{
		TeamName:       "Tim CP Detail",
		Supervisor:     "Dosen",
		SupervisorNIDN: "1234567890",
		JoinCode:       joinCode,
		Event:          "cp",
		ID_LeadTeam:    1,
	}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team failed: %v", err)
	}

	cpTeam := entity.CPTeam{
		IDTeam:           team.ID_Team,
		Stage:            "Registered",
		Status:           "Registration",
		DomjudgeUsername: "cpuser01",
		DomjudgePassword: "cppass01",
	}
	if err := db.Create(&cpTeam).Error; err != nil {
		t.Fatalf("create cp team failed: %v", err)
	}

	cpSvc := service.NewCpService(db)
	result, err := cpSvc.Get(joinCode)
	if err != nil {
		t.Fatalf("cp detail retrieval failed: %v", err)
	}

	if result.DomjudgeUsername != "cpuser01" || result.DomjudgePassword != "cppass01" {
		t.Fatalf("credential mismatch: %+v", result)
	}
	if result.Stage != "Registered" || result.Status != "Registration" {
		t.Fatalf("stage/status mismatch: %+v", result)
	}
}
