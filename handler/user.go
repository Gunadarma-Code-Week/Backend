package handler

import (
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/helper/logging"
	"gcw/service"
	"net/http"
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
