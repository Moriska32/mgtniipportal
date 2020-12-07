package training

import (
	"PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//Posttrainingtopic Post training topic
func Posttrainingtopic(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	is_active := c.PostForm("is_active")
	is_external := c.PostForm("is_external")
	type_id := c.PostForm("type_id")
	title := c.PostForm("title")
	descr := c.PostForm("descr")

	sql := fmt.Sprintf(`INSERT INTO public.trainingtopic
	(is_active, is_external, type_id, title, descr)
	VALUES(%s, %s, %s, '%s', '%s');
	`, is_active, is_external, type_id, title, descr)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//Updatetrainingtopic Update training topic
func Updatetrainingtopic(c *gin.Context) {

	id := c.PostForm("id")

	dbConnect := config.Connect()
	defer dbConnect.Close()

	is_active := c.PostForm("is_active")
	is_external := c.PostForm("is_external")
	type_id := c.PostForm("type_id")
	title := c.PostForm("title")
	descr := c.PostForm("descr")

	sql := fmt.Sprintf(`UPDATE public.trainingtopic
	SET is_active=%s, is_external=%s, type_id=%s, title='%s', descr='%s'
	WHERE id= %s;
	`, is_active, is_external, type_id, title, descr, id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//Gettrainingtopic Get training to
func Gettrainingtopic(c *gin.Context) {

	id := c.PostForm("id")

	dbConnect := config.Connect()
	todo := fmt.Sprintf(`SELECT is_active, is_external, type_id, title, descr, id
	FROM public.trainingtopic WHERE id= %s;
	`, id)

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

//Gettrainingtopicslimit Get training to
func Gettrainingtopicslimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT is_active, is_external, type_id, title, descr, id
	FROM public.trainingtopic limit %s offset %s order by is_active desc;
	`, limit, offset)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	todo = fmt.Sprintf(`SELECT ceil(count(*)::real/%s::real) as pages_length
	FROM public.trainingtopic;
	`, limit)

	count, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

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
		"count":  count,
	})
	dbConnect.Close()
	return

}
