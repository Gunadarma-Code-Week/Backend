package handler

import (
	"gcw/service"

	"github.com/gin-gonic/gin"
)

type CompetitiveHandler struct {
	cpService *service.CpService
}

func GateCompetitiveHandler(s *service.CpService) *CompetitiveHandler {
	return &CompetitiveHandler{
		cpService: s,
	}
}

// @Summary Get CP Details
// @Description Get CP details by join code
// @Tags CP
// @Accept  json
// @Produce  json
// @Param join_code path string true "Join Code"
// @Success 200 {object} helper.Response{data=dto.CpDetailDto}
// @Failure 400 {object} helper.Response{data=string} "Invalid join code"
// @Failure 404 {object} helper.Response{data=string} "CP details not found"
// @Failure 500 {object} helper.Response{data=string} "Internal server error"
// @Router /cp/{join_code} [get]
func (h *CompetitiveHandler) GetDetail(c *gin.Context) {
	join_code := c.Param("join_code")
	result, err := h.cpService.Get(join_code)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": result})
}
