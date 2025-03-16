package dto

type RegistrationResponseDTO struct {
	Name           string `json:"name"`
	Supervisor     string `json:"supervisor"`
	SupervisorNIDN string `json:"supervisor_nidn"`
	ID_LeadTeam    uint64 `json:"id_lead"`
}

type RegistrationResponseWithJoinCode struct {
	Name           string `json:"name"`
	Supervisor     string `json:"supervisor"`
	SupervisorNIDN string `json:"supervisor_nidn"`
	ID_LeadTeam    uint64 `json:"id_lead"`
	JoinCode       string `json:"join_code,omitempty"`
}

type RegistrationResponseHackathon struct {
	Stage            string `json:"stage"`
	Status           string `json:"status"`
	ProposalUrl      string `json:"proposal_url"`
	GithubProjectUrl string `json:"github_project_url"`
	PitchDeckUrl     string `json:"pitch_deck_url"`
	IDTeam           uint64 `json:"id_team,omitempty"`
}

type RegistrationRequestHackathon struct {
	RegistrationResponseWithJoinCode
	RegistrationResponseHackathon
}

type RegistrationCombinedResponse struct {
	Registration  RegistrationResponseWithJoinCode `json:"registration"`
	HackathonTeam RegistrationResponseHackathon    `json:"hackathon_team"`
}
