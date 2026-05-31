package handler

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/service"
	"net/http"
	"strings"
	"time"

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
		if oldSettings.SeminarRegistrationDisabled != input.SeminarRegistrationDisabled {
			status := "dibuka"
			if input.SeminarRegistrationDisabled {
				status = "ditutup"
			}
			changes = append(changes, fmt.Sprintf("pendaftaran seminar %s", status))
		}
		if oldSettings.SeminarRequireVerification != input.SeminarRequireVerification {
			status := "dinonaktifkan"
			if input.SeminarRequireVerification {
				status = "diaktifkan"
			}
			changes = append(changes, fmt.Sprintf("syarat verifikasi pendaftaran seminar %s", status))
		}
		if isStringPtrChanged(oldSettings.HackathonProposalDeadline, input.HackathonProposalDeadline) {
			changes = append(changes, fmt.Sprintf("deadline proposal diubah dari \"%s\" menjadi \"%s\"",
				formatDeadline(oldSettings.HackathonProposalDeadline),
				formatDeadline(input.HackathonProposalDeadline)))
		}
		if isStringPtrChanged(oldSettings.HackathonVideoDeadline, input.HackathonVideoDeadline) {
			changes = append(changes, fmt.Sprintf("deadline video diubah dari \"%s\" menjadi \"%s\"",
				formatDeadline(oldSettings.HackathonVideoDeadline),
				formatDeadline(input.HackathonVideoDeadline)))
		}
		if isStringPtrChanged(oldSettings.HackathonFinalDeadline, input.HackathonFinalDeadline) {
			changes = append(changes, fmt.Sprintf("deadline final diubah dari \"%s\" menjadi \"%s\"",
				formatDeadline(oldSettings.HackathonFinalDeadline),
				formatDeadline(input.HackathonFinalDeadline)))
		}
		if isStringPtrChanged(oldSettings.ProfileUpdateDeadline, input.ProfileUpdateDeadline) {
			changes = append(changes, fmt.Sprintf("deadline pembaruan profil diubah dari \"%s\" menjadi \"%s\"",
				formatDeadline(oldSettings.ProfileUpdateDeadline),
				formatDeadline(input.ProfileUpdateDeadline)))
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
		HackathonProposalDeadline:     input.HackathonProposalDeadline,
		HackathonVideoDeadline:        input.HackathonVideoDeadline,
		HackathonFinalDeadline:        input.HackathonFinalDeadline,
		ProfileUpdateDeadline:         input.ProfileUpdateDeadline,
		SeminarRegistrationDisabled:   input.SeminarRegistrationDisabled,
		SeminarRequireVerification:    input.SeminarRequireVerification,
	}

	updatedSettings, err := h.settingService.UpdateSettings(newSettings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", updatedSettings))
}

func isStringPtrChanged(oldVal, newVal *string) bool {
	if oldVal == nil && newVal == nil {
		return false
	}
	if oldVal == nil || newVal == nil {
		return true
	}
	return *oldVal != *newVal
}

func formatDeadline(d *string) string {
	if d == nil || *d == "" || *d == "null" {
		return "tidak ada"
	}
	val := *d
	// Try parsing with RFC3339
	t, err := time.Parse(time.RFC3339, val)
	if err != nil {
		// Try standard ISO-like formats
		layouts := []string{
			"2006-01-02T15:04:05",
			"2006-01-02T15:04",
			"2006-01-02 15:04:05",
			"2006-01-02 15:04",
			"2006-01-02",
		}
		for _, layout := range layouts {
			if parsed, parseErr := time.Parse(layout, val); parseErr == nil {
				t = parsed
				err = nil
				break
			}
		}
	}

	if err != nil {
		// Fallback to raw value if parsing fails
		return val
	}

	// Format in Indonesian style: DD MMMM YYYY, HH:mm WIB
	var indoMonths = map[time.Month]string{
		time.January:   "Januari",
		time.February:  "Februari",
		time.March:     "Maret",
		time.April:     "April",
		time.May:       "Mei",
		time.June:      "Juni",
		time.July:      "Juli",
		time.August:    "Agustus",
		time.September: "September",
		time.October:   "Oktober",
		time.November:  "November",
		time.December:  "Desember",
	}

	return fmt.Sprintf("%d %s %d, %02d:%02d WIB", t.Day(), indoMonths[t.Month()], t.Year(), t.Hour(), t.Minute())
}
