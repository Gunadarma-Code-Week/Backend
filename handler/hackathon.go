package handler

import "github.com/gin-gonic/gin"

type hackathonHandler struct {
}

type HackathonHandler interface {
}

func GateHackathonHandler() HackathonHandler {
	return &hackathonHandler{}
}

func (h *hackathonHandler) Ping(c *gin.Context) {
}
