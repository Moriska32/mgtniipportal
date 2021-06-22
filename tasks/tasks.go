package tasks

import (
	"PortalMGTNIIP/api"
	"PortalMGTNIIP/config"
	"PortalMGTNIIP/user"
	"fmt"
	"log"
	"net/http"

	"strings"
	"time"

	js "encoding/json"

	"net/url"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
)

var Operatortasks string = api.Operatortasks

type TasksJSON struct {
	ID                   string `json:"id"`
	TypeID               int    `json:"type_id"`
	Number               int    `json:"number"`
	Description          string `json:"description"`
	AuthorID             string `json:"author_id"`
	OperatorID           string `json:"operator_id"`
	ExecutorID           string `json:"executor_id"`
	Phone                string `json:"phone"`
	OperatorAcceptTime   string `json:"operator_accept_time"`
	OperatorDeclineTime  string `json:"operator_decline_time"`
	ExecuteStartTime     string `json:"execute_start_time"`
	ExecuteEndTime       string `json:"execute_end_time"`
	ExecuteStartPlanTime string `json:"execute_start_plan_time"`
	ExecuteEndPlanTime   string `json:"execute_end_plan_time"`
	OperatorComment      string `json:"operator_comment"`
	ExecutorComment      string `json:"executor_comment"`
	ExecuteAcceptTime    string `json:"execute_accept_time"`
	ExecuteDeclineTime   string `json:"execute_decline_time"`
	Source               string `json:"sourse"`
}

//Insert it tasks in bd
func PostTasks(c *gin.Context) {
	data := jwt.ExtractClaims(c)

	var json TasksJSON

	pool, _ := c.GetRawData()

	err := js.Unmarshal([]byte(pool), &json)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Json parse err: %s", err.Error()))
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`INSERT INTO public.tasks
	(type_id, description, phone,author_id, create_time, source)
	VALUES('%d', '%s', '%s',%s, '%s', '%s')
	RETURNING id;
	`, json.TypeID, json.Description, json.Phone, data["user_id"], time.Now().Format("2006-01-02 15:04:05"), json.Source)

	theCase := "lower"
	id, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	task := make(map[string]string)
	task["descr"] = json.Description
	task["id"] = id[0]["id"]
	task["user_id"] = fmt.Sprintf("%s", data["user_id"])

	err = api.SendLongMail(task)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	return
}

//Update it tasks in bd
func UpdateTasks(c *gin.Context) {
	//data := jwt.ExtractClaims(c)

	id := c.Query("id")
	theCase := "lower"
	var json api.TasksJSON

	pool, _ := c.GetRawData()

	err := js.Unmarshal(pool, &json)

	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Get file name err: %s", err.Error()))
	}
	log.Println(json)
	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`UPDATE public.tasks
	SET type_id = %d,operator_id=%s, executor_id=%s, operator_accept_time='%s',
	  operator_decline_time='%s', execute_start_time='%s', execute_end_time='%s',
	   execute_start_plan_time='%s', execute_end_plan_time='%s', operator_comment='%s',
	    executor_comment='%s', execute_accept_time ='%s', execute_decline_time='%s',
		source='%s'
	WHERE id='%s';	
	`, json.TypeID, json.OperatorID, json.ExecutorID, json.OperatorAcceptTime, json.OperatorDeclineTime,
		json.ExecuteStartTime, json.ExecuteEndTime, json.ExecuteStartPlanTime, json.ExecuteEndPlanTime,
		json.OperatorComment, json.ExecutorComment, json.ExecuteAcceptTime, json.ExecuteDeclineTime, json.Source, id)

	sql = strings.ReplaceAll(sql, "=,", "= null,")
	sql = strings.ReplaceAll(sql, "=''", "= null")

	_, err = dbConnect.Exec(sql)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

	switch {
	case json.ExecutorID != "" && json.ExecuteStartPlanTime != "" && json.ExecuteEndPlanTime != "" && json.ExecuteAcceptTime == "":

		err = api.SendLongMailAny(json)
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
		}
	case json.ExecutorID != "" && json.ExecuteStartTime != "" && json.ExecuteEndTime == "" && json.ExecuteAcceptTime != "" && json.OperatorAcceptTime != "":

		var jsonMail api.SendMailITJSON
		//Исполнитель
		sql = fmt.Sprintf(`select login,fam, name, user_id, userrole, tasks_role from public.tuser where
	user_id=%s;`, json.ExecutorID)
		executer, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

		sql := fmt.Sprintf(`UPDATE public.tasks
	SET execute_accept_time='%s'
	WHERE id='%s';`,
			time.Now().Format("2006-01-02 15:04:05"), id)

		_, err = dbConnect.Exec(sql)
		if err != nil {
			fmt.Println(err)
		}

		token := user.Refresher(executer[0])
		jsonMail.Subject = fmt.Sprintf(`IT-%s: вы приступили к выполнению заявки.`, json.Number)

		jsonMail.HTML = fmt.Sprintf(`

	<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta name="viewport" content="width=device-width,initial-scale=1.0">
	  <title>Письмо</title>
	</head>
	<body style="font-size: 16px;"> 
	
	<div style="margin-bottom: 20px;">Вы приступили к выполнению заявки IT-%s</div>
	   
	   
	   <a href="http://portal.mgtniip.ru:4747/v1/api/accepttaskany?token=%s&id=%s&start=0" style="display: block; padding: 10px; background: #090; color: #fff; cursor: pointer; border: none; text-decoration: none; font-size: 24px; text-align: center;">Завершить</a>
	 
	 </body>
	 </html>
`, json.Number, token, id)
		jsonMail.To = []string{executer[0]["login"]}
		api.MailSender(jsonMail)
		tokenow, _ := c.Get("JWT_TOKEN")

		inserttoken := fmt.Sprintf(`INSERT INTO public.logout
	("token")
	VALUES('%s');`, tokenow)

		_, err = dbConnect.Exec(inserttoken)

		if err != nil {
			log.Fatal("Insert token:" + err.Error())
		}

	}

	return
}

//Delete it tasks in bd
func DeleteTasks(c *gin.Context) {
	data := jwt.ExtractClaims(c)

	id := c.Query("id")

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`SELECT *
	FROM public.tasks where id = '%s' order by create_time;
	`, id)

	theCase := "lower"
	dataDB, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	if err != nil {
		log.Printf("Error while getting a single sql, Reason: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
		})
		return
	}

	if dataDB[0]["author_id"] != data["user_id"] && data["userrole"] != "2" {
		c.String(http.StatusNotAcceptable, "You are not Admin!")

		return
	} else if data["userrole"] == "2" || dataDB[0]["author_id"] == data["user_id"] {

		c.String(http.StatusAccepted, "OK")

	}

	sql = fmt.Sprintf(`DELETE FROM public.tasks
	WHERE id='%s';	
	`, id)

	_, err = dbConnect.Exec(sql)
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
	case tmin == "" && tmax == "":
		sql = fmt.Sprintf(`SELECT *
	FROM public.tasks where create_time > (now() - INTERVAL '7 DAY') order by create_time;
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

//GetTasksRoles Get Tasks Roles
func GetTasksRoles(c *gin.Context) {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`SELECT id, "role"
	FROM public.tasks_role;`)

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

func AcceptTask(c *gin.Context) {

	id := c.Query("id")

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`UPDATE public.tasks
	SET operator_accept_time='%s', operator_id = %s
	WHERE id='%s';`,
		time.Now().Format("2006-01-02 15:04:05"), Operatortasks, id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	//Данные по автору  номеру задачи
	sql = fmt.Sprintf(`select author_id, number from public.tasks where
	id='%s';`, id)
	theCase := "lower"
	task, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	//Почта юзера заявки
	sql = fmt.Sprintf(`select login from public.tuser where
	user_id=%s;`, task[0]["author_id"])

	user, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	//Оператор
	sql = fmt.Sprintf(`select fam, name from public.tuser where
	user_id=%s;`, Operatortasks)
	operator, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	var json api.SendMailITJSON

	json.To = []string{user[0]["login"]}
	json.HTML = fmt.Sprintf(`
	<!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1.0">
    <title>Письмо</title>
  </head>
  <body style="font-size: 16px;">

    <div>Вашу заявку IT-%s принял оператор %s %s. Ожидайте назначения исполнителя.</div>
    
    <a href="http://portal.mgtniip.ru/tasks">Все заявки</a>
  
  </body>
  </html>
	`, task[0]["number"], operator[0]["name"], operator[0]["fam"])

	json.Subject = fmt.Sprintf(`Вашу заявку IT-%s принял оператор %s %s.`, task[0]["number"], operator[0]["name"], operator[0]["fam"])

	api.MailSender(json)

	token, _ := c.Get("JWT_TOKEN")

	inserttoken := fmt.Sprintf(`INSERT INTO public.logout
	("token")
	VALUES('%s');`, token)

	_, err = dbConnect.Exec(inserttoken)

	if err != nil {
		log.Fatal("Insert token:" + err.Error())
	}
	loc := url.URL{Path: "http://newportal.mgtniip.ru/tasks"}
	c.Redirect(http.StatusFound, loc.RequestURI())

	return

}

func TestTest(c *gin.Context) {

	claims := jwt.ExtractClaims(c)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   claims,
	})

}

func AcceptTaskAny(c *gin.Context) {

	id := c.Query("id")
	start := c.Query("start")
	token, _ := c.Get("JWT_TOKEN")

	claims := jwt.ExtractClaims(c)

	dbConnect := config.Connect()
	defer dbConnect.Close()

	sql := fmt.Sprintf(`UPDATE public.tasks
	SET execute_accept_time='%s', executor_id = %s
	WHERE id='%s';`,
		time.Now().Format("2006-01-02 15:04:05"), claims["user_id"], id)

	_, err := dbConnect.Exec(sql)
	if err != nil {
		fmt.Println(err)
	}
	theCase := "lower"
	//данные по испольнителю

	//Исполнитель
	sql = fmt.Sprintf(`select login,fam, name, user_id, userrole, tasks_role from public.tuser where
	user_id=%s;`, claims["user_id"])
	executer, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	//Данные по автору  номеру задачи = author_id, number
	sql = fmt.Sprintf(`select author_id, number, execute_start_plan_time, execute_end_plan_time from public.tasks where
	id='%s';`, id)

	task, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	//Почта заявителя = login
	sql = fmt.Sprintf(`select login from public.tuser where
	user_id=%s;`, task[0]["author_id"])

	usersite, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	//Оператор
	sql = fmt.Sprintf(`select login,fam, name from public.tuser where
	user_id=%s;`, Operatortasks)
	operator, err := gosqljson.QueryDbToMap(dbConnect, theCase, sql)

	var json api.SendMailITJSON

	switch {
	case start == "":
		//Письмо пользователю
		json.Subject = fmt.Sprintf(`IT-%s: назначен исполнитель.`, task[0]["number"])
		json.To = []string{usersite[0]["login"]}
		json.HTML = fmt.Sprintf(`
	<!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1.0">
    <title>Письмо</title>
  </head>
  <body style="font-size: 16px;">

    <div>%s %s назначен исполнителем на заявку IT-%s</div>
	<div>Заявка будет выполнена %s с %s до %s</div>
    
    <a href="http://portal.mgtniip.ru/tasks">Все заявки</a>
  
  </body>
  </html>
	`, executer[0]["name"], executer[0]["fam"], task[0]["number"], task[0]["execute_start_plan_time"][0:10],
			task[0]["execute_start_plan_time"][11:16], task[0]["execute_end_plan_time"][11:16])
		api.MailSender(json)

		//Письмо оператору
		json.Subject = fmt.Sprintf(`IT-%s: исполнитель принял заявку.`, task[0]["number"])
		json.HTML = fmt.Sprintf(`
	<!DOCTYPE html>
  <html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width,initial-scale=1.0">
    <title>Письмо</title>
  </head>
  <body style="font-size: 16px;">

    <div>%s %s принял заявку IT-%s</div>
	
    
    <a href="http://portal.mgtniip.ru/tasks">Все заявки</a>
  
  </body>
  </html>
	`, executer[0]["name"], executer[0]["fam"], task[0]["number"])
		json.To = []string{operator[0]["login"]}
		api.MailSender(json)

		//Письмо  Исполнителю
		token = user.Refresher(executer[0])
		json.Subject = fmt.Sprintf(`IT-%s: вы приняли заявку к исполнению.`, task[0]["number"])

		json.HTML = fmt.Sprintf(`

	<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta name="viewport" content="width=device-width,initial-scale=1.0">
	  <title>Письмо</title>
	</head>
	<body style="font-size: 16px;"> 
	
	<div style="margin-bottom: 20px;">Вы приняли к исполнению заявку IT-%s</div>
	   
	   
	   <a href="http://portal.mgtniip.ru:4747/v1/api/accepttaskany?token=%s&id=%s&start=1" style="display: block; padding: 10px; background: #090; color: #fff; cursor: pointer; border: none; text-decoration: none; font-size: 24px; text-align: center;">Начать выполнение</a>
	 
	 </body>
	 </html>
`, task[0]["number"], token, id)
		json.To = []string{executer[0]["login"]}
		api.MailSender(json)
		inserttoken := fmt.Sprintf(`INSERT INTO public.logout
	("token")
	VALUES('%s');`, token)

		_, err = dbConnect.Exec(inserttoken)

		if err != nil {
			log.Fatal("Insert token:" + err.Error())
		}

		//
		//
		//Старт работы исполнителя
	case start == "1":
		//Письмо исполнителю
		sql := fmt.Sprintf(`UPDATE public.tasks
	SET execute_start_time='%s'
	WHERE id='%s';`,
			time.Now().Format("2006-01-02 15:04:05"), id)

		_, err := dbConnect.Exec(sql)
		if err != nil {
			fmt.Println(err)
		}

		token = user.Refresher(executer[0])
		json.Subject = fmt.Sprintf(`IT-%s: вы приступили к выполнению заявки.`, task[0]["number"])

		json.HTML = fmt.Sprintf(`

	<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta name="viewport" content="width=device-width,initial-scale=1.0">
	  <title>Письмо</title>
	</head>
	<body style="font-size: 16px;"> 
	
	<div style="margin-bottom: 20px;">Вы приступили к выполнению заявки IT-%s</div>
	   
	   
	   <a href="http://portal.mgtniip.ru:4747/v1/api/accepttaskany?token=%s&id=%s&start=0" style="display: block; padding: 10px; background: #090; color: #fff; cursor: pointer; border: none; text-decoration: none; font-size: 24px; text-align: center;">Завершить</a>
	 
	 </body>
	 </html>
`, task[0]["number"], token, id)
		json.To = []string{executer[0]["login"]}
		api.MailSender(json)

		inserttoken := fmt.Sprintf(`INSERT INTO public.logout
		("token")
		VALUES('%s');`, token)

		_, err = dbConnect.Exec(inserttoken)

		if err != nil {
			log.Fatal("Insert token:" + err.Error())
		}

		//Письмо оператору
		json.Subject = fmt.Sprintf(`IT-%s: %s %s начал выполнение заявки.`, task[0]["number"],
			executer[0]["name"], executer[0]["fam"])
		json.HTML = fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Письмо</title>
</head>
<body style="font-size: 16px;">

<div>%s %s начал выполнение заявки IT-%s</div>


<a href="http://portal.mgtniip.ru/tasks">Все заявки</a>

</body>
</html>
`, executer[0]["name"], executer[0]["fam"], task[0]["number"])
		json.To = []string{operator[0]["login"]}
		api.MailSender(json)

		//Остановка работы исполнителя
	case start == "0":

		sql := fmt.Sprintf(`UPDATE public.tasks
	SET execute_end_time='%s'
	WHERE id='%s';`,
			time.Now().Format("2006-01-02 15:04:05"), id)

		_, err := dbConnect.Exec(sql)
		if err != nil {
			fmt.Println(err)
		}
		//Письмо оператору
		json.Subject = fmt.Sprintf(`IT-%s: %s %s завершил выполнение заявки.`, task[0]["number"],
			executer[0]["name"], executer[0]["fam"])
		json.HTML = fmt.Sprintf(`
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Письмо</title>
</head>
<body style="font-size: 16px;">

<div>%s %s завершил выполнение заявки IT-%s</div>


<a href="http://portal.mgtniip.ru/tasks">Все заявки</a>

</body>
</html>
`, executer[0]["name"], executer[0]["fam"], task[0]["number"])
		json.To = []string{operator[0]["login"]}
		api.MailSender(json)
		inserttoken := fmt.Sprintf(`INSERT INTO public.logout
	("token")
	VALUES('%s');`, token)

		_, err = dbConnect.Exec(inserttoken)

		if err != nil {
			log.Fatal("Insert token:" + err.Error())
		}

	}
	loc := url.URL{Path: "http://newportal.mgtniip.ru/tasks"}
	c.Redirect(http.StatusFound, loc.RequestURI())

	return

}
