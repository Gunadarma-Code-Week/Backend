package entity

import "time"

type SystemSetting struct {
	ID                            uint       `gorm:"primaryKey" json:"id"`
	HackathonRegistrationDisabled bool       `gorm:"default:false" json:"hackathon_registration_disabled"`
	CPRegistrationDisabled        bool       `gorm:"default:false" json:"cp_registration_disabled"`
	CTFRegistrationDisabled       bool       `gorm:"default:false" json:"ctf_registration_disabled"`
	HackathonProposalDisabled     bool       `gorm:"default:false" json:"hackathon_proposal_disabled"`
	HackathonVideoDisabled        bool       `gorm:"default:false" json:"hackathon_video_disabled"`
	HackathonFinalDisabled        bool       `gorm:"default:false" json:"hackathon_final_disabled"`
	HackathonProposalDeadline     *time.Time `gorm:"default:null" json:"hackathon_proposal_deadline"`
	HackathonVideoDeadline        *time.Time `gorm:"default:null" json:"hackathon_video_deadline"`
	HackathonFinalDeadline        *time.Time `gorm:"default:null" json:"hackathon_final_deadline"`
}
