package routes

import (
	api "ProtalMGTNIIP/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Routes pool of func
func Routes(router *gin.Engine) {
	router.GET("/", welcome)
	router.GET("/dep/:id", api.Dep)
	router.GET("/deps", api.Deps)
	router.GET("/img", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_file.html", gin.H{})
	})
	router.POST("/upload", api.Upload)
	router.StaticFS("/file", http.Dir("public"))
	router.POST("/fileslist", api.Fileslist)
	router.POST("/mkrm", api.Mkrm)
	router.POST("/postnews", api.Postnews)
	router.GET("/getnews", api.Getnews)
	router.POST("/newuser", api.Newuser)

}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}
