package tasks

import (
	"PortalMGTNIIP/config"
	"fmt"
	"log"
	"net/http"
	"time"

	js "encoding/json"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

type TasksJSON struct {
	TypeID               int    `json:"type_id"`
	Description          string `json:"description"`
	AuthorID             string `json:"author_id"`
	OperatorID           string `json:"operator_id"`
	ExecutorID           string `json:"executor_id"`
	Phone                int    `json:"phone"`
	OperatorAcceptTime   string `json:"operator_accept_time"`
	OperatorDeclineTime  string `json:"operator_decline_time"`
	ExecuteStartTime     string `json:"execute_start_time"`
	ExecuteEndTime       string `json:"execute_end_time"`
	ExecuteStartPlanTime string `json:"execute_start_plan_time"`
	ExecuteEndPlanTime   string `json:"execute_end_plan_time"`
	OperatorComment      string `json:"operator_comment"`
	ExecutorComment      string `json:"executor_comment"`
}

//Insert it tasks in bd
func PostTasks(c *gin.Context) {
	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	var json TasksJSON

	pool := c.PostForm("json")

	err := js.Unmarshal([]byte(pool), &json)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`INSERT INTO public.tasks
	(type_id, description, phone,author_id, create_time)
	VALUES('%d', '%s', %d,%s, '%s')
	;
	`, json.TypeID, json.Description, json.Phone, data["user_id"], time.Now().Format("2006-01-02 15:04:05"))

	_, err = dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return
}

//Update it tasks in bd
func UpdateTasks(c *gin.Context) {
	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	id := c.Params.ByName("id")

	var json TasksJSON

	pool := c.PostForm("json")

	err := js.Unmarshal([]byte(pool), &json)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`UPDATE public.tasks
	SET operator_id=%s, executor_id=%s, operator_accept_time='%s',
	  operator_decline_time='%s', execute_start_time='%s', execute_end_time='%s',
	   execute_start_plan_time='%s', execute_end_plan_time='%s', operator_comment='%s',
	    executor_comment='%s'
	WHERE id='%s';	
	`, json.OperatorID, json.ExecutorID, json.OperatorAcceptTime, json.OperatorDeclineTime,
		json.ExecuteStartTime, json.ExecuteEndTime, json.ExecuteStartPlanTime, json.ExecuteEndPlanTime,
		json.OperatorComment, json.ExecutorComment, id)

	_, err = dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return
}

//Delete it tasks in bd
func DeleteTasks(c *gin.Context) {
	data := jwt.ExtractClaims(c)

	if data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	}

	id := c.Query("id")

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`DELETE FROM public.tasks
	WHERE id='%s';	
	`, id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return
}

//Get it tasks in bd
func GetTasks(c *gin.Context) {

	tmin := c.Query("tmin")
	tmax := c.Query("tmax")
	id := c.Query("id")
	log.Println(tmin, tmax)
	sql := ""
	switch {

	case id != "":
		sql = fmt.Sprintf(`SELECT *
	FROM public.tasks where id = '%s' order by create_time;
	`, id)
	case tmin != "" || tmax != "":
		sql = fmt.Sprintf(`SELECT *
	FROM public.tasks where create_time between '%s' and '%s' order by create_time;
	`, tmin, tmax)
	case tmin == "" || tmax == "":
		sql = fmt.Sprintf(`SELECT *
	FROM public.tasks order by create_time;
	`)

	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	if err != nil {
		log.Printf("Error while getting a single sql, Reason: %v\n", err)
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

//Get by id it tasks in bd
func GetTasksByID(c *gin.Context) {
	id := c.Param("id")
	log.Println(id)
	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`SELECT *
	FROM public.tasks where id = '%s' order by create_time;
	`, id)
	log.Println(sql)
	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	if err != nil {
		log.Printf("Error while getting a single sql, Reason: %v\n", err)
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
