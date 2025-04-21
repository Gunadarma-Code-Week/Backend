package handler

import (
	"gcw/dto"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type newsletterHandler struct {
	service service.NewsletterService
}

func NewNewsletterHandler(s service.NewsletterService) *newsletterHandler {
	return &newsletterHandler{service: s}
}

func (h *newsletterHandler) CreateNewsletter(c *gin.Context) {
	var input dto.CreateNewsLetterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Invalid input", err.Error()))
		return
	}

	newsletter, err := h.service.Create(input)
	if err != nil {
		logging.Warn("CreateNewsletter", "error create newsletter", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Create failed", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateSuccessResponse("Newsletter created", newsletter))
}

func (h *newsletterHandler) GetNewsLetter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Invalid ID", err.Error()))
		return
	}

	newsletter, err := h.service.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, helper.CreateErrorResponse("Not found", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Success", newsletter))
}

func (h *newsletterHandler) UpdateNewsLetter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Invalid ID", err.Error()))
		return
	}

	var input dto.UpdateNewsLetterDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Invalid input", err.Error()))
		return
	}

	newsletter, err := h.service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Update failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Newsletter updated", newsletter))
}

func (h *newsletterHandler) DeleteNewsLetter(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Invalid ID", err.Error()))
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("Delete failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Newsletter deleted", nil))
}
