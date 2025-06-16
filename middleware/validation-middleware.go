package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"

	"gcw/helper"
)

// ValidateDTO is a middleware that validates the request body against the provided DTO struct
// It uses the helper.ValidateStruct function to perform validation
func ValidateDTO(dto interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a new instance of the DTO type
		dtoType := reflect.TypeOf(dto).Elem()
		dtoValue := reflect.New(dtoType).Interface()

		// Bind the request body to the DTO
		if err := c.ShouldBindJSON(dtoValue); err != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("BAD_REQUEST", err.Error()))
			c.Abort()
			return
		}

		// Validate the DTO
		if errors := helper.ValidateStruct(dtoValue); errors != nil {
			c.JSON(http.StatusBadRequest, helper.CreateErrorResponse("VALIDATION_ERROR", errors))
			c.Abort()
			return
		}

		// Set the validated DTO in the context
		c.Set("dto", dtoValue)
		c.Next()
	}
}