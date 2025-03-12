package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gcw/helper"
	"io"
	"net/http"
	"os"
)

type DomJudgeService struct {
	DomJudgeUrl       string
	DomJudgeContestID string
	DomJudgeAuth      string
}

func NewDomJudgeService() *DomJudgeService {
	domJudgeUsername := os.Getenv("DOMJUDGE_USERNAME")
	domJudgePassword := os.Getenv("DOMJUDGE_PASSWORD")

	auth := domJudgeUsername + ":" + domJudgePassword
	auth = base64.StdEncoding.EncodeToString([]byte(auth))

	return &DomJudgeService{
		DomJudgeUrl:       os.Getenv("DOMJUDGE_URL"),
		DomJudgeContestID: os.Getenv("DOMJUDGE_CONTEST_ID"),
		DomJudgeAuth:      auth,
	}
}

func (s *DomJudgeService) CreateTeam(
	id string, // use team id
	name string, // use team name
) (string, error) {
	client := &http.Client{}

	body := map[string]any{
		"id":           id,
		"label":        id,
		"name":         name,
		"display_name": name,
		"group_ids": []string{
			"3",
		},
	}
	jsonBody, _ := json.Marshal(body)
	bufferBody := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", s.DomJudgeUrl+"/api/v4/teams?cid="+s.DomJudgeContestID, bufferBody)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Basic "+s.DomJudgeAuth)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", err
	}

	if response.StatusCode != 201 {
		return "", errors.New(data["message"].(string))
	}

	return data["id"].(string), nil
}

func (s *DomJudgeService) CreateUser(
	username string, // use random string
	name string, // use team name
	email string, // use lead email
	password string, // use random string
	teamId string, // use team id result from CreateTeam
) (string, error) {
	client := &http.Client{}

	body := map[string]any{
		"username": username,
		"name":     name,
		"email":    email,
		"password": password,
		"enabled":  true,
		"team_id":  teamId,
		"roles": []string{
			"team",
		},
	}
	jsonBody, _ := json.Marshal(body)
	bufferBody := bytes.NewBuffer(jsonBody)

	req, err := http.NewRequest("POST", s.DomJudgeUrl+"/api/v4/users", bufferBody)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Basic "+s.DomJudgeAuth)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(responseBody, &data); err != nil {
		return "", err
	}

	if response.StatusCode != 201 {
		return "", errors.New(data["message"].(string))
	}

	return data["id"].(string), nil
}

func (s *DomJudgeService) CreateDomJudgeTeamUser(
	teamId string,
	teamName string,
	leadEmail string,
) (username string, password string, err error) {
	username = helper.RandomString(10)
	password = helper.RandomString(10)

	domJudgeTeamId, err := s.CreateTeam(teamId, teamName)
	if err != nil {
		return
	}

	_, err = s.CreateUser(username, teamName, leadEmail, password, domJudgeTeamId)
	if err != nil {
		return
	}

	return
}
