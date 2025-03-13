package handler

import "github.com/gin-gonic/gin"

type competitiveHandler struct {
}

type CompetitiveHandler interface {
}

func GateCompetitiveHandler() CompetitiveHandler {
	return &competitiveHandler{}
}

func (h *competitiveHandler) Ping(c *gin.Context) {
}
