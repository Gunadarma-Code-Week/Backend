package dto

type UpdateSystemSettingDTO struct {
	HackathonRegistrationDisabled bool   `json:"hackathon_registration_disabled"`
	CPRegistrationDisabled        bool   `json:"cp_registration_disabled"`
	CTFRegistrationDisabled       bool   `json:"ctf_registration_disabled"`
	HackathonProposalDisabled     bool   `json:"hackathon_proposal_disabled"`
	HackathonVideoDisabled        bool   `json:"hackathon_video_disabled"`
	HackathonFinalDisabled        bool   `json:"hackathon_final_disabled"`
	HackathonProposalDeadline     string `json:"hackathon_proposal_deadline" binding:"max=50"`
	HackathonVideoDeadline        string `json:"hackathon_video_deadline" binding:"max=50"`
	HackathonFinalDeadline        string `json:"hackathon_final_deadline" binding:"max=50"`
}
