package handler

import (
	"errors"
	"gcw/helper"
	"gcw/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CTFHandler struct {
	ctfService *service.CtfService
}

func NewCTFHandler(s *service.CtfService) *CTFHandler {
	return &CTFHandler{
		ctfService: s,
	}
}

func (h *CTFHandler) GetDetail(c *gin.Context) {
	join_code := c.Param("join_code")
	result, err := h.ctfService.Get(join_code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("Detail tim CTF tidak ditemukan"))
			return
		}
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Gagal mengambil detail tim CTF", helper.FormatValidationError(err)))
		return
	}
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", result))
}

