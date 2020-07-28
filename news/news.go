package news

import (
	config "PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elgs/gosqljson"

	"github.com/gin-gonic/gin"
)

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
