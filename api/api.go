package api

import (
	config "ProtalMGTNIIP/config"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bdwilliams/go-jsonify/jsonify"
	"github.com/gin-gonic/gin"
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
	dbConnect.Close()
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
	dbConnect.Close()
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

func Upload(c *gin.Context) {

	type Filepaths struct {
		Filepath []string
	}

	var paths []string

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	files := form.File["files[]"]

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	print(subfolder, len(files))

	os.Mkdir(fmt.Sprintf("public/%s/", folder), os.ModePerm)

	for _, file := range files {

		if err := c.SaveUploadedFile(file, fmt.Sprintf("public/%s/%s", folder, file.Filename)); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		paths = append(paths, fmt.Sprintf("/file/%s/%s", folder, file.Filename))

	}
	filepath := Filepaths{
		Filepath: paths,
	}

	jsonData, err := json.Marshal(filepath)
	_ = jsonData
	if err != nil {
		log.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"filepath": string(jsonData)})

}

func Fileslist(c *gin.Context) {

	var files []string

	folder := c.PostForm("folder")
	subfolder := c.PostForm("subfolder")

	switch {
	case len(string(subfolder)) < 1:
		root := fmt.Sprintf("public/%s/%s/", folder, subfolder)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"filepath": files})
	case len(string(subfolder)) > 0:
		root := fmt.Sprintf("public/%s/", folder)
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{"filepath": files})
	}

}

func Newnews(c *gin.Context) {

	dbConnect := config.Connect()

	todo := "SELECT max(n_id) FROM public.tnews;"
	_ = todo
	rows, _ := dbConnect.Exec(todo)

	print(rows)
}
