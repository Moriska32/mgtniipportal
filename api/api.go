package api

import (
	config "ProtalMGTNIIP/config"
	"encoding/json"
	"fmt"
	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var print = fmt.Println

type Filetypes struct {
	Filetype string `json: "filetype"`
}

func Dep(c *gin.Context) {
	dbConnect := config.Connect()
	ID := c.Param("id")
	todo := "SELECT dep_id, name, parent_id FROM public.tdep where parent_id = " + ID + ";"
	rows, err := dbConnect.Query(todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    jsonify.Jsonify(rows),
	})
	return

}
func Deps(c *gin.Context) {
	dbConnect := config.Connect()
	todo := "SELECT dep_id, name, parent_id FROM public.tdep;"
	rows, err := dbConnect.Query(todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Todo not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Single Todo",
		"data":    jsonify.Jsonify(rows),
	})
	return

}

//Load Files
func Upload(c *gin.Context) {
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
		log.Fatal("Create file : %s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	filepath := fmt.Sprintf("/file/%s/%s", name, filename)
	c.JSON(http.StatusOK, gin.H{"filepath": filepath})
}

func Uploadmany(c *gin.Context) {

	type Filepaths struct {
		Filepaths []string
	}

	var paths []string

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["files"]
	name := c.PostForm("name")

	os.Mkdir(fmt.Sprintf("public/%s/", name), os.ModePerm)

	for _, file := range files {
		filename := filepath.Base(file.Filename)
		out, err := os.Create(fmt.Sprintf("public/%s/%s", name, filename))
		if err != nil {
			log.Fatal("Create file : %s", err)
		}
		defer out.Close()

		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		paths := append(paths, fmt.Sprintf("/file/%s/%s", name, filename))
		_ = paths

	}

	filepaths := Filepaths{

		Filepaths: paths,
	}

	var jsonData []byte
	jsonData, err = json.Marshal(filepaths)
	if err != nil {
		log.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{"filepath": jsonData})
}
