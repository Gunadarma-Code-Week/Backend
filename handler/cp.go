package handler

import "github.com/gin-gonic/gin"

type CompetitiveHandler struct {
}

func GateCompetitiveHandler() *CompetitiveHandler {
	return &CompetitiveHandler{}
}

func (h *CompetitiveHandler) Ping(c *gin.Context) {
}
