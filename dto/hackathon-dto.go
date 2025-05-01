package dto

type RequestHackathon struct {
	LinkDrive string `json:"link_drive"`
}

type HackatonStageStatus struct {
	Stage1 bool `json:"stage_1"`
	Stage2 bool `json:"stage_2"`
	Final  bool `json:"final"`
}
