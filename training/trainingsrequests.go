package training

import (
	"PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

//PostTrainingRequest Post Training Requset
func PostTrainingRequest(c *gin.Context) {

	data := jwt.ExtractClaims(c)

	dbConnect := config.Connect()
	defer dbConnect.Close()

	user_id := data["user_id"]
	training_id := c.PostForm("training_id")
	status_req := c.PostForm("status_req")
	date_send_req := c.PostForm("date_send_req")
	date_answer_req := c.PostForm("date_answer_req")

	sql := fmt.Sprintf(`INSERT INTO public.trainingsrequests
	( user_id, training_id, status_req, date_send_req, date_answer_req)
	VALUES(%s, %s, %s, '%s', '%s');
	`, user_id, training_id, status_req, date_send_req, date_answer_req)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//UpdateTrainingRequest Update Training Request
func UpdateTrainingRequest(c *gin.Context) {

	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()
	req_id := c.PostForm("req_id")
	user_id := c.PostForm("user_id")
	training_id := c.PostForm("training_id")
	status_req := c.PostForm("status_req")
	date_send_req := c.PostForm("date_send_req")
	date_answer_req := c.PostForm("date_answer_req")

	sql := fmt.Sprintf(`UPDATE public.trainingsrequests
	SET req_id=%s, 
	user_id=%s, training_id=%s, status_req=%s, date_send_req='%s', date_answer_req='%s';
	;
	`, req_id, user_id, training_id, status_req, date_send_req, date_answer_req)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return

}

//GetTrainingRequestsLimit Get Training Requests Limit
func GetTrainingRequestsLimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT *
	FROM public.trainingsrequests limit %s offset %s ;
	`, limit, offset)

	log.Println(todo)

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
	FROM public.trainingsrequests;
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
		"status":       http.StatusOK,
		"data":         data,
		"pages_length": count[0]["pages_length"],
	})
	dbConnect.Close()
	return

}
