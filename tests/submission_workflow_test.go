package tests

import (
	"fmt"
	"testing"
	"time"

	"gcw/dto"
	"gcw/entity"
	"gcw/service"
)

func TestUTSUB001_CreateAndGetSubmissionWorkflow(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	joinCode := fmt.Sprintf("HC%06d", time.Now().UnixNano()%1000000)
	team := entity.Team{
		TeamName:       "UT Hack Team",
		Supervisor:     "Supervisor",
		SupervisorNIDN: "1234567890",
		JoinCode:       joinCode,
		Event:          "hackathon",
		ID_LeadTeam:    1,
	}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team failed: %v", err)
	}

	hackTeam := entity.HackathonTeam{
		IDTeam: team.ID_Team,
		Stage:  "Registered",
		Status: "Registration",
	}
	if err := db.Create(&hackTeam).Error; err != nil {
		t.Fatalf("create hackathon team failed: %v", err)
	}

	svc := service.NewSubmissionService(db)

	if _, err := svc.Create(joinCode, "stage1", dto.RequestHackathon{LinkDrive: "https://drive.test/proposal"}); err != nil {
		t.Fatalf("stage1 submission failed: %v", err)
	}
	if _, err := svc.Create(joinCode, "stage2", dto.RequestHackathon{LinkDrive: "https://drive.test/pitch"}); err != nil {
		t.Fatalf("stage2 submission failed: %v", err)
	}
	if _, err := svc.Create(joinCode, "final", dto.RequestHackathon{LinkDrive: "https://github.com/example/project"}); err != nil {
		t.Fatalf("final submission failed: %v", err)
	}

	status, err := svc.Get(joinCode)
	if err != nil {
		t.Fatalf("get submission status failed: %v", err)
	}

	if !status.Stage1 || status.Stage1Url != "https://drive.test/proposal" {
		t.Fatalf("stage1 status mismatch: %+v", status)
	}
	if !status.Stage2 || status.Stage2Url != "https://drive.test/pitch" {
		t.Fatalf("stage2 status mismatch: %+v", status)
	}
	if !status.Final || status.FinalUrl != "https://github.com/example/project" {
		t.Fatalf("final status mismatch: %+v", status)
	}
}

func TestUTSUB002_CreateSubmission_InvalidJoinCode(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	svc := service.NewSubmissionService(db)

	_, err := svc.Create("NOTFOUND", "stage1", dto.RequestHackathon{LinkDrive: "https://drive.test/proposal"})
	if err == nil {
		t.Fatal("expected error for unknown join_code, got nil")
	}
}

func TestUTUPL004_SubmitSampleSolutionRecorded(t *testing.T) {
	db, cleanup := setupIsolatedPostgresDB(t)
	defer cleanup()

	joinCode := fmt.Sprintf("HS%06d", time.Now().UnixNano()%1000000)
	team := entity.Team{
		TeamName:       "UT Submit Team",
		Supervisor:     "Supervisor",
		SupervisorNIDN: "1234567890",
		JoinCode:       joinCode,
		Event:          "hackathon",
		ID_LeadTeam:    1,
	}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team failed: %v", err)
	}

	hackTeam := entity.HackathonTeam{
		IDTeam: team.ID_Team,
		Stage:  "Registered",
		Status: "Registration",
	}
	if err := db.Create(&hackTeam).Error; err != nil {
		t.Fatalf("create hackathon team failed: %v", err)
	}

	svc := service.NewSubmissionService(db)
	if _, err := svc.Create(joinCode, "stage1", dto.RequestHackathon{LinkDrive: "https://drive.test/sample-solution"}); err != nil {
		t.Fatalf("submit sample solution failed: %v", err)
	}

	var stored entity.HackathonTeam
	if err := db.Where("id_team = ?", team.ID_Team).First(&stored).Error; err != nil {
		t.Fatalf("reload submission failed: %v", err)
	}
	if stored.ProposalUrl != "https://drive.test/sample-solution" {
		t.Fatalf("expected submission stored in proposal_url, got %q", stored.ProposalUrl)
	}
}
