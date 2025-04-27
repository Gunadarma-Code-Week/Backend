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

// @Summary Get Newsletter
// @Description Get a newsletter by ID
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param id path uint64 true "Newsletter ID"
// @Success 200 {object} helper.Response{data=entity.NewsLetter}
// @Failure 400 {object} helper.Response{data=string} "Invalid ID"
// @Failure 404 {object} helper.Response{data=string} "Not found"
// @Router /newsletter/{id} [get]
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

// @Summary Create Newsletter
// @Description Create a new newsletter
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param request body dto.CreateNewsLetterDTO true "Create Newsletter"
// @Success 201 {object} helper.Response{data=entity.NewsLetter}
// @Failure 400 {object} helper.Response{data=string} "Invalid input"
// @Failure 500 {object} helper.Response{data=string} "Create failed"
// @Router /newsletter [post]
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

// @Summary Update Newsletter
// @Description Update a newsletter by ID
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param id path uint64 true "Newsletter ID"
// @Param request body dto.UpdateNewsLetterDTO true "Update Newsletter"
// @Success 200 {object} helper.Response{data=entity.NewsLetter}
// @Failure 400 {object} helper.Response{data=string} "Invalid input or ID"
// @Failure 500 {object} helper.Response{data=string} "Update failed"
// @Router /newsletter/{id} [put]
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

// @Summary Delete Newsletter
// @Description Delete a newsletter by ID
// @Tags Newsletter
// @Accept  json
// @Produce  json
// @Param id path uint64 true "Newsletter ID"
// @Success 200 {object} helper.Response{data=string} "Newsletter deleted"
// @Failure 400 {object} helper.Response{data=string} "Invalid ID"
// @Failure 500 {object} helper.Response{data=string} "Delete failed"
// @Router /newsletter/{id} [delete]
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
