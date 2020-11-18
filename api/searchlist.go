package api

import (
	"PortalMGTNIIP/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//Search in bd
func Search(c *gin.Context) {
	log.Println("поиск")

	param := c.PostForm("param")

	dbConnect := config.Connect()

	defer dbConnect.Close()

	log.Println(param)
	if param == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	refresh := `REFRESH MATERIALIZED VIEW public.searchlist;`

	_, err := dbConnect.Query(refresh)

	todo := `select pool from (
	SELECT row_to_json::text as pool
	FROM public.searchlist) as pool where lower(pool) like lower('%` + param + `%');`

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	dir := "documents"

	var resultf []string

	fileInfo, _, err := FilePathWalkDir("public/" + dir + "/")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range fileInfo {

		if strings.Contains(strings.ToLower(strings.Split(item, "\\")[len(strings.Split(item, "\\"))-1]), strings.ToLower(param)) {

			resultf = append(resultf, strings.Replace(item, "public", "file", 1))

		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
		"files":  resultf,
	})

	return

}

//FilePathWalkDir func
func FilePathWalkDir(root string) ([]string, []string, error) {
	var files []string
	var dirs []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		if info.IsDir() {
			dirs = append(files, path)
		}

		return nil
	})
	return files, dirs, err
}

//SearchInFolder files ind folder
func SearchInFolder(c *gin.Context) {

	dir := c.PostForm("dir")
	name := c.PostForm("name")
	_ = name

	var result []string

	fileInfo, _, err := FilePathWalkDir("public/" + dir + "/")
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range fileInfo {

		if strings.Contains(strings.ToLower(strings.Split(item, "\\")[len(strings.Split(item, "\\"))-1]), strings.ToLower(name)) {

			result = append(result, strings.Replace(item, "public", "file", 1))

		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   result,
	})

}
