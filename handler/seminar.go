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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "user not found in context"))
		return
	}

	// Call service
	response, err := h.seminarService.JoinSeminar(userAuth.ID, request)
	if err != nil {
		logging.Low("SeminarHandler.JoinSeminar", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Berhasil bergabung ke seminar", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "user not found in context"))
		return
	}

	// Call service
	response, err := h.seminarService.GetTicketDetail(userAuth.ID)
	if err != nil {
		if err.Error() == "tiket seminar tidak ditemukan" {
			c.JSON(http.StatusNotFound, helper.CreateErrorResponse("error", err.Error()))
			return
		}
		logging.Low("SeminarHandler.GetMyTicket", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Detail tiket seminar", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "ticket_id is required"))
		return
	}

	// Call service
	response, err := h.seminarService.GetTicketByID(ticketID)
	if err != nil {
		if err.Error() == "tiket seminar tidak ditemukan" {
			c.JSON(http.StatusNotFound, helper.CreateErrorResponse("error", err.Error()))
			return
		}
		logging.Low("SeminarHandler.GetTicketByID", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Detail tiket seminar", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	// Call service
	response, err := h.seminarService.AdminAddParticipant(request.UserID)
	if err != nil {
		logging.Low("SeminarHandler.AdminAddParticipant", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Berhasil menambahkan participant ke seminar", response))
}