package user

import (
	config "PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elgs/gosqljson"

	"github.com/gin-gonic/gin"
)

//Newuser on BD
func Newuser(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["file"]
	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
	var path string

	for _, file := range files {

		if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		path = fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename)

	}

	login := c.PostForm("login")
	pass := c.PostForm("pass")
	fam := c.PostForm("fam")
	name := c.PostForm("name")
	otch := c.PostForm("otch")
	birthday := c.PostForm("birthday")
	foto := path
	hobby := c.PostForm("hobby")
	profskills := c.PostForm("profskills")
	drecrut := c.PostForm("drecrut")
	depid := c.PostForm("dep_id")
	chief := c.PostForm("chief")
	tel := c.PostForm("tel")

	userrole := c.PostForm("userrole")
	del := c.PostForm("del")

	dbConnect := config.Connect()

	insertuser := fmt.Sprintf("INSERT INTO public.tuser (login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, '%s', %s, %s);", login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, userrole, del)

	_, err = dbConnect.Exec(insertuser)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}

//Loginpass check login and pass
func Loginpass(c *gin.Context) {

	login := c.PostForm("login")
	pass := c.PostForm("pass")

	dbConnect := config.Connect()

	loginpass := fmt.Sprintf("SELECT user_id FROM public.tuser where login = %s AND pass = %s;", login, pass)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, loginpass)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

	return
}

//Deleteuser Delete users
func Deleteuser(c *gin.Context) {

	users := c.PostFormArray("user_id")
	for _, user := range users {
		dbConnect := config.Connect()

		deletetuser := fmt.Sprintf("DELETE FROM public.tuser WHERE user_id = %s;", user)

		_, err := dbConnect.Exec(deletetuser)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}

//Updateuser on BD
func Updateuser(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["file"]
	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
	var path string

	for _, file := range files {

		if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		path = fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename)

	}
	user := c.PostForm("user_id")
	login := c.PostForm("login")
	pass := c.PostForm("pass")
	fam := c.PostForm("fam")
	name := c.PostForm("name")
	otch := c.PostForm("otch")
	birthday := c.PostForm("birthday")
	foto := path
	hobby := c.PostForm("hobby")
	profskills := c.PostForm("profskills")
	drecrut := c.PostForm("drecrut")
	depid := c.PostForm("dep_id")
	chief := c.PostForm("chief")
	tel := c.PostForm("tel")
	userrole := c.PostForm("userrole")
	del := c.PostForm("del")

	dbConnect := config.Connect()

	insertuser := fmt.Sprintf("UPDATE public.tuser SET login='%s', pass='%s', fam='%s', name='%s', otch='%s', birthday='%s', foto='%s', hobby='%s', profskills='%s', drecrut='%s', dep_id=%s, chief=%s, tel='%s', userrole=%s, del=%s,  WHERE user_id=%s;", login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, userrole, del, user)

	_, err = dbConnect.Exec(insertuser)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}
