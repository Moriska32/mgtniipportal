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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://172.20.0.82:35766/", "http://172.20.0.82:5885/", "*"},
		AllowHeaders:     []string{"Access-Control-Request-Method", "Access-Control-Request-Headers", "X-Requested-With", "Access-Control-Allow-Headers", "Origin", "Authorization", "Content-Type", "Content-Length", "Accept", "Accept-Encoding", "X-HttpRequest"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"},
		ExposeHeaders:    []string{"Access-Control-Allow-Origin", "Content-Length"},
		AllowCredentials: true,
		MaxAge:           5600,
	}))
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
