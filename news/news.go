package news

import (
	config "PortalMGTNIIP/config"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/elgs/gosqljson"

	"github.com/gin-gonic/gin"
)

//BUFFERSIZE buffer for file
var BUFFERSIZE int64

//Copy files
func Copy(src, dst string, BUFFERSIZE int64) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file.", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return fmt.Errorf("File %s already exists.", dst)
	}

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	if err != nil {
		panic(err)
	}

	buf := make([]byte, BUFFERSIZE)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return err
}

//Getnews get news
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
		BUFFERSIZE, err := strconv.ParseInt("1000", 10, 64)
		if err != nil {
			fmt.Printf("Invalid buffer size: %q\n", err)
			return
		}
		destination := "public/photos/Новости/"
		err = Copy(filepath, destination, BUFFERSIZE)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))]
		print(filename)
		path = destination + strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))]
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

	}
	insertphoto := fmt.Sprintf("INSERT INTO public.tnews_file (n_id, nf_name, nf_path, nf_type) VALUES(%s, '%s', '%s', 0);", Nid, filename, path)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
}

//Updatenews news
func Updatenews(c *gin.Context) {

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
		BUFFERSIZE, err := strconv.ParseInt("1000", 10, 64)
		if err != nil {
			fmt.Printf("Invalid buffer size: %q\n", err)
			return
		}
		destination := "public/photos/Новости/"
		filepath = strings.Replace(filepath, "file", "public", 1)
		err = Copy(filepath, destination, BUFFERSIZE)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))]
		print(filename)
		path = destination + strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))]
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

	insertphoto := fmt.Sprintf("UPDATE public.tnews_file SET nf_name='%s', nf_path='%s' WHERE nf_id= %s;", filename, path, nid)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
}
