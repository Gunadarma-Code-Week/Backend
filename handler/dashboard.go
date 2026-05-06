package handler

import (
	"fmt"
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
	fmt.Println("[dashboard.GetAllDashboard] START:", c.Request.URL.Path)
	acara := c.Param("acara")
	strCount := c.Param("count")
	strPage := c.Param("page")
	startDate := c.Param("start_date")
	endDate := c.Param("end_date")
	search := c.Query("search")

	fmt.Printf("[dashboard.GetAllDashboard] Params: acara=%s, start=%s, end=%s, count=%s, page=%s\n", 
		acara, startDate, endDate, strCount, strPage)

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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "count and page must be integers"))
		return
	}


	fmt.Println("[dashboard.GetAllDashboard] calling service for", acara)
	var respondData interface{}

	switch acara {
	case "seminar":
		fmt.Println("[dashboard.GetAllDashboard] entering seminar branch")
		data, err := h.Service.GetAllSeminar(count, page, search)
		if err != nil {
			fmt.Println("[dashboard.GetAllDashboard] seminar service error:", err)
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "hackaton", "hackathon":
		fmt.Println("[dashboard.GetAllDashboard] entering hackaton branch")
		data, err := h.Service.GetAllHackaton(count, page)
		if err != nil {
			fmt.Println("[dashboard.GetAllDashboard] hackaton service error:", err)
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "cp":
		fmt.Println("[dashboard.GetAllDashboard] entering cp branch")
		data, err := h.Service.GetAllCp(count, page)
		if err != nil {
			fmt.Println("[dashboard.GetAllDashboard] cp service error:", err)
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	case "ctf":
		fmt.Println("[dashboard.GetAllDashboard] entering ctf branch")
		data, err := h.Service.GetAllCtf(count, page)
		if err != nil {
			fmt.Println("[dashboard.GetAllDashboard] ctf service error:", err)
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "error service"))
			return
		}

		respondData = data

	default:
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "kegiatan not found"))
		return
	}

	fmt.Println("[dashboard.GetAllDashboard] sending response for", acara, "with data:", respondData)
	
	// Final safety check: if respondData is still nil, return empty object instead of nothing
	if respondData == nil {
		fmt.Println("[dashboard.GetAllDashboard] respondData is NIL, falling back to empty object")
		respondData = map[string]interface{}{}
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("SUCCESS", respondData))
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
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		targetName, err := h.Service.UpdateSeminarService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}
		c.Set("target_name", targetName)

	case "hackathon":
		var input dto.Hackaton
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		targetName, err := h.Service.UpdateHackatonService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}
		c.Set("target_name", targetName)

	case "cp":
		var input dto.Cp
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		targetName, err := h.Service.UpdateCpService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}
		c.Set("target_name", targetName)

	case "ctf":
		var input dto.Ctf
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "BAD_REQUEST"))
			return
		}

		targetName, err := h.Service.UpdateCtfService(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("ERROR", "error service"))
			return
		}
		c.Set("target_name", targetName)

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
// @Param acara path string true "Event type (seminar, hackaton, cp, ctf)"
// @Param id path string true "Event ID"
// @Success 200 {object} helper.Response{data=string}
// @Failure 400 {object} helper.Response{message=string}
// @Router /dashboard/{acara}/{id} [delete]
func (h *dashboardController) Delete(c *gin.Context) {
	acara := c.Param("acara")
	id := c.Param("id")

	if acara != "seminar" && acara != "hackaton" && acara != "hackathon" && acara != "cp" && acara != "ctf" {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "kegiatan not found"))
		return
	}

	targetName, err := h.Service.DeletePesertaService(acara, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "Error delete service"))
		return
	}
	c.Set("target_name", targetName)

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("UPDATED", targetName))
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
