package main

import (
	"context"
	"fmt"
	"gcw/config"
	"gcw/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"

	_ "gcw/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// @title GCW API
// @version 1.0
// @description This is a sample server for GCW API.

// @host localhost:8000
// @BasePath /api/v1/gcw/resources

func main() {
	database := config.SetupDatabaseConnection()
	defer config.CloseDatabaseConnection(database)

	file, err := createLogFile()

	if err != nil {
		log.Fatal("Error creating log file:", err)
	}

	defer file.Close()

	log.SetOutput(file)

	r := gin.Default()

	ginSwagger.URL("http://localhost:8000/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// CORS ORIGIN
	// r.Use(middleware.CORSMiddleware())

	origin := os.Getenv("CORS_ORIGIN")
	if origin == "" {
		origin = "http://localhost:3000"
	}
	origins := strings.Split(origin, ",")
	for i, o := range origins {
		origins[i] = strings.TrimSpace(o)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = origins
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "x-token", "cache-control", "Authorization"}
	corsConfig.AllowCredentials = true
	corsConfig.AllowMethods = []string{"POST", "DELETE", "GET", "PUT", "PATCH", "OPTIONS"}

	r.Use(cors.New(corsConfig))

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

// just for trigger action