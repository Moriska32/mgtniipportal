package routes

import (
	"PortalMGTNIIP/api"
	files "PortalMGTNIIP/files"
	"PortalMGTNIIP/geomap"
	"PortalMGTNIIP/meetingroom"
	news "PortalMGTNIIP/news"
	projects "PortalMGTNIIP/project"
	"PortalMGTNIIP/training"
	user "PortalMGTNIIP/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

//Routes pool of func
func Routes(router *gin.Engine) {
	router.StaticFS("/file", http.Dir("public"))

	authMiddleware := user.Auth()

	auth := router.Group("/v1")
	auth.Use(authMiddleware.MiddlewareFunc())
	auth.Use(user.Blacklist)
	{
		root := auth.Group("/api")
		{
			//TOKEN
			root.GET("/logout", user.Logout)
			root.GET("/tokendescr", user.Token)

			//GET

			root.GET("/dep/:id", api.Dep)
			root.GET("/deps", api.Deps)
			root.GET("/chiefs", api.Chief)
			root.GET("/orgstructure", api.Orgstructure)
			root.GET("/img", func(c *gin.Context) {
				c.HTML(http.StatusOK, "select_file.html", gin.H{})
			})

			root.GET("/post/:id", api.Post)
			root.GET("/posts", api.Posts)
			root.GET("/cbrdaily", api.Cbrdaily)
			root.GET("/weather", api.Weathers)
			root.GET("/weathersss", api.Weathers)
			//Files
			root.POST("/upload", files.Upload)

			root.POST("/fileslist", files.Fileslist)
			root.POST("/rmfiles", files.Rmfiles)
			root.POST("/mkrmsubfolders", files.Mkrmsubfolders)
			//NEWS
			root.POST("/postnews", news.Postnews)
			root.POST("/getnewslist", news.Getnewslist)
			root.GET("/getnews", news.Getnews)
			root.POST("/updatenews", news.Updatenews)
			root.POST("/deletenews", news.Deletenews)
			root.POST("/getonenews", news.GetOneNews)
			root.POST("/getnewslimit", news.GetnewsLimit)
			root.POST("/getnewslimitcount", news.GetnewsLimitCount)
			root.GET("/newsthemes", news.Newsthemes)
			root.POST("/getnewsbytheme", news.GetNewsByTheme)
			root.POST("/getnewsbytime", news.GetNewsByTime)

			root.POST("/searchinnews", news.SearchInNews)

			//User
			root.POST("/newuser", user.Newuser)
			root.POST("/loginpass", user.Loginpass)
			root.GET("/getusers", user.Getusers)
			root.GET("/getportalusers", user.GetUsersNotPass)
			root.POST("/getuser", user.Getuser)
			root.POST("/getportaluser", user.GetuserNotPass)
			root.GET("/getsuperusers", user.Getsuperuser)
			root.POST("/deleteusers", user.Deleteuser)
			root.POST("/updateuser", user.Updateuser)
			root.POST("/getuserslimit", user.Getuserslimit)
			root.POST("/getusersobj", user.Getusersobj)
			root.POST("/getusersletter", user.Getusersletter)
			root.POST("/getusersbyobj", user.Getusersbyobj)
			root.POST("/getuserslimitcount", user.Getuserslimitcount)
			root.GET("/getusersadmins", user.Getusersadmins)
			root.GET("/getusersletters", user.Getusersletters)
			root.GET("/getuserstime", user.Getuserstime)
			root.POST("/searchinusers", user.SearchInUsers)
			root.POST("/updatepass", user.UpdatePass)
			root.POST("/updatephoto", user.UpdatePhoto)

			//Object
			root.GET("/objectstype", api.Objectstype)
			root.GET("/objects/:id", api.Objects)

			//Project
			root.POST("/postproject", projects.Postprojects)
			root.POST("/updateproject", projects.UpdateProjects)
			root.POST("/getproject", projects.GetProject)
			root.GET("/getprojects", projects.GetProjects)
			root.GET("/getprojectsdirections", projects.GetProjectsDirection)
			root.POST("/deleteprojects", projects.DeleteProjects)
			root.POST("/getprojectsbydirectionid", projects.GetProjectsByID)
			root.POST("/getprojectslimit", projects.GetProjectsLimit)
			root.POST("/getprojectslimitcount", projects.GetProjectsLimitCount)
			root.POST("/searchinprojects", projects.SearchInProjects)

			//Meetingroom
			root.GET("/meetingrooms", api.Meetingrooms)
			root.POST("/newmeet", meetingroom.Newmeet)
			root.POST("/getmeets", meetingroom.Getmeets)
			root.POST("/deletemeet", meetingroom.Deletemeet)
			root.POST("/updatemeet", meetingroom.Updatemeet)
			root.GET("/getallmeets", meetingroom.GetAllMeets)

			root.POST("/getmeetslimit", meetingroom.GetMeetsLimit)

			//Mail sender
			root.POST("/sendmail", api.SendMail)
			root.POST("/sendrequest", api.SendRequest)
			root.POST("/getrequests", api.GetRequest)
			//root.POST("/sendmailit", api.SendMailIT)

			root.POST("/getrequestlimit", api.GetRequestLimit)

			//HH
			root.POST("/posthh", api.PostHH)
			root.POST("/updatehh", api.UpdateHH)
			root.POST("/deletehh", api.DeleteHH)
			root.GET("/gethhs", api.GetHHs)
			root.POST("/gethh", api.GetHH)

			//Map
			root.POST("/geombyfloor", geomap.Map)

			//search
			root.POST("/search", api.Search)
			root.POST("/searchinfolder", api.SearchInFolder)

			//TrainingTopic
			root.POST("/posttrainingtopic", training.Posttrainingtopic)
			root.POST("/updatetrainingtopic", training.Updatetrainingtopic)
			root.POST("/gettrainingtopic", training.Gettrainingtopic)
			root.POST("/gettrainingstopicslimit", training.Gettrainingstopicslimit)
			root.POST("/deletetrainingtopic", training.Deletetrainingtopic)

			root.POST("/posttraining", training.Posttraining)
			root.POST("/updatetraining", training.Updatetraining)
			root.POST("/deletetrainings", training.Deletetrainings)
			root.POST("/gettraining", training.Gettraining)
			root.POST("/gettrainingslimit", training.Gettrainingslimit)

			root.GET("/gettrainingstopicstypes", training.Gettrainingstopicstypes)
			root.GET("/getactivetrainings", training.Getactivetrainings)
			root.GET("/getpasttrainings", training.Getpasttrainings)

			root.GET("/gettrainingreqstatuses", training.Gettrainingreqstatuses)

			//TrainingAnakitic
			root.GET("/getpooltrainingbyyear", training.Getpooltrainingbyyear)
			root.GET("/getpoolusersbydep", training.Getpoolusersbydep)
			root.GET("/getexcelanaliticstraining", training.GetExelAnaliticsTraining)

			//TrainingRequest

			root.POST("/posttrainingrequest", training.PostTrainingRequest)
			root.POST("/updatetrainingrequest", training.UpdateTrainingRequest)
			root.POST("/gettrainingsrequestslimit", training.GetTrainingRequestsLimit)
			root.GET("/getuserwithtrainingsandrequests", training.GetUserWithTrainingsAndRequests)

			//root.GET("/test", user.GetTokenInfo)

		}
	}

}

func welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "Welcome To API",
	})
	return
}
