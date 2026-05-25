package dto

import "time"

type UpdateSystemSettingDTO struct {
	HackathonRegistrationDisabled bool       `json:"hackathon_registration_disabled"`
	CPRegistrationDisabled        bool       `json:"cp_registration_disabled"`
	CTFRegistrationDisabled       bool       `json:"ctf_registration_disabled"`
	HackathonProposalDisabled     bool       `json:"hackathon_proposal_disabled"`
	HackathonVideoDisabled        bool       `json:"hackathon_video_disabled"`
	HackathonFinalDisabled        bool       `json:"hackathon_final_disabled"`
	HackathonProposalDeadline     *time.Time `json:"hackathon_proposal_deadline"`
	HackathonVideoDeadline        *time.Time `json:"hackathon_video_deadline"`
	HackathonFinalDeadline        *time.Time `json:"hackathon_final_deadline"`
}
