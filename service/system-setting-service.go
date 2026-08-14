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
	err := s.db.First(&setting).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Auto save to DB if it doesn't exist
			if createErr := s.db.Create(&setting).Error; createErr != nil {
				return setting, createErr
			}
			return setting, nil
		}
		return setting, err
	}
	return setting, nil
}

func (s *SystemSettingService) UpdateSettings(newSettings entity.SystemSetting) (entity.SystemSetting, error) {
	var setting entity.SystemSetting
	err := s.db.First(&setting).Error
	if err != nil {
		// If settings don't exist yet, create the first record
		err = s.db.Create(&newSettings).Error
		return newSettings, err
	}
	setting.HackathonRegistrationDisabled = newSettings.HackathonRegistrationDisabled
	setting.CPRegistrationDisabled = newSettings.CPRegistrationDisabled
	setting.CTFRegistrationDisabled = newSettings.CTFRegistrationDisabled
	setting.HackathonProposalDisabled = newSettings.HackathonProposalDisabled
	setting.HackathonVideoDisabled = newSettings.HackathonVideoDisabled
	setting.HackathonFinalDisabled = newSettings.HackathonFinalDisabled
	setting.ProfileUpdateDisabled = newSettings.ProfileUpdateDisabled
	setting.HackathonProposalDeadline = newSettings.HackathonProposalDeadline
	setting.HackathonVideoDeadline = newSettings.HackathonVideoDeadline
	setting.HackathonFinalDeadline = newSettings.HackathonFinalDeadline
	setting.CPRegistrationDeadline = newSettings.CPRegistrationDeadline
	setting.CTFRegistrationDeadline = newSettings.CTFRegistrationDeadline
	setting.ProfileUpdateDeadline = newSettings.ProfileUpdateDeadline
	setting.SeminarRegistrationDisabled = newSettings.SeminarRegistrationDisabled
	setting.SeminarRequireVerification = newSettings.SeminarRequireVerification
	setting.HackathonProposalChecklist = newSettings.HackathonProposalChecklist
	setting.HackathonVideoChecklist = newSettings.HackathonVideoChecklist
	setting.HackathonFinalChecklist = newSettings.HackathonFinalChecklist
	err = s.db.Save(&setting).Error
	return setting, err
}
