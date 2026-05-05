package tests

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"gcw/dto"
	"gcw/entity"
	"gcw/service"
)

func TestUTJDG002_UpdateCpStageWorkflow(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	team := entity.Team{
		TeamName:       "Tim CP Lama",
		Supervisor:     "Dosen",
		SupervisorNIDN: "1234567890",
		JoinCode:       "991122",
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
		DomjudgeUsername: "cpuser-update",
		DomjudgePassword: "cppass-update",
	}
	if err := db.Create(&cpTeam).Error; err != nil {
		t.Fatalf("create cp team failed: %v", err)
	}

	idTeam := team.ID_Team
	member := entity.User{
		Name:   "Member No Email",
		Email:  "",
		Role:   "user",
		IDTeam: &idTeam,
	}
	if err := db.Create(&member).Error; err != nil {
		t.Fatalf("create member failed: %v", err)
	}

	svc := service.DashboardServices{
		DB:           db,
		EmailService: &service.EmailService{},
	}

	err := svc.UpdateCpService(
		strconv.FormatUint(cpTeam.ID_CPTeam, 10),
		dto.Cp{NamaTim: "Tim CP Lolos", Stage: "Final"},
	)
	if err != nil {
		t.Fatalf("update cp stage failed: %v", err)
	}

	var updatedCP entity.CPTeam
	if err := db.First(&updatedCP, cpTeam.ID_CPTeam).Error; err != nil {
		t.Fatalf("reload cp team failed: %v", err)
	}
	if updatedCP.Stage != "Final" {
		t.Fatalf("expected cp stage Final, got %q", updatedCP.Stage)
	}

	var updatedTeam entity.Team
	if err := db.First(&updatedTeam, team.ID_Team).Error; err != nil {
		t.Fatalf("reload team failed: %v", err)
	}
	if updatedTeam.TeamName != "Tim CP Lolos" {
		t.Fatalf("expected team name updated, got %q", updatedTeam.TeamName)
	}
}

func TestUTUPL005_JudgingResultGenerated(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	joinCode := fmt.Sprintf("CP%06d", time.Now().UnixNano()%1000000)
	team := entity.Team{
		TeamName:       "Tim CP Submit",
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
		DomjudgeUsername: "cp-submit",
		DomjudgePassword: "cp-pass",
	}
	if err := db.Create(&cpTeam).Error; err != nil {
		t.Fatalf("create cp team failed: %v", err)
	}

	svc := service.DashboardServices{
		DB:           db,
		EmailService: &service.EmailService{},
	}

	if err := svc.UpdateCpService(strconv.FormatUint(cpTeam.ID_CPTeam, 10), dto.Cp{
		NamaTim: "Tim CP Submit",
		Stage:   "Qualified",
	}); err != nil {
		t.Fatalf("update cp stage failed: %v", err)
	}

	cpSvc := service.NewCpService(db)
	result, err := cpSvc.Get(joinCode)
	if err != nil {
		t.Fatalf("get cp detail failed: %v", err)
	}
	if result.Stage != "Qualified" {
		t.Fatalf("expected stage Qualified after judging, got %q", result.Stage)
	}
}
