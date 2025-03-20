package main

import (
	"context"
	"fmt"
	"gcw/config"
	"gcw/middleware"
	"gcw/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
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
	// if err := godotenv.Load(); err != nil {
	// 	panic("Failed to load env file")
	// }

	database := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(database)

	file, err := createLogFile()

	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	defer file.Close()

	log.SetOutput(file)

	r := gin.Default()

	// CORS ORIGIN
	r.Use(middleware.CORSMiddleware())

	router.SetupRouter(r)

	// GRACEFULL SHUTDOWN
	srv := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
