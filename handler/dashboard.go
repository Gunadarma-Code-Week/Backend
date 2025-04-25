package handler

import (
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
	Seminar(*gin.Context)
	Hackaton(*gin.Context)
	Cp(*gin.Context)
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

func (h *dashboardController) Seminar(c *gin.Context) {}

func (h *dashboardController) Hackaton(c *gin.Context) {}

func (h *dashboardController) Cp(c *gin.Context) {}
