package handler

import (
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type paymentHandler struct {
	midtransService *service.MidtransService
	registrationService *service.RegistrationService
}

func NewPaymentHandler(ms *service.MidtransService, rs *service.RegistrationService) *paymentHandler {
	return &paymentHandler{
		midtransService: ms,
		registrationService: rs,
	}
}

// @Summary Midtrans Payment Notification
// @Description Callback endpoint for Midtrans to notify payment status
// @Tags Payment
// @Accept json
// @Produce json
// @Router /payment/notification [post]
func (h *paymentHandler) Notification(c *gin.Context) {
	var notificationPayload map[string]interface{}
	if err := c.ShouldBindJSON(&notificationPayload); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	orderID, _ := notificationPayload["order_id"].(string)
	
	// Verify notification from Midtrans
	res, err := h.midtransService.Client.CheckTransaction(orderID)
	if err != nil {
		logging.Low("PaymentHandler.Notification", "INTERNAL_SERVER_ERROR", "CheckTransaction failed: "+err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to verify transaction"))
		return
	}

	if res != nil {
		// Handle transaction status
		// For QRIS, settlement means payment is successful
		if res.TransactionStatus == "settlement" {
			err := h.registrationService.UpdatePaymentStatus(orderID, "Paid")
			if err != nil {
				logging.Low("PaymentHandler.Notification", "INTERNAL_SERVER_ERROR", "UpdatePaymentStatus failed: "+err.Error())
				c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to update payment status"))
				return
			}
		} else if res.TransactionStatus == "expire" || res.TransactionStatus == "cancel" {
			_ = h.registrationService.UpdatePaymentStatus(orderID, "Failed")
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *paymentHandler) ManualCheckTransaction(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Order ID is required"))
		return
	}

	// If it's a CP team (manual payment), skip Midtrans check
	if len(orderID) > 3 && orderID[:3] == "CP-" {
		// Just return the current status from DB
		team := &entity.Team{}
		err := h.registrationService.Repository().DB.Where("order_id = ?", orderID).First(team).Error
		if err != nil {
			c.JSON(http.StatusNotFound, helper.CreateErrorResponse("error", "Team not found"))
			return
		}

		c.JSON(http.StatusOK, helper.CreateSuccessResponse("Manual payment status", gin.H{
			"status":             team.PaymentStatus,
			"transaction_status": "manual",
		}))
		return
	}

	res, err := h.midtransService.Client.CheckTransaction(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to check transaction: "+err.Error()))
		return
	}

	status := "Pending"
	if res.TransactionStatus == "settlement" || res.TransactionStatus == "capture" {
		status = "Paid"
	} else if res.TransactionStatus == "expire" || res.TransactionStatus == "cancel" || res.TransactionStatus == "deny" {
		status = "Failed"
	}

	errUpdate := h.registrationService.UpdatePaymentStatus(orderID, status)
	if errUpdate != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", "Failed to update internal status"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Transaction checked", gin.H{
		"status": status,
		"transaction_status": res.TransactionStatus,
	}))
}

