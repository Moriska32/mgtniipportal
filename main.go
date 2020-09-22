package main

import (
	routes "PortalMGTNIIP/routes"
	auth "PortalMGTNIIP/user"
	"log"
	_ "os"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	authMiddleware := auth.Auth()

	router.POST("/login", authMiddleware.LoginHandler)
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	auth := router.Group("/auth")
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)

	router.LoadHTMLGlob("template/*")
	router.Use(cors.Default())
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
