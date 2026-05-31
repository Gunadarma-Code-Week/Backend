package dto



type UpdateSystemSettingDTO struct {
	HackathonRegistrationDisabled bool       `json:"hackathon_registration_disabled"`
	CPRegistrationDisabled        bool       `json:"cp_registration_disabled"`
	CTFRegistrationDisabled       bool       `json:"ctf_registration_disabled"`
	HackathonProposalDisabled     bool       `json:"hackathon_proposal_disabled"`
	HackathonVideoDisabled        bool       `json:"hackathon_video_disabled"`
	HackathonFinalDisabled        bool       `json:"hackathon_final_disabled"`
	ProfileUpdateDisabled         bool       `json:"profile_update_disabled"`
	HackathonProposalDeadline     *string `json:"hackathon_proposal_deadline"`
	HackathonVideoDeadline        *string `json:"hackathon_video_deadline"`
	HackathonFinalDeadline        *string `json:"hackathon_final_deadline"`
	HackathonProposalChecklist    *string `json:"hackathon_proposal_checklist"`
	HackathonVideoChecklist       *string `json:"hackathon_video_checklist"`
	HackathonFinalChecklist       *string `json:"hackathon_final_checklist"`
	ProfileUpdateDeadline         *string `json:"profile_update_deadline"`
	SeminarRegistrationDisabled   bool    `json:"seminar_registration_disabled"`
	SeminarRequireVerification    bool    `json:"seminar_require_verification"`
}
