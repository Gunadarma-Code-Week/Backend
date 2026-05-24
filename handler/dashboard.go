package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/service"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type dashboardController struct {
	Service service.DashboardServices
}

type DashboardControllerInterface interface {
	Statistics(*gin.Context)
	GetAllDashboard(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	GetTeamNames(*gin.Context)
	SendBulkEmail(*gin.Context)
	GetEmailTemplates(*gin.Context)
	// GetEvent(*gin.Context)
}

func DashboardController(db *gorm.DB) DashboardControllerInterface {
	return &dashboardController{
		Service: service.NewDashboardServices(db),
	}
}

func (h *dashboardController) Statistics(c *gin.Context) {}

// @Summary Get All Dashboard
// @Description Retrieve all dashboard data based on the specified event type (seminar, hackaton, cp, ctf).
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param acara path string true "Event type (seminar, hackaton, cp, ctf)"
// @Param count path int true "Number of items per page"
// @Param page path int true "Page number"
// @Param search query string false "Search by id_tiket, name, or email"
// @Success 200 {object} helper.Response{data=interface{}}
// @Failure 400 {object} helper.Response{message=string}
// @Router /dashboard/{acara}/{count}/{page} [get]
func (h *dashboardController) GetAllDashboard(c *gin.Context) {
	acara := c.Param("acara")
	strCount := c.Param("count")
	strPage := c.Param("page")

	// Fallback logic if path params are missing
	if strCount == "" {
		strCount = c.Query("count")
	}
	if strPage == "" {
		strPage = c.Query("page")
	}

	// Default values if still empty
	if strCount == "" {
		strCount = "10"
	}
	if strPage == "" {
		strPage = "0"
	}

	count, errCount := strconv.Atoi(strCount)
	page, errPage := strconv.Atoi(strPage)

	if errCount != nil || errPage != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Parameter count dan page harus berupa angka"))
		return
	}

	var respondData interface{}
	switch acara {
	case "seminar":
		search := c.Query("search")
		data, err := h.Service.GetAllSeminar(count, page, search)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}
		respondData = data
	case "hackaton", "hackathon":
		search := c.Query("search")
		stage := c.Query("stage")
		data, err := h.Service.GetAllHackaton(count, page, search, stage)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}
		respondData = data
	case "cp":
		search := c.Query("search")
		stage := c.Query("stage")
		data, err := h.Service.GetAllCp(count, page, search, stage)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}
		respondData = data
	case "ctf":
		search := c.Query("search")
		stage := c.Query("stage")
		data, err := h.Service.GetAllCtf(count, page, search, stage)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}
		respondData = data
	default:
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Acara tidak ditemukan"))
		return
	}

	// Final safety check: if respondData is still nil, return empty object instead of nothing
	if respondData == nil {
		respondData = map[string]interface{}{}
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", respondData))
}

// @Summary Update Dashboard Event
// @Description Update a specific dashboard event based on the event type and ID.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param acara path string true "Event type (seminar, hackaton, cp, ctf)"
// @Param id path string true "Event ID"
// @Param request body interface{} true "Event data to update"
// @Success 200 {object} helper.Response{data=string}
// @Failure 400 {object} helper.Response{message=string}
// @Failure 500 {object} helper.Response{message=string}
// @Router /dashboard/{acara}/{id} [put]
func (h *dashboardController) Update(c *gin.Context) {
	acara := c.Param("acara")
	id := c.Param("id")

	switch acara {
	case "seminar":
		var input dto.Seminar
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}

		targetName, err := h.Service.UpdateSeminarService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal memperbarui seminar", helper.FormatValidationError(err)))
			return
		}
		c.Set("target_name", targetName)

	case "hackathon":
		var input dto.Hackaton
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}

		targetName, err := h.Service.UpdateHackatonService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal memperbarui hackathon", helper.FormatValidationError(err)))
			return
		}
		c.Set("target_name", targetName)

	case "cp":
		var input dto.Cp
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}

		targetName, err := h.Service.UpdateCpService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal memperbarui CP", helper.FormatValidationError(err)))
			return
		}
		c.Set("target_name", targetName)

	case "ctf":
		var input dto.Ctf
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
			return
		}

		targetName, err := h.Service.UpdateCtfService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal memperbarui CTF", helper.FormatValidationError(err)))
			return
		}
		c.Set("target_name", targetName)

	default:
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Acara tidak ditemukan"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateMutationResponse("Data berhasil diperbarui", id))
}

// @Summary Delete Dashboard Event
// @Description Delete a specific dashboard event based on the event type and ID.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param acara path string true "Event type (seminar, hackaton, cp, ctf)"
// @Param id path string true "Event ID"
// @Success 200 {object} helper.Response{data=string}
// @Failure 400 {object} helper.Response{message=string}
// @Router /dashboard/{acara}/{id} [delete]
func (h *dashboardController) Delete(c *gin.Context) {
	acara := c.Param("acara")
	id := c.Param("id")

	if acara != "seminar" && acara != "hackaton" && acara != "hackathon" && acara != "cp" && acara != "ctf" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Acara tidak ditemukan"))
		return
	}

	var deleteRequest dto.DeleteTeamRequest
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	targetName, err := h.Service.DeletePesertaService(acara, id, deleteRequest.Alasan)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Gagal menghapus peserta", helper.FormatValidationError(err)))
		return
	}
	c.Set("target_name", targetName)

	c.JSON(http.StatusOK, helper.CreateMutationResponse("Peserta berhasil dihapus", targetName))
}

// @Summary Get User Events
// @Description Retrieve all events associated with a specific user.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param id_user path string true "User ID"
// @Success 200 {object} helper.Response{data=dto.ResponseEvents}
// @Failure 400 {object} helper.Response{message=string}
// @Router /dashboard/events/{id_user} [get]
// func (h *dashboardController) GetEvent(c *gin.Context) {
// 	idUser := c.Param("id_user")

// 	dataEvent, err := h.Service.GetEventSevice(idUser)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "data not found"))
// 		return
// 	}

// 	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Data berhasil ditemukan", dataEvent))
// }

func (h *dashboardController) GetTeamNames(c *gin.Context) {
	event := c.Query("event")
	if event == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Parameter event wajib diisi"))
		return
	}

	names, err := h.Service.GetTeamNames(event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terjadi kesalahan pada server", helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", names))
}

func (h *dashboardController) SendBulkEmail(c *gin.Context) {
	var input struct {
		Event      string `json:"event" binding:"max=50"`
		Stage      string `json:"stage" binding:"max=50"`
		TargetRole string `json:"target_role" binding:"max=50"`
		Subject    string `json:"subject" binding:"max=150"`
		Content    string `json:"content"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	if input.Event == "" || input.Subject == "" || input.Content == "" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Parameter event, subject, dan content wajib diisi"))
		return
	}

	err := h.Service.SendBulkEmail(input.Event, input.Stage, input.TargetRole, input.Subject, input.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Permintaan berhasil diproses",
		"data":    "Email berhasil dimasukkan ke dalam antrean pengiriman",
		"errors":  "nil",
	})
}

func (h *dashboardController) GetEmailTemplates(c *gin.Context) {
	templates, err := h.Service.GetEmailTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Terjadi kesalahan pada server", helper.FormatValidationError(err)))
		return
	}
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", templates))
}
