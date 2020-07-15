package main

import (
	routes "ProtalMGTNIIP/routes"
	"log"
	_ "os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	// Init Router
	router := gin.Default()

	router.LoadHTMLGlob("template/*")
	// Route Handlers / Endpoints
	router.Use(cors.Default())
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
