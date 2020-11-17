package user

import (
	config "PortalMGTNIIP/config"
	fl "PortalMGTNIIP/files"
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
	if len(folder)+len(subfolder) == 0 {

		path = "/file/photos/Пользователи/default-user-avatar.jpg"

	} else {
		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			path = fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename)
			fl.Resize(fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename))

		}
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
	workplace := c.PostForm("workplace")
	userrole := c.PostForm("userrole")
	del := c.PostForm("del")
	postid := c.PostForm("post_id")
	dbConnect := config.Connect()

	insertuser := fmt.Sprintf(`INSERT INTO public.tuser 
	(login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id) 
	VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, '%s',%s, %s, %s, %s);
		`, login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, workplace, userrole, del, postid)

	_, err = dbConnect.Exec(insertuser)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	dbConnect.Close()
}

//Loginpass check login and pass
func Loginpass(c *gin.Context) {

	login := c.PostForm("login")
	pass := c.PostForm("pass")
	superadminarray := make([]map[string]int, 0)
	superadmin := make(map[string]int)
	superadmin["user_id"] = 1

	superadminarray = append(superadminarray, superadmin)

	if login == "superadmin" && pass == "superadmin12345" {

		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"data":   superadminarray,
		})

		return
	}

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
	dbConnect.Close()
	return
}

//Deleteuser Delete users
func Deleteuser(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	users := c.PostFormArray("user_ids")
	for _, user := range users {

		todo := fmt.Sprintf("select foto FROM public.tuser WHERE user_id = %s;", user)

		theCase := "lower"
		data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
		}
		if data[0]["foto"] != "/file/photos/Пользователи/default-user-avatar.jpg" {
			err = os.Remove(strings.Replace(data[0]["foto"], "/file", "public", 1))
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("Can't delete file: %s", err.Error()))

			}
		}
		deletetuser := fmt.Sprintf("DELETE FROM public.tuser WHERE user_id = %s;", user)

		_, err = dbConnect.Exec(deletetuser)
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
	user := c.PostForm("user_id")
	var path, filename string
	os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)

	dbConnect := config.Connect()
	defer dbConnect.Close()

	switch {

	case (len(folder) + len(subfolder) + len(filepath)) == 0:

		path = "/file/photos/Пользователи/default-user-avatar.jpg"

	case len(files) > 0:
		os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)

		todo := fmt.Sprintf("select foto FROM public.tuser WHERE user_id = %s;", user)

		theCase := "lower"
		data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
		}
		if data[0]["foto"] != "/file/photos/Пользователи/default-user-avatar.jpg" {
			err = os.Remove(strings.Replace(data[0]["foto"], "/file", "public", 1))
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("Can't delete file: %s", err.Error()))

			}
		}

		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			path = fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename)
			fl.Resize(fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename))
			filename = file.Filename

		}
	case len(filepath) > 1 && len(newfullname) < 2:

		if err != nil {
			fmt.Printf("Invalid buffer size: %q\n", err)
			return
		}

		filepath = strings.Replace(filepath, "/file", "public", 1)
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
		destination := "public/photos/Пользователи/" + filename
		fl.Resize(fmt.Sprintf(destination))
		err = Copy(filepath, destination)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}

		print(filename)
		path = strings.Replace(destination, "public", "/file", 1)

	case len(newfullname) > 1:

		filepath = strings.Replace(filepath, "/file", "public", 1)
		err := os.Rename(filepath, "public/photos/Пользователи/"+newfullname)

		fl.Resize("public/photos/Пользователи/" + newfullname)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("rename file err: %s", err.Error()))
		}
		filename = newfullname

		path = "/file/photos/Пользователи/" + filename

	case len(filepath) > 0 && len(newfullname) < 1 && len(files) == 0:

		path = filepath
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]

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
	workplace := c.PostForm("workplace")
	userrole := c.PostForm("userrole")
	postid := c.PostForm("post_id")
	del := c.PostForm("del")

	insertuser := fmt.Sprintf("UPDATE public.tuser SET login='%s', pass='%s', fam='%s', name='%s', otch='%s', birthday='%s', foto='%s', hobby='%s', profskills='%s', drecrut='%s', dep_id=%s, chief=%s, tel='%s', workplace = %s, userrole=%s, del=%s, post_id = %s WHERE user_id = %s ;", login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, workplace, userrole, del, postid, user)

	_, err = dbConnect.Exec(insertuser)

	if err != nil {
		print(err)
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}

//Getusers get news
func Getusers(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT * FROM public.tuser where login not in ('admin', 'moder', 'user');"

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
	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

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
	var data []map[string]string
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
	data[0]["foto_min"] = strings.Replace(strings.Replace(data[0]["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//GetuserNotPass get Get user Not Pass
func GetuserNotPass(c *gin.Context) {
	id := c.PostForm("user_id")
	dbConnect := config.Connect()
	todo := fmt.Sprintf(`SELECT user_id, login, fam, "name", otch, birthday, foto, 
	hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del,
	 post_id from public.tuser where user_id = %s and login not in ('admin', 'moder', 'user');`, id)

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

	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//GetUsersNotPass get news
func GetUsersNotPass(c *gin.Context) {
	dbConnect := config.Connect()
	todo := `SELECT user_id, login, fam, "name", otch, birthday, foto, 
	hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id FROM public.tuser where login not in ('admin', 'moder', 'user');`

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

	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

	return
}

//Getsuperuser 'admin', 'moder', 'user'
func Getsuperuser(c *gin.Context) {

	dbConnect := config.Connect()
	todo := `SELECT user_id, login, fam, "name", otch, birthday, foto, 
	hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id FROM public.tuser where login in ('admin', 'moder', 'user');`

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

	return

}

//Getuserslimit get user by limit
func Getuserslimit(c *gin.Context) {
	dbConnect := config.Connect()

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	todo := fmt.Sprintf(`SELECT user_id, login, fam, "name", otch, birthday, foto, 
	hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id
	FROM public.tuser where login not in ('admin', 'user', 'moder') limit %s offset %s;`, limit, offset)

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
	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//Getusersobj get user by limit
func Getusersobj(c *gin.Context) {
	dbConnect := config.Connect()

	objid := c.PostForm("obj_id")

	todo := fmt.Sprintf(`
	SELECT t.user_id, t.login, t.fam, t."name", t.otch, t.birthday, t.foto, t.hobby, t.profskills, t.drecrut, t.dep_id, t.chief, t.tel, t.workplace, t.userrole, t.del, t.post_id
		FROM public.tobject floor, public.tobject room,public.tobject cabinet, public.tuser t
		where floor.object_id = %s and floor.object_id = room.container_id and room.object_id = cabinet.container_id and t.workplace = cabinet.object_id;`, objid)

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
	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//Getusersletter get user by letter
func Getusersletter(c *gin.Context) {
	dbConnect := config.Connect()

	letter := c.PostForm("letter")

	todo := fmt.Sprintf(`SELECT user_id, login, fam, "name", otch, birthday, foto,
	 hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id
	FROM public.tuser where upper(substr(fam,1,1)) = upper('%s') and login not in ('admin', 'user', 'moder');`, letter)

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
	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//Getusersbyobj get user by limit
func Getusersbyobj(c *gin.Context) {
	dbConnect := config.Connect()

	objid := c.PostForm("obj_id")

	todo := fmt.Sprintf(`SELECT user_id, login, fam, "name", otch, birthday, foto, hobby,
	profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id
	FROM public.tuser where workplace = %s;`, objid)

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
	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//Getuserslimitcount get count of users by limit
func Getuserslimitcount(c *gin.Context) {
	dbConnect := config.Connect()

	limit := c.PostForm("limit")

	todo := fmt.Sprintf(`SELECT ceil(count(*)::real/%s::real) as pages_length from public.tuser where login not in ('admin', 'user', 'moder');`, limit)

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

//Getusersadmins get count of users by limit
func Getusersadmins(c *gin.Context) {
	dbConnect := config.Connect()

	todo := `select user_id, login, fam, "name", otch, birthday, foto, hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id from tuser
	where user_id in ( 72,
253,
272,
259,
308,
303,
380,
1,
261,
100,
194,
321,
309,
415,
262);`

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

	for _, items := range data {

		items["foto_min"] = strings.Replace(strings.Replace(items["foto"], ".jpg", "-min.jpg", 1), "Пользователи", "Пользователи-min", 1)

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//Getusersletters get user by letter
func Getusersletters(c *gin.Context) {
	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`SELECT string_agg(distinct(substr(fam,1,1)), ',') as letter
	FROM public.tuser where fam ~ '[а-я]+' order by letter;`)

	theCase := "lower"
	rus, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	todo = fmt.Sprintf(`SELECT string_agg(distinct(substr(fam,1,1)), ',') as letter
	FROM public.tuser where fam ~ '[a-z]+' and login not in ('admin', 'user', 'moder') order by letter;`)

	eng, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"rus":    rus[0]["letter"],
		"eng":    eng[0]["letter"],
	})
	dbConnect.Close()
	return
}
