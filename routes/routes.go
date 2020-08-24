package routes

import (
	api "PortalMGTNIIP/api"
	files "PortalMGTNIIP/files"
	"PortalMGTNIIP/meetingroom"
	news "PortalMGTNIIP/news"
	projects "PortalMGTNIIP/project"
	user "PortalMGTNIIP/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Routes pool of func
func Routes(router *gin.Engine) {
	//GET
	router.GET("/", welcome)
	router.GET("/dep/:id", api.Dep)
	router.GET("/deps", api.Deps)
	router.GET("/orgstructure", api.Orgstructure)
	router.GET("/img", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_file.html", gin.H{})
	})

	router.GET("/post/:id", api.Post)
	router.GET("/posts", api.Posts)
	router.GET("/cbrdaily", api.Cbrdaily)
	router.GET("/weather", api.Weather)
	router.GET("/weathersss", api.Weathers)
	//Files
	router.POST("/upload", files.Upload)
	router.StaticFS("/file", http.Dir("public"))
	router.POST("/fileslist", files.Fileslist)
	router.POST("/rmfiles", files.Rmfiles)
	router.POST("/mkrmsubfolders", files.Mkrmsubfolders)
	//NEWS
	router.POST("/postnews", news.Postnews)
	router.POST("/getnewslist", news.Getnewslist)
	router.GET("/getnews", news.Getnews)
	router.POST("/newuser", user.Newuser)
	router.POST("/updatenews", news.Updatenews)
	router.POST("/deletenews", news.Deletenews)
	//User
	router.POST("/loginpass", user.Loginpass)
	router.GET("/getusers", user.Getusers)
	router.POST("/getuser", user.Getuser)
	router.POST("/deleteusers", user.Deleteuser)
	router.POST("/updateuser", user.Updateuser)
	//Object
	router.GET("/objectstype", api.Objectstype)
	//Project
	router.POST("/updateprojects", projects.Updateprojects)
	router.POST("/deleteprojects", projects.Deleteprojects)
	router.POST("/postprojects", projects.Postprojects)

	//Meetingroom
	router.GET("/meetingrooms", api.Meetingrooms)
	router.POST("/Newmeet", meetingroom.Newmeet)
}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}
