package news

import (
	config "PortalMGTNIIP/config"
	js "encoding/json"
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

//Project insert in BD
type Project struct {
	ProjName  string `json:"proj_name"`
	PdID      int    `json:"pd_id"`
	ProjDecsr string `json:"proj_decsr"`
	Drealiz   string `json:"drealiz"`
	ProjID    int    `json:"proj_id"`
	PfName    string `json:"pf_name"`
	PfPath    string `json:"pf_path"`
	PfType    int    `json:"pf_type"`
}

//Postprojects on BD
func Postprojects(c *gin.Context) {

	var (
		json Project
	)

	pool := c.PostForm("json")
	err := js.Unmarshal([]byte(pool), &json)
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

	insertproject := `INSERT INTO public.tproject
	(proj_name, pd_id, proj_decsr, drealiz)
	VALUES('%s', %s, '%s', '%s');`

	_, err = dbConnect.Query(insertproject, json.ProjName, json.PdID, json.ProjDecsr, json.Drealiz)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

	insertfile := fmt.Sprintf(`INSERT INTO public.tproject_file
	(proj_id, pf_name, pf_path, pf_type)
	VALUES(%s, '%s', '%s', %s);`, string(json.ProjID), json.PfName, path, string(json.PfType))

	_, err = dbConnect.Exec(insertfile)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}

//UpdateProjects Projects
func UpdateProjects(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

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
		destination := "public/photos/Проекты/" + filename
		err = Copy(filepath, destination)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}

		print(filename)
		path = strings.Replace(destination, "public", "/file", 1)

	case len(newfullname) > 1:

		filepath = strings.Replace(filepath, "/file", "public", 1)
		err := os.Rename(filepath, "public/photos/Проекты/"+newfullname)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("rename file err: %s", err.Error()))
		}
		filename = newfullname

		path = "/file/photos/Проекты/" + filename

	}

	var json Project

	pool := c.PostForm("json")
	err = js.Unmarshal([]byte(pool), &json)

	dbConnect := config.Connect()

	insertnews := fmt.Sprintf(`UPDATE public.tproject
	SET proj_name='%s', pd_id=%s, proj_decsr='%s', drealiz='%s'
	WHERE proj_id=%s;`, json.ProjName, string(json.PdID), json.ProjDecsr, json.Drealiz, string(json.ProjID))

	_, err = dbConnect.Exec(insertnews)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	deletetnews := fmt.Sprintf(`DELETE FROM public.tproject_file
	WHERE proj_id=%s;`, string(json.ProjID))
	_, err = dbConnect.Exec(deletetnews)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("delete: %s", err.Error()))
	}

	insertfile := fmt.Sprintf(`INSERT INTO public.tproject_file
	(proj_id, pf_name, pf_path, pf_type)
	VALUES(%s, '%s', '%s', %s);`, string(json.ProjID), filename, path, string(json.PfType))

	_, err = dbConnect.Exec(insertfile)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}
	dbConnect.Close()
}

//GetProjectsDirection Get projects Direction
func GetProjectsDirection(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := `SELECT pd_id, pd_name
	FROM public.tproject_direction;
	;`

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

//GetProject get project
func GetProject(c *gin.Context) {

	projid := c.PostForm("proj_id")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT tproject.*, tproject_file.*
	FROM public.tproject tproject, public.tproject_file tproject_file
	WHERE 
		tproject_file.proj_id = tproject.proj_id and tproject.proj_id = %s;`, string(projid))

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

//GetProjects get project
func GetProjects(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT tproject.*, tproject_file.*
	FROM public.tproject tproject, public.tproject_file tproject_file
	WHERE 
		tproject_file.proj_id = tproject.proj_id;`)

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

//DeleteProjects delete Projects by id
func DeleteProjects(c *gin.Context) {

	projids := c.PostFormArray("proj_ids")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	for _, id := range projids {
		deletetProjectsfile := fmt.Sprintf(`
		DELETE FROM public.tproject_file
		WHERE proj_id=%s;
		"`, id)

		deletetProjects := fmt.Sprintf(`DELETE FROM public.tproject
		WHERE proj_id=%s;`, id)

		_, err := dbConnect.Exec(deletetProjectsfile)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
		_, err = dbConnect.Exec(deletetProjects)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}
