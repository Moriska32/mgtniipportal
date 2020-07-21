package api

import (
	config "ProtalMGTNIIP/config"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elgs/gosqljson"

	"github.com/gin-gonic/gin"
)

var print = fmt.Println

//OSReadDir List of Folders
func OSReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if !strings.Contains(file.Name(), ".") {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

//OSReadFile List of Files
func OSReadFile(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		if strings.Contains(file.Name(), ".") {
			files = append(files, strings.Replace(root, "public", "file", 1)+file.Name())
		}
	}
	return files, nil
}

//Dep list if dep by id
func Dep(c *gin.Context) {
	dbConnect := config.Connect()
	ID := c.Param("id")
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id = " + ID + ";"

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

//Deps List all of deps
func Deps(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT dep_id, name, parent_id FROM public.tdep;"

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

//Uploadtest test Load Files
func Uploadtest(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	name := c.PostForm("name")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename

	os.Mkdir(fmt.Sprintf("public/%s/", name), os.ModePerm)

	out, err := os.Create(fmt.Sprintf("public/%s/%s", name, filename))
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := fmt.Sprintf("/file/%s/%s", name, filename)
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

//Upload files on server
func Upload(c *gin.Context) {

	type Filepaths struct {
		Filepath []string
	}
	var filepath Filepaths
	_ = filepath
	var paths []string

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["files"]

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	print(subfolder, len(files))

	os.Mkdir(fmt.Sprintf("public/%s/", folder), os.ModePerm)

	switch {
	case len(subfolder) < 1:
		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s", folder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			paths = append(paths, fmt.Sprintf("/file/%s/%s", folder, file.Filename))

		}
		filepath = Filepaths{
			Filepath: paths,
		}
	case len(subfolder) > 0:
		os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
		for _, file := range files {

			if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
				return
			}

			paths = append(paths, fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename))

		}
		filepath = Filepaths{
			Filepath: paths,
		}
	}

	jsonData, err := json.Marshal(filepath)
	_ = jsonData
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{"filepath": string(jsonData)})

}

//Fileslist of files on server
func Fileslist(c *gin.Context) {

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	switch {
	case len(string(subfolder)) < 1:
		root := fmt.Sprintf("public/%s/", folder)

		dir, err := OSReadDir(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}

		files, err := OSReadFile(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
		c.JSON(http.StatusOK, gin.H{"subfolders": dir, "files": files})
	case len(string(subfolder)) > 0:

		root := fmt.Sprintf("public/%s/%s/", folder, subfolder)
		dir, err := OSReadDir(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}

		files, err := OSReadFile(root)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
		c.JSON(http.StatusOK, gin.H{"subfolders": dir, "files": files})
	}

}

//Getnews from bd
func Getnews(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file WHERE tnews_file.n_id = tnews.n_id;"

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

//Deletenews delete news by id
func Deletenews(c *gin.Context) {

	nid := c.PostForm("n_id")

	dbConnect := config.Connect()

	deletetnewsfile := fmt.Sprintf("DELETE FROM public.tnews_file WHERE n_id = %s;", nid)

	deletetnews := fmt.Sprintf("DELETE FROM public.tnews WHERE n_id = %s;", nid)

	_, err := dbConnect.Exec(deletetnewsfile)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	_, err = dbConnect.Exec(deletetnews)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}

//Postnews on BD
func Postnews(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["file"]

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
	var path, filename string

	for _, file := range files {

		if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file.Filename)); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		path = fmt.Sprintf("/file/%s/%s/%s", folder, subfolder, file.Filename)
		filename = file.Filename

	}

	date := c.PostForm("date")
	title := c.PostForm("title")
	text := c.PostForm("text")

	dbConnect := config.Connect()

	todo := "SELECT max(n_id) as n_id FROM public.tnews;"

	insertnews := fmt.Sprintf("INSERT INTO public.tnews (n_date, autor, title, textshort, textfull, dep_id) VALUES('%s', '', '%s', '', '%s', 1);", date, title, text)

	_, err = dbConnect.Exec(insertnews)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

	rows, _ := dbConnect.Query(todo)

	defer rows.Close()
	print(rows)
	var Nid string
	for rows.Next() {
		if err := rows.Scan(&Nid); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.cd
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
		print(Nid)
	}
	insertphoto := fmt.Sprintf("INSERT INTO public.tnews_file (n_id, nf_name, nf_path, nf_type) VALUES(%s, '%s', '%s', 0);", Nid, filename, path)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
}

//Mkrm Remove and make folder and files
func Mkrm(c *gin.Context) {
	doit := c.PostForm("doit")
	folder := c.PostForm("folder")
	subfolders := c.PostFormArray("subfolders")
	file := c.PostFormArray("file")
	switch {
	case doit == "rm":
		for _, subfolder := range subfolders {
			print(subfolder)
			err := os.Remove(fmt.Sprintf("public/%s/%s", folder, subfolder))
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			}
		}
	case doit == "mk":
		for _, subfolder := range subfolders {
			print(subfolder)
			err := os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			}
		}
	case doit == "rm" && len(file) > 0:
		for _, subfolder := range subfolders {
			err := os.Remove(fmt.Sprintf("public/%s/%s/%s", folder, subfolder, file))
			if err != nil {
				c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			}
		}
	}

}

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
	workplace := c.PostForm("workplace")
	userrole := c.PostForm("userrole")
	del := c.PostForm("del")
	postid := c.PostForm("post_id")

	dbConnect := config.Connect()

	insertuser := fmt.Sprintf("INSERT INTO public.tuser (login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, dep_id, chief, tel, workplace, userrole, del, post_id) VALUES('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', %s, %s, '%s', %s, %s, %s, %s);", login, pass, fam, name, otch, birthday, foto, hobby, profskills, drecrut, depid, chief, tel, workplace, userrole, del, postid)

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
