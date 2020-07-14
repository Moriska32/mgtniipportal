package routes

import (
	api "ProtalMGTNIIP/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/dep/:id", api.Dep)

}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}
