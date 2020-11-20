package news

import (
	config "PortalMGTNIIP/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
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

	for _, item := range data {

		if item["nf_type"] == "1" {

			item["screen"] = PictureFromVideo(item["nf_path"])
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
	dbConnect.Close()
	return
}

//Deletenews delete news by id
func Deletenews(c *gin.Context) {

	nids := c.PostFormArray("n_ids")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	for _, nid := range nids {

		todo := fmt.Sprintf("SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file WHERE tnews_file.n_id = tnews.n_id AND tnews_file.n_id = %s order by tnews.n_date desc;", nid)

		theCase := "lower"
		data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
		}

		err = os.Remove(strings.Replace(data[0]["nf_path"], "/file", "public", 1))
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't delete file: %s", err.Error()))
		}
		deletetnewsfile := fmt.Sprintf("DELETE FROM public.tnews_file WHERE n_id = %s;", nid)

		deletetnews := fmt.Sprintf("DELETE FROM public.tnews WHERE n_id = %s;", nid)

		_, err = dbConnect.Exec(deletetnewsfile)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Delete file from BD err: %s", err.Error()))
		}
		_, err = dbConnect.Exec(deletetnews)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Delete news err: %s", err.Error()))
		}
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
	}

	date := c.PostForm("date")
	title := c.PostForm("title")
	text := c.PostForm("text")
	theme := c.PostForm("theme")

	dbConnect := config.Connect()

	todo := "SELECT max(n_id) as n_id FROM public.tnews;"

	insertnews := fmt.Sprintf("INSERT INTO public.tnews (n_date, autor, title, textshort, textfull, dep_id, theme) VALUES('%s', '', '%s', '', '%s', 1, %s);", date, title, text, theme)

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

	nftype := c.PostForm("nf_type")

	insertphoto := fmt.Sprintf("INSERT INTO public.tnews_file (n_id, nf_name, nf_path, nf_type) VALUES(%s, '%s', '%s', %s);", Nid, filename, path, nftype)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}
	dbConnect.Close()
}

//Updatenews news
func Updatenews(c *gin.Context) {

	form, err := c.MultipartForm()

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

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

		todo := fmt.Sprintf("SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file WHERE tnews_file.n_id = tnews.n_id AND tnews_file.n_id = %s order by tnews.n_date desc;", nid)

		theCase := "lower"
		data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
		}

		err = os.Remove(strings.Replace(data[0]["nf_path"], "/file", "public", 1))
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Can't delete file: %s", err.Error()))
		}

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
		destination := "public/" + folder + "/" + subfolder + "/" + filename
		err = Copy(filepath, destination)
		if err != nil {
			fmt.Printf("File copying failed: %q\n", err)
		}

		print(filename)
		path = strings.Replace(destination, "public", "/file", 1)

	case len(newfullname) > 1:

		filepath = strings.Replace(filepath, "/file", "public", 1)
		err := os.Rename(filepath, "public/"+folder+"/"+subfolder+"/"+newfullname)

		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("rename file err: %s", err.Error()))
		}
		filename = newfullname

		path = "file/" + folder + "/" + subfolder + "/" + newfullname

	case len(filepath) > 0 && len(newfullname) < 1 && len(files) == 0:

		path = filepath
		filename = strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
	}

	date := c.PostForm("date")
	title := c.PostForm("title")
	text := c.PostForm("text")

	log.Println(nid)

	nftype := c.PostForm("nf_type")
	insertphoto := fmt.Sprintf("UPDATE public.tnews_file SET nf_name='%s', nf_path='%s', nf_type=%s WHERE n_id= %s;", filename, path, nftype, nid)

	_, err = dbConnect.Exec(insertphoto)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	insertnews := fmt.Sprintf("UPDATE public.tnews SET n_date='%s', autor='', title='%s', textshort='', textfull='%s' WHERE n_id= %s;", date, title, text, nid)

	_, err = dbConnect.Exec(insertnews)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
	}

}

//Getnewslist get news
func Getnewslist(c *gin.Context) {

	nftype := c.PostForm("nf_type")

	dbConnect := config.Connect()
	todo := fmt.Sprintf("SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file WHERE tnews_file.n_id = tnews.n_id AND tnews_file.nf_type = %s order by tnews.n_date desc;", nftype)

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

//GetOneNews get one news
func GetOneNews(c *gin.Context) {

	id := c.PostForm("id")

	dbConnect := config.Connect()
	todo := fmt.Sprintf("SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file WHERE tnews_file.n_id = tnews.n_id AND tnews_file.n_id = %s order by tnews.n_date desc;", id)

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

//PictureFromVideo take picture from video
func PictureFromVideo(filename string) string {

	filename = strings.Replace(filename, "/file", "public", 1)
	log.Println(filename)
	format := strings.Split(filename, ".")[len(strings.Split(filename, "."))-1]
	log.Println(format)
	log.Println("ffmpeg", "-i", filename, "-ss", "00:00:01", "-vframes", "1", strings.Replace(filename, format, "jpg", 1))
	cmd, err := exec.Command("ffmpeg", "-i", filename, "-ss", "00:00:01", "-vframes", "1", strings.Replace(filename, format, "jpg", 1)).Output()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%s\n", cmd)

	filename = strings.Replace((strings.Replace(filename, "public", "/file", 1)), format, "jpg", 1)
	log.Println(filename)

	return filename

}

//GetnewsLimit get news
func GetnewsLimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")
	t := c.PostForm("type")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file 
	WHERE tnews_file.n_id = tnews.n_id and nf_type = %s order by tnews.n_date desc limit %s offset %s;`, t, limit, offset)

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

//GetnewsLimitCount get news
func GetnewsLimitCount(c *gin.Context) {

	limit := c.PostForm("limit")
	t := c.PostForm("type")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`select ceil(count(*)::real/%s::real) as pages_length from
	(SELECT * FROM public.tnews tnews, public.tnews_file tnews_file 
		WHERE tnews_file.n_id = tnews.n_id and nf_type = %s  order by tnews.n_date desc) a ;`, limit, t)

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

//Newsthemes get news theme
func Newsthemes(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT id, theme
	FROM public.newstheme order by id;`)

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

//GetNewsByTheme get news by theme
func GetNewsByTheme(c *gin.Context) {

	theme := c.PostForm("theme")

	dbConnect := config.Connect()
	todo := fmt.Sprintf(`SELECT tnews.*, tnews_file.* FROM public.tnews tnews, public.tnews_file tnews_file 
	WHERE tnews_file.n_id = tnews.n_id AND tnews.theme = %s and tnews.theme in (0,1,2,3,4) order by tnews.n_date desc;`, theme)

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
