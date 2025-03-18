package handler

import (
	"gcw/helper"
	"gcw/helper/logging"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type hackathonHandler struct {
}

type HackathonHandler interface {
	Submission(c *gin.Context)
}

func GateHackathonHandler() HackathonHandler {
	return &hackathonHandler{}
}

func (h *hackathonHandler) Submission(c *gin.Context) {

	destinationDir := "hackathon/nama_team"

	/*
		|--------------------------------------------------------------------------
		| Post Submission to AWS
		|--------------------------------------------------------------------------
	*/

	if err := c.Request.ParseMultipartForm(3 << 20); err != nil { // Max file size of 1MB
		logging.High("HackathonHandler.Submission", "Bad Request", "File Too large")
		response := helper.CreateErrorResponse("ERROR_PARSE_MULTIPART", "File Too large")
		c.JSON(http.StatusConflict, response)
		return
	}

	fileName := helper.UploadFile(c, "file", destinationDir)

	if fileName == "" {
		logging.High("HackathonHandler.Submission", "Internal Server Error", "Failed to upload file")
		c.JSON(http.StatusInternalServerError, helper.CreateErrorResponse("INTERNAL_SERVER_ERROR", "Failed to upload file"))
		return
	}

	fileURL := "https://notarius.s3.amazonaws.com/" + destinationDir + "/" + fileName

	log.Printf("Profile image uploaded: %s", fileURL)

	/*
		|--------------------------------------------------------------------------
		| Update Submission
		|--------------------------------------------------------------------------
	*/

	response := gin.H{
		"image_url": fileURL,
		"file_name": fileName,
		"directory": destinationDir,
	}

	c.JSON(http.StatusOK, helper.CreateSuccessResponse("success", response))

}
