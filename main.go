package main

import (
	"fmt"
	"gcw/config"
	"gcw/router"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	database = config.SetupDatabaseConnection()
)

func createLogFile() (*os.File, error) {
	err := os.MkdirAll("log", os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create log directory: %w", err)
	}

	currentTime := time.Now()
	logFilename := fmt.Sprintf("log/%d-%02d.log", currentTime.Year(), currentTime.Month())

	file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %w", err)
	}

	return file, nil
}

func main() {

	defer config.CloseDatabaseConnection(database)

	file, err := createLogFile()

	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	defer file.Close()

	log.SetOutput(file)

	r := gin.Default()
	router.SetupRouter(r)

	r.Run(":8080")
}
