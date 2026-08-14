package entity

type SystemSetting struct {
	ID                            uint    `gorm:"primaryKey" json:"id"`
	HackathonRegistrationDisabled bool    `gorm:"default:false" json:"hackathon_registration_disabled"`
	CPRegistrationDisabled        bool    `gorm:"default:false" json:"cp_registration_disabled"`
	CTFRegistrationDisabled       bool    `gorm:"default:false" json:"ctf_registration_disabled"`
	HackathonProposalDisabled     bool    `gorm:"default:false" json:"hackathon_proposal_disabled"`
	HackathonVideoDisabled        bool    `gorm:"default:false" json:"hackathon_video_disabled"`
	HackathonFinalDisabled        bool    `gorm:"default:false" json:"hackathon_final_disabled"`
	ProfileUpdateDisabled         bool    `gorm:"default:false" json:"profile_update_disabled"`
	HackathonProposalDeadline     *string `json:"hackathon_proposal_deadline"`
	HackathonVideoDeadline        *string `json:"hackathon_video_deadline"`
	HackathonFinalDeadline        *string `json:"hackathon_final_deadline"`
	CPRegistrationDeadline        *string `json:"cp_registration_deadline"`
	CTFRegistrationDeadline       *string `json:"ctf_registration_deadline"`
	HackathonProposalChecklist    *string `gorm:"type:text" json:"hackathon_proposal_checklist"`
	HackathonVideoChecklist       *string `gorm:"type:text" json:"hackathon_video_checklist"`
	HackathonFinalChecklist       *string `gorm:"type:text" json:"hackathon_final_checklist"`
	ProfileUpdateDeadline         *string `json:"profile_update_deadline"`
	SeminarRegistrationDisabled   bool    `gorm:"default:false" json:"seminar_registration_disabled"`
	SeminarRequireVerification    bool    `gorm:"default:false" json:"seminar_require_verification"`
}
