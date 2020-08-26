package news

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

//Getprojects get news
func Getprojects(c *gin.Context) {
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
	dbConnect.Close()
	return
}

//Deleteprojects delete news by id
func Deleteprojects(c *gin.Context) {

	nids := c.PostFormArray("n_ids")

	dbConnect := config.Connect()
	for _, nid := range nids {
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
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
	dbConnect.Close()

}

//Project insert in BD
type Project struct {
	ProjID int    `json:"proj_id"`
	PfName string `json:"pf_name"`
	PfPath string `json:"pf_path"`
	PfType int    `json:"pf_type"`
}

//Postprojects on BD
func Postprojects(c *gin.Context) {

	var json Project

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["file"]

	filepath := c.PostForm("filepath")

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")
	var path, filename string
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
	case len(filepath) > 1:

		if err != nil {
			fmt.Printf("Invalid buffer size: %q\n", err)
			return
		}

		filepath = strings.Replace(filepath, "/file", "public", 1)
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
		destination := "public/photos/Проекты/" + filename
		err = Copy(filepath, destination)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}

		print(filename)
		path = strings.Replace(destination, "public", "/file", 1)
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := "SELECT max(pf_id) as pf_id FROM public.tproject_file;"

	insertproject := `INSERT INTO public.tproject_file
	(proj_id, pf_name, pf_path, pf_type)
	VALUES(?, '?', '?', ?);`

	_, err = dbConnect.Query(insertproject, json.ProjID, filename, path, json.PfType)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

	rows, _ := dbConnect.Query(todo)

	defer rows.Close()
	print(rows)
	var pfid string
	for rows.Next() {
		if err := rows.Scan(&pfid); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.cd
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}

	}
	insertphoto := fmt.Sprintf("INSERT INTO public.tnews_file (n_id, nf_name, nf_path, nf_type) VALUES(%s, '%s', '%s', 0);", pfid, filename, path)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}

//Updateprojects news
func Updateprojects(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	nid := c.PostForm("n_id")

	files := form.File["file"]
	filepath := c.PostForm("filepath")
	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")
	newfullname := c.PostForm("new_fullname")
	var path, filename string
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
		err := os.Rename(filepath, "public/photos/Новости/"+newfullname)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("rename file err: %s", err.Error()))
		}
		filename = newfullname

		path = "/file/photos/Новости/" + filename

	}

	date := c.PostForm("date")
	title := c.PostForm("title")
	text := c.PostForm("text")

	dbConnect := config.Connect()

	insertnews := fmt.Sprintf("UPDATE public.tnews SET n_date='%s', autor='', title='%s', textshort='', textfull='%s' WHERE n_id= %s;", date, title, text, nid)

	_, err = dbConnect.Exec(insertnews)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	deletetnews := fmt.Sprintf("DELETE FROM public.tnews_file WHERE n_id = %s;", nid)
	_, err = dbConnect.Exec(deletetnews)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("delete: %s", err.Error()))
	}
	insertphoto := fmt.Sprintf("INSERT INTO public.tnews_file (n_id, nf_name, nf_path, nf_type) VALUES(%s, '%s', '%s', 0);", nid, filename, path)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}
	dbConnect.Close()
}
