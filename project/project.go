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
	PdID      string `json:"pd_id"`
	ProjDecsr string `json:"proj_decsr"`
	Drealiz   string `json:"drealiz"`
	ProjID    string `json:"proj_id"`
	PfName    string `json:"pf_name"`
	PfPath    string `json:"pf_path"`
	PfType    string `json:"pf_type"`
}

//Postprojects on BD
func Postprojects(c *gin.Context) {

	var (
		json Project
	)

	pool := c.PostForm("json")
	err := js.Unmarshal([]byte(pool), &json)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

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

	insertproject := fmt.Sprintf(`INSERT INTO public.tproject
	(proj_name, pd_id, proj_decsr, drealiz)
	VALUES('%s', %s, '%s', '%s');`, json.ProjName, json.PdID, json.ProjDecsr, json.Drealiz)

	_, err = dbConnect.Exec(insertproject)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insertproject err: %s", err.Error()))
	}

	todo := `SELECT max(proj_id) as proj_id
	FROM public.tproject;
	`

	rows, _ := dbConnect.Query(todo)

	defer rows.Close()
	print(rows)
	var projid string
	for rows.Next() {
		if err := rows.Scan(&projid); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.cd
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
		}

	}

	insertfile := fmt.Sprintf(`INSERT INTO public.tproject_file
	(proj_id, pf_name, pf_path, pf_type)
	VALUES(%s, '%s', '%s', %s);`, string(projid), filename, path, json.PfType)
	fmt.Printf(insertfile)
	_, err = dbConnect.Exec(insertfile)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insertfile err: %s", err.Error()))
	}

}

//UpdateProjects Projects
func UpdateProjects(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	files := form.File["file"]
	filepath := c.PostForm("filepath")
	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")
	newfullname := c.PostForm("new_fullname")
	pool := c.PostForm("json")

	var json Project

	err = js.Unmarshal([]byte(pool), &json)

	var path, filename string
	switch {
	case len(files) > 0:
		os.Mkdir(fmt.Sprintf("public/%s/%s", folder, subfolder), os.ModePerm)

		todo := fmt.Sprintf(`SELECT tproject.*, tproject_file.*
		FROM public.tproject tproject, public.tproject_file tproject_file
		WHERE 
			tproject_file.proj_id = tproject.proj_id and tproject.proj_id = %s order by tproject.drealiz desc ;`, json.ProjID)

		theCase := "lower"
		data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
		}

		os.Remove(strings.Replace(data[0]["pf_path"], "/file", "public", 1))

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

	case len(filepath) > 0 && len(newfullname) < 1 && len(files) == 0:

		path = filepath
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]

	}

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
		tproject_file.proj_id = tproject.proj_id and tproject.proj_id = %s order by tproject.drealiz desc;`, string(projid))

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
	todo := fmt.Sprintf(`
	SELECT tproject.*, tproject_file.*
		FROM public.tproject tproject, public.tproject_file tproject_file
		WHERE 
			tproject_file.proj_id = tproject.proj_id order by tproject.drealiz desc ;`)

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

		todo := fmt.Sprintf(`SELECT tproject.*, tproject_file.*
		FROM public.tproject tproject, public.tproject_file tproject_file
		WHERE 
			tproject_file.proj_id = tproject.proj_id and tproject.proj_id = %s order by tproject.drealiz desc ;`, id)

		theCase := "lower"
		data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
		}

		os.Remove(strings.Replace(data[0]["pf_path"], "/file", "public", 1))

		deletetProjectsfile := fmt.Sprintf(`
		DELETE FROM public.tproject_file
		WHERE proj_id=%s;`, id)

		deletetProjects := fmt.Sprintf(`DELETE FROM public.tproject
		WHERE proj_id=%s;`, id)

		_, err = dbConnect.Exec(deletetProjectsfile)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("delete project file err: %s", err.Error()))
		}
		_, err = dbConnect.Exec(deletetProjects)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("delete project err: %s", err.Error()))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})

}

//GetProjectsByID get project
func GetProjectsByID(c *gin.Context) {

	id := c.PostForm("pd_id")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT tproject.*, tproject_file.*
	FROM public.tproject tproject, public.tproject_file tproject_file
	WHERE 
		tproject_file.proj_id = tproject.proj_id and tproject.pd_id = %s order by tproject.drealiz desc ;`, id)

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

//GetProjectsLimit get project
func GetProjectsLimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT tproject.*, tproject_file.*
	FROM public.tproject tproject, public.tproject_file tproject_file
	WHERE 
		tproject_file.proj_id = tproject.proj_id order by tproject.drealiz desc limit %s offset %s;`, limit, offset)

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

//GetProjectsLimitCount get project
func GetProjectsLimitCount(c *gin.Context) {

	limit := c.PostForm("limit")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`select ceil(count(*)::real/%s::real) as pages_length from
	(SELECT *
	FROM public.tproject tproject, public.tproject_file tproject_file
	WHERE 
		tproject_file.proj_id = tproject.proj_id order by tproject.drealiz desc) a;`, limit)

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

//SearchInProjects Search In Projects
func SearchInProjects(c *gin.Context) {

	param := c.PostForm("param")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`
	select * from(
	   SELECT row_to_json(u.*)::text AS row_to_json
			   FROM ( SELECT tproject.proj_id,
				tproject.proj_name,
				tproject.pd_id,
				tproject.proj_decsr,
				tproject.drealiz,
				tproject_file.pf_id,
				tproject_file.proj_id,
				tproject_file.pf_name,
				tproject_file.pf_path,
				tproject_file.pf_type
			   FROM tproject tproject,
				tproject_file tproject_file) u) 
					  news where lower(news.row_to_json) like lower('%` + param + `%');`)

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
