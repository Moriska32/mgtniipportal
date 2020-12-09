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

	dbConnect := config.Connect()
	defer dbConnect.Close()

	id := c.PostForm("id")
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

//Gettrainingstopicslimit Get training to
func Gettrainingstopicslimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT is_active, is_external, type_id, title, descr, id
	FROM public.trainingtopic order by is_active desc limit %s offset %s ;
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

//Posttraining Post training
func Posttraining(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	speakers := c.PostForm("speakers")
	users := c.PostForm("users")
	type_id := c.PostForm("type_id")
	topic_id := c.PostForm("topic_id")
	has_free_places := c.PostForm("has_free_places")
	dates_json := c.PostForm("dates_json")

	sql := fmt.Sprintf(`INSERT INTO public.training
	(speakers, users, type_id, topic_id, has_free_places, dates_json)
	VALUES(%s, %s, %s, %s, %s, '%s'::json);
	`, speakers, users, type_id, topic_id, has_free_places, dates_json)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//Updatetraining Update training
func Updatetraining(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	speakers := c.PostForm("speakers")
	users := c.PostForm("users")
	type_id := c.PostForm("type_id")
	topic_id := c.PostForm("topic_id")
	has_free_places := c.PostForm("has_free_places")
	dates_json := c.PostForm("dates_json")
	id := c.PostForm("id")

	sql := fmt.Sprintf(`UPDATE public.training
	SET speakers=%s, users=%s, type_id=%s, topic_id=%s, has_free_places=%s, dates_json='%s'
	WHERE id=%s;
	`, speakers, users, type_id, topic_id, has_free_places, dates_json, id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//Deletetraining Deletet training
func Deletetraining(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	id := c.PostForm("id")

	sql := fmt.Sprintf(`DELETE FROM public.training
	WHERE id= %s
	`, id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//Gettraining Get training by limit
func Gettraining(c *gin.Context) {

	id := c.PostForm("id")

	dbConnect := config.Connect()
	todo := fmt.Sprintf(`SELECT id, speakers, users, type_id, topic_id, has_free_places, dates_json
	FROM public.training WHERE id= %s;
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

//Gettraininglimit Get training by limit
func Gettraininglimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT id, speakers, users, type_id, topic_id, has_free_places, dates_json
	FROM public.training limit %s offset %s;
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
	FROM public.training;
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

//Gettrainingstopicstypes Get trainings topics types
func Gettrainingstopicstypes(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT type_id, "type"
	FROM public.training_type;
	`)

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
