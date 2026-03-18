package tests

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gcw/service"
)

func TestUTUPL002_DomjudgeCreateTeamAndUser_RequestContract(t *testing.T) {
	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret"))

	teamCalled := false
	userCalled := false

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/teams", func(w http.ResponseWriter, r *http.Request) {
		teamCalled = true

		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method for teams: %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != expectedAuth {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "invalid credentials"})
			return
		}
		if got := r.URL.Query().Get("cid"); got != "1" {
			t.Fatalf("unexpected cid: %s", got)
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode team payload failed: %v", err)
		}
		if payload["id"] != "165191" {
			t.Fatalf("unexpected team id payload: %v", payload["id"])
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "165191"})
	})
	mux.HandleFunc("/api/v4/users", func(w http.ResponseWriter, r *http.Request) {
		userCalled = true

		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method for users: %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != expectedAuth {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "invalid credentials"})
			return
		}

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode user payload failed: %v", err)
		}
		if payload["team_id"] != "165191" {
			t.Fatalf("unexpected team_id payload: %v", payload["team_id"])
		}

		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "U-1"})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := &service.DomJudgeService{
		DomJudgeUrl:       server.URL,
		DomJudgeContestID: "1",
		DomJudgeAuth:      base64.StdEncoding.EncodeToString([]byte("admin:secret")),
	}

	username, password, err := svc.CreateDomJudgeTeamUser("165191", "Tim CP", "lead@example.com")
	if err != nil {
		t.Fatalf("CreateDomJudgeTeamUser failed: %v", err)
	}

	if !teamCalled || !userCalled {
		t.Fatalf("expected both endpoints called, team=%v user=%v", teamCalled, userCalled)
	}
	if len(username) != 10 || len(password) != 10 {
		t.Fatalf("expected generated credentials len=10, got username=%d password=%d", len(username), len(password))
	}
}

func TestUTUPL003_DomjudgeCreateTeam_InvalidCredential(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/teams", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{"message": "invalid credentials"})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := &service.DomJudgeService{
		DomJudgeUrl:       server.URL,
		DomJudgeContestID: "1",
		DomJudgeAuth:      base64.StdEncoding.EncodeToString([]byte("wrong:wrong")),
	}

	_, err := svc.CreateTeam("165191", "Tim CP")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got: %v", err)
	}
}

func TestUTJDG001_DomjudgeAPIAuthHeader(t *testing.T) {
	expectedAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte("admin:secret"))

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v4/teams", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != expectedAuth {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]any{"message": "invalid credentials"})
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "TEAM-1"})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	svc := &service.DomJudgeService{
		DomJudgeUrl:       server.URL,
		DomJudgeContestID: "1",
		DomJudgeAuth:      base64.StdEncoding.EncodeToString([]byte("admin:secret")),
	}

	id, err := svc.CreateTeam("TEAM-1", "Team One")
	if err != nil {
		t.Fatalf("expected auth header accepted, got error: %v", err)
	}
	if id != "TEAM-1" {
		t.Fatalf("expected returned team id TEAM-1, got %s", id)
	}
}

func TestUTUPL006_DomjudgeSkipWhenURLUnset(t *testing.T) {
	svc := &service.DomJudgeService{}

	result, err := svc.CreateTeam("165191", "Tim CP")
	if err != nil {
		t.Fatalf("expected nil error when URL unset, got: %v", err)
	}
	if result != "skipped" {
		t.Fatalf("expected skipped result, got: %s", result)
	}
}
