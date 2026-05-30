package handler

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type SystemSettingHandler struct {
	settingService *service.SystemSettingService
}

func NewSystemSettingHandler(ss *service.SystemSettingService) *SystemSettingHandler {
	return &SystemSettingHandler{settingService: ss}
}

// @Summary Get System Settings
// @Tags Settings
// @Produce json
// @Success 200 {object} helper.Response{data=entity.SystemSetting}
// @Router /settings [get]
func (h *SystemSettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingService.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", settings))
}

// @Summary Update System Settings (Admin)
// @Tags Admin Settings
// @Accept json
// @Produce json
// @Param request body dto.UpdateSystemSettingDTO true "Update System Settings"
// @Success 200 {object} helper.Response{data=entity.SystemSetting}
// @Router /admin/settings [put]
func (h *SystemSettingHandler) UpdateSettings(c *gin.Context) {
	var input dto.UpdateSystemSettingDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	oldSettings, err := h.settingService.GetSettings()
	var changes []string
	if err == nil {
		if oldSettings.HackathonRegistrationDisabled != input.HackathonRegistrationDisabled {
			status := "dibuka"
			if input.HackathonRegistrationDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("pendaftaran hackathon %s", status))
		}
		if oldSettings.CPRegistrationDisabled != input.CPRegistrationDisabled {
			status := "dibuka"
			if input.CPRegistrationDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("pendaftaran CP %s", status))
		}
		if oldSettings.CTFRegistrationDisabled != input.CTFRegistrationDisabled {
			status := "dibuka"
			if input.CTFRegistrationDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("pendaftaran CTF %s", status))
		}
		if oldSettings.HackathonProposalDisabled != input.HackathonProposalDisabled {
			status := "dibuka"
			if input.HackathonProposalDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("proposal %s", status))
		}
		if oldSettings.HackathonVideoDisabled != input.HackathonVideoDisabled {
			status := "dibuka"
			if input.HackathonVideoDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("video %s", status))
		}
		if oldSettings.HackathonFinalDisabled != input.HackathonFinalDisabled {
			status := "dibuka"
			if input.HackathonFinalDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("final %s", status))
		}
		if oldSettings.ProfileUpdateDisabled != input.ProfileUpdateDisabled {
			status := "dibuka"
			if input.ProfileUpdateDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("pembaruan profil %s", status))
		}
		if oldSettings.HackathonProposalDeadline != input.HackathonProposalDeadline {
			changes = append(changes, fmt.Sprintf("deadline proposal diubah menjadi %s", input.HackathonProposalDeadline))
		}
		if oldSettings.HackathonVideoDeadline != input.HackathonVideoDeadline {
			changes = append(changes, fmt.Sprintf("deadline video diubah menjadi %s", input.HackathonVideoDeadline))
		}
		if oldSettings.HackathonFinalDeadline != input.HackathonFinalDeadline {
			changes = append(changes, fmt.Sprintf("deadline final diubah menjadi %s", input.HackathonFinalDeadline))
		}
		if oldSettings.ProfileUpdateDeadline != input.ProfileUpdateDeadline {
			deadlineStr := "kosong"
			if input.ProfileUpdateDeadline != nil {
				deadlineStr = *input.ProfileUpdateDeadline
			}
			changes = append(changes, fmt.Sprintf("deadline pembaruan profil diubah menjadi %s", deadlineStr))
		}
	}

	if len(changes) > 0 {
		c.Set("audit_changes", strings.Join(changes, ", "))
	}

	newSettings := entity.SystemSetting{
		HackathonRegistrationDisabled: input.HackathonRegistrationDisabled,
		CPRegistrationDisabled:        input.CPRegistrationDisabled,
		CTFRegistrationDisabled:       input.CTFRegistrationDisabled,
		HackathonProposalDisabled:     input.HackathonProposalDisabled,
		HackathonVideoDisabled:        input.HackathonVideoDisabled,
		HackathonFinalDisabled:        input.HackathonFinalDisabled,
		ProfileUpdateDisabled:         input.ProfileUpdateDisabled,
		HackathonProposalDeadline:     input.HackathonProposalDeadline,
		HackathonVideoDeadline:        input.HackathonVideoDeadline,
		HackathonFinalDeadline:        input.HackathonFinalDeadline,
		ProfileUpdateDeadline:         input.ProfileUpdateDeadline,
	}

	updatedSettings, err := h.settingService.UpdateSettings(newSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", updatedSettings))
}
