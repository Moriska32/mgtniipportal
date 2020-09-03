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
	router.GET("/chief", api.Chief)
	router.GET("/orgstructure", api.Orgstructure)
	router.GET("/img", func(c *gin.Context) {
		c.HTML(http.StatusOK, "select_file.html", gin.H{})
	})

	router.GET("/post/:id", api.Post)
	router.GET("/posts", api.Posts)
	router.GET("/cbrdaily", api.Cbrdaily)
	router.GET("/weather", api.Weathers)
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
	router.GET("/getportalusers", user.GetUsersNotPass)
	router.POST("/getuser", user.Getuser)
	router.POST("/getportaluser", user.GetuserNotPass)

	router.POST("/deleteusers", user.Deleteuser)
	router.POST("/updateuser", user.Updateuser)
	//Object
	router.GET("/objectstype", api.Objectstype)
	router.GET("/objects/:id", api.Objects)
	//Project
	router.POST("/postproject", projects.Postprojects)
	router.POST("/updateproject", projects.UpdateProjects)
	router.POST("/getproject", projects.GetProject)
	router.GET("/getprojects", projects.GetProjects)
	router.GET("/getprojectsdirections", projects.GetProjectsDirection)
	router.POST("/deleteprojects", projects.DeleteProjects)

	//Meetingroom
	router.GET("/meetingrooms", api.Meetingrooms)
	router.POST("/newmeet", meetingroom.Newmeet)
	router.POST("/getmeets", meetingroom.Getmeets)
	router.POST("/deletemeet", meetingroom.Deletemeet)
	router.POST("/updatemeet", meetingroom.Updatemeet)
	router.GET("/getallmeets", meetingroom.GetAllMeets)

	//Mail sender
	router.POST("/sendmail", api.SendMail)
	router.POST("/sendrequest", api.SendRequest)
	router.POST("/getrequest", api.GetRequest)
	//router.POST("/sendmailit", api.SendMailIT)

	//HH
	router.POST("/posthh", api.PostHH)
	router.POST("/updatehh", api.UpdateHH)
	router.POST("/deletehh", api.DeleteHH)
	router.GET("/gethhs", api.GetHHs)
	router.POST("/gethh", api.GetHH)

}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}
