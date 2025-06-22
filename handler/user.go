package handler

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

// @Summary Get My Profile Data
// @Tags Profile
// @Produce  json
// @Success 200 {object} helper.Response{data=dto.UserResponseDTO}
// @Router /profile/my [get]
func (h *UserHandler) GetMyProfile(c *gin.Context) {
	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("AuthHandler.Login", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "user not found in context"))
		return
	}

	user := &dto.UserResponseDTO{}
	smapping.FillStruct(user, smapping.MapFields(userAuth))

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", user))
}

// @Summary Update My Profile Data
// @Tags Profile
// @Accept json
// @Produce  json
// @Param request body dto.UpdateUserProfileDTO true "Update User Profile"
// @Success 200 {object} helper.Response{data=dto.UserResponseDTO}
// @Router /profile/my [post]
func (h *UserHandler) UpdateMyProfile(c *gin.Context) {
	userUpdateDTO := &dto.UpdateUserProfileDTO{}
	if err := c.Bind(userUpdateDTO); err != nil {
		logging.Low("AuthHandler.Login", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("AuthHandler.Login", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "user not found in context"))
		return
	}

	userUpdate := &entity.User{}
	smapping.FillStruct(userUpdate, smapping.MapFields(userUpdateDTO))
	// "YYYY-MM-DD" convert to time.Time
	birthDate, err := time.Parse("2006-01-02", userUpdateDTO.BirthDate)
	if err != nil {
		logging.Low("UserHandler.UpdateMyProfile", "BAD_REQUEST", "invalid birth date format")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "invalid birth date format"))
		return
	}
	userUpdate.BirthDate = &birthDate
	userUpdate.Phone = userUpdateDTO.Phone
	userUpdate.Major = userUpdateDTO.Major

	err = h.userService.Update(userUpdate, userAuth.ID)
	if err != nil {
		logging.High("AuthHandler.Login", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	user, err := h.userService.FindById(userAuth.ID)
	if err != nil {
		logging.High("AuthHandler.Login", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	userResponse := &dto.UserResponseDTO{}
	smapping.FillStruct(userResponse, smapping.MapFields(user))

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", userResponse))
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	startDateStr := c.Param("start_date")
	endDateStr := c.Param("end_date")
	countStr := c.Param("count")
	pageStr := c.Param("page")

	// Convert 'count' and 'page' to integers
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	// Parse start_date and end_date into time.Time
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	// Calculate the offset for pagination
	offset := (page - 1) * count

	// Fetch paginated data from the service
	users, totalUsers, err := h.userService.GetUsersByDateRange(startDate, endDate, count, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Prepare the response with mapped DTOs
	userResponses := []dto.UserResponseDTO{}
	for _, user := range users {
		var userResponse dto.UserResponseDTO
		err := smapping.FillStruct(&userResponse, smapping.MapFields(user))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error mapping user data"})
			return
		}
		userResponses = append(userResponses, userResponse)
	}

	// Calculate total pages
	totalPages := (totalUsers + int64(count) - 1) / int64(count)

	has_more := false
	if totalPages > int64(page) {
		has_more = true
	}

	response := gin.H{
		"status":      "success",
		"message":     "success",
		"data":        userResponses,
		"totalItems":  totalUsers,
		"totalPages":  totalPages + 1,
		"currentPage": page,
		"count":       count,
		"has_more":    has_more,
	}

	// Return a paginated response
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))
}

// @Summary Get User Events
// @Tags Profile
// @Produce json
// @Success 200 {object} helper.Response{data=dto.ResponseEvents}
// @Router /profile/events [get]
func (h *UserHandler) GetEvents(c *gin.Context) {
	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("UserHandler.GetEvents", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "user not found in context"))
		return
	}

	dataEvent, err := h.userService.GetEvents(userAuth.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", "data not found"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("FOUND", dataEvent))
}
