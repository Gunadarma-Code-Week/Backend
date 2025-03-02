package main

import (
	"gcw/config"
	"gcw/router"

	"github.com/gin-gonic/gin"
)

var (
	database = config.SetupDatabaseConnection()
)

func main() {

	defer config.CloseDatabaseConnection(database)

	r := gin.Default()
	router.SetupRouter(r)

	r.Run(":8080")
}
