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
	GetEvent(*gin.Context)
}

func DashboardController(db *gorm.DB) DashboardControllerInterface {
	return &dashboardController{
		Service: service.NewDashboardServices(db),
	}
}

func (h *dashboardController) Statistics(c *gin.Context) {}

func (h *dashboardController) GetAllDashboard(c *gin.Context) {
	acara := c.Param("acara")
	strCount := c.Param("count")
	strPage := c.Param("page")

	count, errCount := strconv.Atoi(strCount)
	page, errPage := strconv.Atoi(strPage)

	if errCount != nil || errPage != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "count or page are error"))
	}

	var respondData interface{}

	switch acara {
	case "seminar":
		data, err := h.Service.GetAllSeminar(count, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "hackaton":
		data, err := h.Service.GetAllHackaton(count, page)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "cp":
		data, err := h.Service.GetAllCp(count, page)
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

func (h *dashboardController) Update(c *gin.Context) {
	acara := c.Param(":acara")
	id := c.Param(":id")

	switch acara {
	case "seminar":
		var input dto.Seminar
		if err := c.ShouldBindJSON(&input).Error; err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		if err := h.Service.UpdateSeminarService(id, input); err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}

	case "hackaton":
		var input dto.Hackaton
		if err := c.ShouldBindJSON(&input).Error; err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		if err := h.Service.UpdateHackatonService(id, input); err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}

	case "cp":
		var input dto.Cp
		if err := c.ShouldBindJSON(&input).Error; err != nil {
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

func (h *dashboardController) GetEvent(c *gin.Context) {
	idUser := c.Param("id_user")

	dataEvent, err := h.Service.GetEventSevice(idUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "data not found"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("FOUND", dataEvent))
}
