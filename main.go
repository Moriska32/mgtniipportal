package main

import (
	routes "ProtalMGTNIIP/routes"
	"log"
	_ "os"

	"github.com/gin-gonic/gin"
)

func main() {

	// Init Router
	router := gin.Default()
	// Route Handlers / Endpoints
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
