package handler

import (
	"gcw/service"

	"github.com/gin-gonic/gin"
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
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": result})
}
