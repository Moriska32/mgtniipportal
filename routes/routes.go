package routes

import (
	api "PortalMGTNIIP/api"
	files "PortalMGTNIIP/files"
	news "PortalMGTNIIP/news"
	user "PortalMGTNIIP/user"
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
	router.POST("/upload", files.Upload)
	router.StaticFS("/file", http.Dir("public"))
	router.POST("/fileslist", files.Fileslist)
	router.POST("/rmfiles", files.Rmfiles)
	router.POST("/mkrmsubfolders", files.Mkrmsubfolders)
	router.POST("/postnews", news.Postnews)
	router.GET("/getnews", news.Getnews)
	router.POST("/newuser", user.Newuser)
	router.POST("/updatenews", news.Updatenews)
	router.POST("/deletenews", news.Deletenews)
	router.POST("/loginpass", user.Loginpass)

}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}
