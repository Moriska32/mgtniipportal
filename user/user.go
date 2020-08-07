package user

import (
	config "PortalMGTNIIP/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elgs/gosqljson"

	"github.com/gin-gonic/gin"
)

//Copy files
func Copy(sourceFile, destinationFile string) error {
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		fmt.Println(err)

	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		fmt.Println("Error creating", destinationFile)
		fmt.Println(err)

	}
	return err
}

//Newuser on BD
func Newuser(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["foto"]
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
	postid := c.PostForm("post_id")
	dbConnect := config.Connect()

	insertuser := fmt.Sprintf(`INSERT INTO public.tuser 
	(login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, dep_id, chief, tel, userrole, del, post_id) 
	VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, '%s', %s, %s, %s);
		`, login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, userrole, del, postid)

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

	loginpass := fmt.Sprintf("SELECT user_id FROM public.tuser where login = '%s' AND pass = '%s';", login, pass)

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

	users := c.PostFormArray("user_ids")
	for _, user := range users {
		dbConnect := config.Connect()
		print(user)
		deletetuser := fmt.Sprintf("DELETE FROM public.tuser WHERE user_id = %s;", user)

		_, err := dbConnect.Exec(deletetuser)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("delete user: %s", err.Error()))
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"user":   users,
	})

}

//Updateuser on BD
func Updateuser(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["foto"]
	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")
	newfullname := c.PostForm("new_fullname")
	filepath := c.PostForm("filepath")
	var path, filename string
	os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)

	switch {
	case len(files) > 0:
		os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)

		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			path = fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename)
			filename = file.Filename

		}
	case len(filepath) > 1 && len(newfullname) < 2:

		if err != nil {
			fmt.Printf("Invalid buffer size: %q\n", err)
			return
		}

		filepath = strings.Replace(filepath, "/file", "public", 1)
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
		destination := "public/photos/Новости/" + filename
		err = Copy(filepath, destination)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}

		print(filename)
		path = strings.Replace(destination, "public", "/file", 1)

	case len(newfullname) > 1:

		filepath = strings.Replace(filepath, "/file", "public", 1)
		err := os.Rename(filepath, "public/photos/Пользователи/"+newfullname)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("rename file err: %s", err.Error()))
		}
		filename = newfullname

		path = "/file/photos/Пользователи/" + filename

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
	postid := c.PostForm("post_id")
	del := c.PostForm("del")

	dbConnect := config.Connect()

	insertuser := fmt.Sprintf("UPDATE public.tuser SET login='%s', pass='%s', fam='%s', name='%s', otch='%s', birthday='%s', foto='%s', hobby='%s', profskills='%s', drecrut='%s', dep_id=%s, chief=%s, tel='%s', userrole=%s, del=%s, post_id = %s WHERE user_id=%s;", login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, userrole, del, postid, user)

	_, err = dbConnect.Exec(insertuser)

	if err != nil {
		print(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}

//Getusers get news
func Getusers(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT * FROM public.tuser;"

	defer dbConnect.Close()

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

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
	dbConnect.Close()
	return
}

//Getuser get news
func Getuser(c *gin.Context) {
	id := c.PostForm("user_id")
	dbConnect := config.Connect()
	todo := fmt.Sprintf("SELECT * from public.tuser where user_id = %s;", id)

	defer dbConnect.Close()

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

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
	dbConnect.Close()
	return
}
