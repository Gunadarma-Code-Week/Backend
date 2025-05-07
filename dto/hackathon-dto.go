package dto

type RequestHackathon struct {
	LinkDrive string `json:"link_drive"`
}

type HackatonStageStatus struct {
	Stage1    bool   `json:"stage_1"`
	Stage1Url string `json:"stage_1_url"`
	Stage2    bool   `json:"stage_2"`
	Stage2Url string `json:"stage_2_url"`
	Final     bool   `json:"final"`
	FinalUrl  string `json:"final_url"`
}
