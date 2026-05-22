package handler

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SeminarHandler struct {
	seminarService *service.SeminarService
}

func NewSeminarHandler(ss *service.SeminarService) *SeminarHandler {
	return &SeminarHandler{
		seminarService: ss,
	}
}

// @Summary Join Seminar
// @Description User bergabung ke seminar (ID tiket akan di-generate otomatis)
// @Tags Seminar
// @Accept json
// @Produce json
// @Success 201 {object} helper.Response{data=dto.JoinSeminarResponse}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 500 {object} helper.Response{data=string} "Internal Server Error"
// @Router /seminar/join [post]
func (h *SeminarHandler) JoinSeminar(c *gin.Context) {
	// Tidak perlu bind request karena tidak ada field yang diperlukan
	var request dto.JoinSeminarRequest

	// Get user from context
	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("SeminarHandler.JoinSeminar", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "User tidak ditemukan di context"))
		return
	}

	// Call service
	response, err := h.seminarService.JoinSeminar(userAuth.ID, request)
	if err != nil {
		logging.Low("SeminarHandler.JoinSeminar", "BAD_REQUEST", err.Error())
		if err.Error() == "user sudah terdaftar di seminar" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("User sudah terdaftar di seminar"))
			return
		}
		if err.Error() == "seminar sudah penuh, maksimal 100 participant" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Seminar sudah penuh"))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(err.Error(), helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Berhasil bergabung ke seminar", response))
}

// @Summary Get My Seminar Ticket
// @Description Mendapatkan detail tiket seminar user yang sedang login
// @Tags Seminar
// @Produce json
// @Success 200 {object} helper.Response{data=dto.SeminarTicketDetail}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 404 {object} helper.Response{data=string} "Not Found"
// @Router /seminar/my-ticket [get]
func (h *SeminarHandler) GetMyTicket(c *gin.Context) {
	// Get user from context
	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("SeminarHandler.GetMyTicket", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "User tidak ditemukan di context"))
		return
	}

	// Call service
	response, err := h.seminarService.GetTicketDetail(userAuth.ID)
	if err != nil {
		if err.Error() == "tiket seminar tidak ditemukan" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Tiket seminar tidak ditemukan"))
			return
		}
		logging.Low("SeminarHandler.GetMyTicket", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal mengambil tiket seminar", helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Detail tiket seminar berhasil ditemukan", response))
}

// @Summary Get Seminar Ticket by ID
// @Description Mendapatkan detail tiket seminar berdasarkan ID tiket (untuk admin)
// @Tags Seminar
// @Produce json
// @Param ticket_id path string true "Ticket ID"
// @Success 200 {object} helper.Response{data=dto.SeminarTicketDetail}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 404 {object} helper.Response{data=string} "Not Found"
// @Router /seminar/ticket/{ticket_id} [get]
func (h *SeminarHandler) GetTicketByID(c *gin.Context) {
	ticketID := c.Param("ticket_id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"ticket_id": {"IS_REQUIRED"}}))
		return
	}

	// Call service
	response, err := h.seminarService.GetTicketByID(ticketID)
	if err != nil {
		if err.Error() == "tiket seminar tidak ditemukan" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Tiket seminar tidak ditemukan"))
			return
		}
		logging.Low("SeminarHandler.GetTicketByID", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal mengambil tiket seminar", helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Detail tiket seminar berhasil ditemukan", response))
}

// @Summary Admin Add Participant to Seminar
// @Description Admin menambahkan participant ke seminar berdasarkan user ID
// @Tags Seminar
// @Accept json
// @Produce json
// @Param request body dto.AdminAddParticipantRequest true "User ID"
// @Success 201 {object} helper.Response{data=dto.AdminAddParticipantResponse}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 500 {object} helper.Response{data=string} "Internal Server Error"
// @Router /seminar/admin/add-participant [post]
func (h *SeminarHandler) AdminAddParticipant(c *gin.Context) {
	var request dto.AdminAddParticipantRequest

	// Bind request
	if err := c.ShouldBindJSON(&request); err != nil {
		logging.Low("SeminarHandler.AdminAddParticipant", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	// Call service
	response, err := h.seminarService.AdminAddParticipant(request.UserID)
	if err != nil {
		logging.Low("SeminarHandler.AdminAddParticipant", "BAD_REQUEST", err.Error())
		if err.Error() == "user tidak ditemukan" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("User tidak ditemukan"))
			return
		}
		if err.Error() == "user sudah terdaftar di seminar" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("User sudah terdaftar di seminar"))
			return
		}
		if err.Error() == "seminar sudah penuh, maksimal 100 participant" {
			c.JSON(http.StatusConflict, helper.CreateConflictResponse("Seminar sudah penuh"))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse(err.Error(), helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Berhasil menambahkan participant ke seminar", response))
}