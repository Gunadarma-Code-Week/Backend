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

// Admin User Management Handlers

// @Summary Get All Users (Admin)
// @Description Get all users with pagination, filtering, and sorting for admin
// @Tags Admin Users
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Param q query string false "Search query"
// @Param sortBy query string false "Sort field" Enums(id,institusi,id_team,nim,soc_med_document,profile_has_updated,data_has_verified)
// @Param sortOrder query string false "Sort order" Enums(ASC,DESC)
// @Success 200 {object} helper.Response{data=dto.AdminUsersListResponseDTO}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 403 {object} helper.Response{data=string} "Forbidden"
// @Router /admin/users [get]
func (h *UserHandler) AdminGetAllUsers(c *gin.Context) {
	var query dto.AdminGetUsersQueryDTO

	// Set defaults
	query.Page = 1
	query.Limit = 10
	query.SortBy = "id"
	query.SortOrder = "ASC"

	// Bind query parameters
	if err := c.ShouldBindQuery(&query); err != nil {
		logging.Low("UserHandler.AdminGetAllUsers", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	// Call service
	response, err := h.userService.AdminGetAllUsers(query)
	if err != nil {
		logging.High("UserHandler.AdminGetAllUsers", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Users retrieved successfully", response))
}

// @Summary Get User by ID (Admin)
// @Description Get a single user by ID for admin
// @Tags Admin Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} helper.Response{data=dto.AdminUserResponseDTO}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 404 {object} helper.Response{data=string} "Not Found"
// @Router /admin/users/{id} [get]
func (h *UserHandler) AdminGetUserById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid user ID"))
		return
	}

	response, err := h.userService.AdminGetUserById(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, helper.CreateErrorResponse("error", "User not found"))
			return
		}
		logging.High("UserHandler.AdminGetUserById", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("User retrieved successfully", response))
}

// @Summary Update User (Admin)
// @Description Update an existing user by admin
// @Tags Admin Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.AdminUpdateUserDTO true "Update User Data"
// @Success 200 {object} helper.Response{data=dto.AdminUserResponseDTO}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 404 {object} helper.Response{data=string} "Not Found"
// @Router /admin/users/{id} [put]
func (h *UserHandler) AdminUpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid user ID"))
		return
	}

	var updateData dto.AdminUpdateUserDTO
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logging.Low("UserHandler.AdminUpdateUser", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	response, err := h.userService.AdminUpdateUser(id, updateData)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, helper.CreateErrorResponse("error", "User not found"))
			return
		}
		logging.High("UserHandler.AdminUpdateUser", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("User updated successfully", response))
}

// @Summary Delete User (Admin)
// @Description Delete a user by ID (admin only)
// @Tags Admin Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} helper.Response{data=string}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Failure 404 {object} helper.Response{data=string} "Not Found"
// @Router /admin/users/{id} [delete]
func (h *UserHandler) AdminDeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", "Invalid user ID"))
		return
	}

	err = h.userService.AdminDeleteUser(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, helper.CreateErrorResponse("error", "User not found"))
			return
		}
		logging.High("UserHandler.AdminDeleteUser", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("User deleted successfully", "User has been deleted"))
}

// @Summary Get User Growth Analytics (Admin)
// @Description Get user growth analytics between two dates
// @Tags Admin Users
// @Accept json
// @Produce json
// @Param startDate query string true "Start date (YYYY-MM-DD)"
// @Param endDate query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} helper.Response{data=[]dto.UserGrowthResponseDTO}
// @Failure 400 {object} helper.Response{data=string} "Bad Request"
// @Router /admin/users/analytics/growth [get]
func (h *UserHandler) AdminGetUserGrowthAnalytics(c *gin.Context) {
	var query dto.UserGrowthAnalyticsDTO

	if err := c.ShouldBindQuery(&query); err != nil {
		logging.Low("UserHandler.AdminGetUserGrowthAnalytics", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	response, err := h.userService.AdminGetUserGrowthAnalytics(query)
	if err != nil {
		logging.High("UserHandler.AdminGetUserGrowthAnalytics", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("User growth analytics retrieved successfully", response))
}
