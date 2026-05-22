package handler

import (
	"gcw/helper"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type paymentHandler struct {
	registrationService *service.RegistrationService
}

func NewPaymentHandler(rs *service.RegistrationService) *paymentHandler {
	return &paymentHandler{
		registrationService: rs,
	}
}


func (h *paymentHandler) UpdateTeamDetails(c *gin.Context) {
	var payload struct {
		OrderID        string `json:"order_id" binding:"required,max=50"`
		TeamName       string `json:"team_name" binding:"required,max=50"`
		Supervisor     string `json:"supervisor" binding:"max=50"`
		SupervisorNIDN string `json:"supervisor_nidn" binding:"max=50"`
		ReceiptLink    string `json:"receipt_link" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	if payload.OrderID == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"order_id": {"IS_REQUIRED"}}))
		return
	}

	// Update DB
	changes, err := h.registrationService.UpdateTeamDetails(payload.OrderID, payload.TeamName, payload.Supervisor, payload.SupervisorNIDN, payload.ReceiptLink)
	if err != nil {
		if err.Error() == "Nama Tim Sudah Digunakan" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Nama tim sudah digunakan"))
			return
		}
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal memperbarui detail tim", helper.FormatValidationError(err)))
		return
	}

	// Set target name and changes for audit log
	c.Set("target_name", payload.TeamName)
	c.Set("audit_changes", changes)

	c.JSON(http.StatusOK, helper.CreateMutationResponse("Detail tim berhasil diperbarui", gin.H{
		"team_name":       payload.TeamName,
		"supervisor":      payload.Supervisor,
		"supervisor_nidn": payload.SupervisorNIDN,
		"receipt_link":    payload.ReceiptLink,
		"changes":         changes,
	}))
}
