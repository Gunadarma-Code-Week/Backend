package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/service"
	"net/http"
	"strconv"
	"time"

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
	// GetEvent(*gin.Context)
}

func DashboardController(db *gorm.DB) DashboardControllerInterface {
	return &dashboardController{
		Service: service.NewDashboardServices(db),
	}
}

func (h *dashboardController) Statistics(c *gin.Context) {}

// @Summary Get All Dashboard
// @Description Retrieve all dashboard data based on the specified event type (seminar, hackaton, cp).
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param acara path string true "Event type (seminar, hackaton, cp)"
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
	startDateStr := c.Param("start_date")
	endDateStr := c.Param("end_date")
	search := c.Query("search")

	count, errCount := strconv.Atoi(strCount)
	page, errPage := strconv.Atoi(strPage)

	if errCount != nil || errPage != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "count or page are error"))
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	var respondData interface{}

	switch acara {
	case "seminar":
		data, err := h.Service.GetAllSeminar(startDate, endDate, count, page, search)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "hackaton":
		data, err := h.Service.GetAllHackaton(startDate, endDate, count, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "cp":
		data, err := h.Service.GetAllCp(startDate, endDate, count, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	default:
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "kegiatan not found"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("SUCCESS", respondData))
}

// @Summary Update Dashboard Event
// @Description Update a specific dashboard event based on the event type and ID.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param acara path string true "Event type (seminar, hackaton, cp)"
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
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		if err := h.Service.UpdateSeminarService(id, input); err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}

	case "hackathon":
		var input dto.Hackaton
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		if err := h.Service.UpdateHackatonService(id, input); err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}

	case "cp":
		var input dto.Cp
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		if err := h.Service.UpdateCpService(id, input); err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}

	default:
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "kegiatan not found"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("UPDATED", id))
}

// @Summary Delete Dashboard Event
// @Description Delete a specific dashboard event based on the event type and ID.
// @Tags Dashboard
// @Accept  json
// @Produce  json
// @Param acara path string true "Event type (seminar, hackaton, cp)"
// @Param id path string true "Event ID"
// @Success 200 {object} helper.Response{data=string}
// @Failure 400 {object} helper.Response{message=string}
// @Router /dashboard/{acara}/{id} [delete]
func (h *dashboardController) Delete(c *gin.Context) {
	acara := c.Param(":acara")
	id := c.Param(":id")

	if acara != "seminar" && acara != "hackaton" && acara != "cp" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "kegiatan not found"))
		return
	}

	id_user, err := h.Service.DeletePesertaService(acara, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "Error delete service"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("UPDATED", id_user))
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
// 		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "data not found"))
// 		return
// 	}

// 	c.JSON(http.StatusOK, helper.CreateSuccessResponse("FOUND", dataEvent))
// }
