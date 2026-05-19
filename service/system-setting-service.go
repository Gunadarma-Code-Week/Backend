package service

import (
	"gcw/entity"
	"gorm.io/gorm"
)

type SystemSettingService struct {
	db *gorm.DB
}

func NewSystemSettingService(db *gorm.DB) *SystemSettingService {
	return &SystemSettingService{db: db}
}

func (s *SystemSettingService) GetSettings() (entity.SystemSetting, error) {
	var setting entity.SystemSetting
	err := s.db.First(&setting, 1).Error
	return setting, err
}

func (s *SystemSettingService) UpdateSettings(newSettings entity.SystemSetting) (entity.SystemSetting, error) {
	var setting entity.SystemSetting
	err := s.db.First(&setting, 1).Error
	if err != nil {
		// If settings don't exist yet, create the first record
		newSettings.ID = 1
		err = s.db.Create(&newSettings).Error
		return newSettings, err
	}
	setting.HackathonRegistrationDisabled = newSettings.HackathonRegistrationDisabled
	setting.CPRegistrationDisabled = newSettings.CPRegistrationDisabled
	setting.CTFRegistrationDisabled = newSettings.CTFRegistrationDisabled
	setting.HackathonProposalDisabled = newSettings.HackathonProposalDisabled
	setting.HackathonVideoDisabled = newSettings.HackathonVideoDisabled
	setting.HackathonFinalDisabled = newSettings.HackathonFinalDisabled
	setting.HackathonProposalDeadline = newSettings.HackathonProposalDeadline
	setting.HackathonVideoDeadline = newSettings.HackathonVideoDeadline
	setting.HackathonFinalDeadline = newSettings.HackathonFinalDeadline
	err = s.db.Save(&setting).Error
	return setting, err
}
