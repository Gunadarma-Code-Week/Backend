package handler

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"
	"strconv"
	"strings"
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user": {"IS_INVALID"}}))
		return
	}

	user := &dto.UserResponseDTO{}
	smapping.FillStruct(user, smapping.MapFields(userAuth))
	user.HasPassword = userAuth.Password != ""

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", user))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("AuthHandler.Login", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user": {"IS_INVALID"}}))
		return
	}

	userUpdate := &entity.User{}
	smapping.FillStruct(userUpdate, smapping.MapFields(userUpdateDTO))
	userUpdate.Phone = userUpdateDTO.Phone
	userUpdate.Major = userUpdateDTO.Major

	oldUser, err := h.userService.FindById(userAuth.ID)
	if err == nil {
		if oldUser.DataHasVerified {
			errs := make(map[string][]string)
			if userUpdateDTO.NIM != "" && (oldUser.NIM == nil || userUpdateDTO.NIM != *oldUser.NIM) {
				errs["nim"] = []string{"IS_ALREADY"}
			}
			if userUpdateDTO.Major != "" && userUpdateDTO.Major != oldUser.Major {
				errs["major"] = []string{"IS_ALREADY"}
			}
			if userUpdateDTO.Institusi != "" && userUpdateDTO.Institusi != oldUser.Institusi {
				errs["institusi"] = []string{"IS_ALREADY"}
			}
			if userUpdateDTO.SocMedDocument != "" && userUpdateDTO.SocMedDocument != oldUser.SocMedDocument {
				errs["socmed_document"] = []string{"IS_ALREADY"}
			}

			if len(errs) > 0 {
				logging.Low("UserHandler.UpdateMyProfile", "CONFLICT", "Data has been verified, cannot update restricted fields")
				c.JSON(http.StatusConflict, helper.CreateConflictResponse("Data sudah terverifikasi, tidak dapat mengubah field ini"))
				return
			}
		}

		var changes []string
		if userUpdateDTO.Name != "" && userUpdateDTO.Name != oldUser.Name { changes = append(changes, "nama") }
		if userUpdateDTO.NIM != "" && (oldUser.NIM == nil || userUpdateDTO.NIM != *oldUser.NIM) { changes = append(changes, "ktm/krs") }
		if userUpdateDTO.Phone != "" && userUpdateDTO.Phone != oldUser.Phone { changes = append(changes, "nomor telepon") }
		if userUpdateDTO.Major != "" && userUpdateDTO.Major != oldUser.Major { changes = append(changes, "jenjang") }
		if userUpdateDTO.Institusi != "" && userUpdateDTO.Institusi != oldUser.Institusi { changes = append(changes, "institusi") }
		if userUpdateDTO.SocMedDocument != "" && userUpdateDTO.SocMedDocument != oldUser.SocMedDocument { changes = append(changes, "dokumen sosial media") }

		if len(changes) > 0 {
			c.Set("audit_changes", strings.Join(changes, ", "))
		}
	}

	err = h.userService.Update(userUpdate, userAuth.ID)
	if err != nil {
		logging.High("AuthHandler.Login", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	user, err := h.userService.FindById(userAuth.ID)
	if err != nil {
		logging.High("AuthHandler.Login", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userResponse := &dto.UserResponseDTO{}
	smapping.FillStruct(userResponse, smapping.MapFields(user))
	userResponse.HasPassword = user.Password != ""

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", userResponse))
}

// @Summary Change Password
// @Tags Profile
// @Accept json
// @Produce  json
// @Param request body dto.ChangePasswordDTO true "Change Password"
// @Success 200 {object} helper.Response
// @Router /profile/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	changePasswordDTO := &dto.ChangePasswordDTO{}
	if err := c.Bind(changePasswordDTO); err != nil {
		logging.Low("UserHandler.ChangePassword", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	userAuth, ok := c.MustGet("user").(*entity.User)
	if !ok {
		logging.Low("UserHandler.ChangePassword", "BAD_REQUEST", "user not found in context")
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user": {"IS_INVALID"}}))
		return
	}

	err := h.userService.ChangePassword(userAuth.ID, changePasswordDTO.OldPassword, changePasswordDTO.NewPassword)
	if err != nil {
		logging.Low("UserHandler.ChangePassword", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("Data berhasil diperbarui", "Password berhasil diubah"))
}

func (h *UserHandler) GetAllUser(c *gin.Context) {
	startDateStr := c.Param("start_date")
	endDateStr := c.Param("end_date")
	countStr := c.Param("count")
	pageStr := c.Param("page")

	// Convert 'count' and 'page' to integers
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Invalid count parameter"))
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Invalid page parameter"))
		return
	}

	// Parse start_date and end_date into time.Time
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Invalid start_date format"))
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", "Invalid end_date format"))
		return
	}

	// Calculate the offset for pagination
	offset := (page - 1) * count

	// Fetch paginated data from the service
	users, totalUsers, err := h.userService.GetUsersByDateRange(startDate, endDate, count, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	// Prepare the response with mapped DTOs
	userResponses := []dto.UserResponseDTO{}
	for _, user := range users {
		var userResponse dto.UserResponseDTO
		err := smapping.FillStruct(&userResponse, smapping.MapFields(user))
		if err != nil {
			c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
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
	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Permintaan berhasil diproses", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user": {"IS_INVALID"}}))
		return
	}

	dataEvent, err := h.userService.GetEvents(userAuth.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"data": {"IS_INVALID"}}))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Data berhasil ditemukan", dataEvent))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	// Call service
	response, err := h.userService.AdminGetAllUsers(query)
	if err != nil {
		logging.High("UserHandler.AdminGetAllUsers", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Daftar user berhasil diambil", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user_id": {"IS_INVALID"}}))
		return
	}

	response, err := h.userService.AdminGetUserById(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("User tidak ditemukan"))
			return
		}
		logging.High("UserHandler.AdminGetUserById", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("User berhasil ditemukan", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user_id": {"IS_INVALID"}}))
		return
	}

	var updateData dto.AdminUpdateUserDTO
	if err := c.ShouldBindJSON(&updateData); err != nil {
		logging.Low("UserHandler.AdminUpdateUser", "BAD_REQUEST", err.Error())
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	response, err := h.userService.AdminUpdateUser(id, updateData)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("User tidak ditemukan"))
			return
		}
		logging.High("UserHandler.AdminUpdateUser", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.Set("target_email", response.Email)
	c.JSON(http.StatusCreated, helper.CreateMutationResponse("User berhasil diperbarui", response))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"user_id": {"IS_INVALID"}}))
		return
	}

	var deleteRequest dto.AdminDeleteUserDTO
	if err := c.ShouldBindJSON(&deleteRequest); err != nil {
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", map[string][]string{"alasan": {"IS_REQUIRED"}}))
		return
	}

	// Get user name for audit log before deletion
	user, err := h.userService.AdminGetUserById(id)
	if err == nil {
		c.Set("target_email", user.Email)
	}

	err = h.userService.AdminDeleteUser(id, deleteRequest.Alasan)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, helper.CreateNotFoundResponse("User tidak ditemukan"))
			return
		}
		logging.High("UserHandler.AdminDeleteUser", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.JSON(http.StatusCreated, helper.CreateMutationResponse("User berhasil dihapus", "User berhasil dihapus"))
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
		c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("Terdapat kesalahan pada permintaan", helper.FormatValidationError(err)))
		return
	}

	response, err := h.userService.AdminGetUserGrowthAnalytics(query)
	if err != nil {
		logging.High("UserHandler.AdminGetUserGrowthAnalytics", "INTERNAL_SERVER_ERROR", err.Error())
		c.JSON(http.StatusInternalServerError, helper.CreateInternalErrorResponse("Terdapat kesalahan pada permintaan"))
		return
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("Analitik pertumbuhan user berhasil diambil", response))
}
