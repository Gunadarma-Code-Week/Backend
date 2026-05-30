package entity



type SystemSetting struct {
	ID                            uint       `gorm:"primaryKey" json:"id"`
	HackathonRegistrationDisabled bool       `gorm:"default:false" json:"hackathon_registration_disabled"`
	CPRegistrationDisabled        bool       `gorm:"default:false" json:"cp_registration_disabled"`
	CTFRegistrationDisabled       bool       `gorm:"default:false" json:"ctf_registration_disabled"`
	HackathonProposalDisabled     bool       `gorm:"default:false" json:"hackathon_proposal_disabled"`
	HackathonVideoDisabled        bool       `gorm:"default:false" json:"hackathon_video_disabled"`
	HackathonFinalDisabled        bool       `gorm:"default:false" json:"hackathon_final_disabled"`
	ProfileUpdateDisabled         bool       `gorm:"default:false" json:"profile_update_disabled"`
	HackathonProposalDeadline     *string `json:"hackathon_proposal_deadline"`
	HackathonVideoDeadline        *string `json:"hackathon_video_deadline"`
	HackathonFinalDeadline        *string `json:"hackathon_final_deadline"`
	ProfileUpdateDeadline         *string `json:"profile_update_deadline"`
}
