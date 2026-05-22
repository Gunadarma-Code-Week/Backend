package handler

import (
	"errors"
	"gcw/helper"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Detail tim kompetitif tidak ditemukan"))
			return
		}
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Gagal mengambil detail tim kompetitif", helper.FormatValidationError(err)))
		return
	}
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", result))
}

