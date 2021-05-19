package api

import (
	"PortalMGTNIIP/config"

	"crypto/tls"
	"encoding/json"
	js "encoding/json"
	"fmt"
	"log"
	"net/http"

	user "PortalMGTNIIP/user"

	"github.com/elgs/gosqljson"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

// SendMail Send mail somebody
func SendMail(c *gin.Context) {

	Header := c.PostForm("header")
	Body := c.PostForm("body")
	To := c.PostFormArray("to")
	m := gomail.NewMessage()

	m.SetHeader("From", "portal@mgtniip.ru")
	m.SetHeader("To", To...)
	m.SetHeader("Subject", Header)

	m.SetBody("text/html", Body)

	from := "portal@mgtniip.ru"
	password := "London106446"

	emailDialer := gomail.NewDialer("exchange.mgtniip.ru", 25, from, password)
	emailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := emailDialer.DialAndSend(m); err != nil {
		panic(err)
	}

}

//SendMailSITJSON json for BD
type SendMailSITJSON []SendMailITJSON

//SendMailITJSON json for BD
type SendMailITJSON struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	HTML    string   `json:"html"`
	UserID  int      `json:"user_id"`
	Type    string   `json:"type"`
	Date    string   `json:"date"`
	DepID   int      `json:"dep_id"`
}

//SendRequest отправка сообщений с записью в БД
func SendRequest(c *gin.Context) {

	var json SendMailITJSON

	m := gomail.NewMessage()

	pool := c.PostForm("json")

	typeid := c.PostForm("type_id")
	err := js.Unmarshal([]byte(pool), &json)

	m.SetHeader("From", "portal@mgtniip.ru")
	m.SetHeader("To", json.To...)
	m.SetHeader("Subject", json.Subject)

	m.SetBody("text/html", json.HTML)

	from := "portal@mgtniip.ru"
	password := "London106446"

	emailDialer := gomail.NewDialer("exchange.mgtniip.ru", 25, from, password)
	emailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := emailDialer.DialAndSend(m); err != nil {
		panic(err)
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()
	sql, err := js.Marshal(json)
	_ = sql
	todo := fmt.Sprintf(`INSERT INTO public.mail
	("json", type_id)
	VALUES('%s'::json, %s) RETURNING id;
	`, pool, typeid)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		log.Printf("Error while getting a single todo, Reason: %v\n", err)

	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

	return
}

//GetRequest mail
func GetRequest(c *gin.Context) {

	typeid := c.PostForm("type_id")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT mail.json
	FROM public.mail mail, public.mail_type mail_type
	WHERE 
		mail_type.type_id = mail.type_id and mail.type_id = %s order by cast(mail.json ->> 'date' as timestamp) desc;`, typeid)
	var (
		pool string
		data SendMailITJSON
	)
	sql, _ := dbConnect.Query(todo)

	var result SendMailSITJSON

	for sql.Next() {
		sql.Scan(&pool)
		//pool = strings.Replace(pool, `\`, ``, 1)
		err := json.Unmarshal([]byte(pool), &data)

		if err != nil {
			panic(err.Error())
		}
		result = append(result, data)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   result,
	})

	return

}

//GetRequestLimit mail
func GetRequestLimit(c *gin.Context) {

	limit := c.PostForm("limit")
	offset := c.PostForm("offset")
	t := c.PostForm("type")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT mail.json, mail.id
	FROM public.mail mail, public.mail_type mail_type
	WHERE 
		mail_type.type_id = mail.type_id and mail.type_id = %s order by cast(mail.json ->> 'date' as timestamp) desc limit %s offset %s;`, t, limit, offset)

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	todo = fmt.Sprintf(`SELECT ceil(count(*)::real/%s::real) as pages_length from
	(SELECT mail.json
	FROM public.mail mail, public.mail_type mail_type
	WHERE 
		mail_type.type_id = mail.type_id and mail.type_id = %s) a;`, limit, t)

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
		"count":  count[0]["pages_length"],
	})

	return

}

func MailSender(json SendMailITJSON) {

	m := gomail.NewMessage()

	m.SetHeader("From", "portal@mgtniip.ru")
	m.SetHeader("To", json.To...)
	m.SetHeader("Subject", json.Subject)

	m.SetBody("text/html", json.HTML)

	from := "portal@mgtniip.ru"
	password := "London106446"

	emailDialer := gomail.NewDialer("exchange.mgtniip.ru", 25, from, password)
	emailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := emailDialer.DialAndSend(m); err != nil {
		panic(err)
	}
	return
}

//SendLongMail send long mail
func SendLongMail(task map[string]string) error {

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`SELECT user_id, login, fam, "name", otch, tel, userrole, tasks_role
	FROM public.tuser where user_id = %s;`, "507")

	theCase := "lower"
	data, err := gosqljson.QueryDbToMap(dbConnect, theCase, todo)

	if err != nil {
		return err

	}

	token := user.Refresher(data[0])

	textmail := fmt.Sprintf(`

	<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="utf-8">
	  <meta http-equiv="X-UA-Compatible" content="IE=edge">
	  <meta name="viewport" content="width=device-width,initial-scale=1.0">
	  <title>Письмо</title>
	</head>
	<body style="font-size: 16px;">
	 
	   <div style="margin-bottom: 20px;">%s %s создал новую заявку: %s</div>
	   <div style="margin-bottom: 20px;">Обратная связь: %s</div>
	   
	   <a href="http://172.20.82:4747/v1/api/accepеttask?token=%s&id=%s" style="display: block; padding: 10px; background: #090; color: #fff; cursor: pointer; border: none; text-decoration: none; font-size: 24px; text-align: center;">Принять</a>
	 
	 </body>
	 </html>
`, data[0]["name"], data[0]["fam"], task["descr"],
		data[0]["tel"], token, task["id"])

	log.Println(textmail)

	var json SendMailITJSON

	json.HTML = textmail
	json.To = []string{data[0]["login"]}
	json.Subject = fmt.Sprintf(`Новая заявка %s от  %s %s :  %s `,
		task["number"], data[0]["name"], data[0]["fam"], task["descr"])

	MailSender(json)

	// inserttoken := fmt.Sprintf(`INSERT INTO public.logout
	// ("token")
	// VALUES('%s');`, token)

	// _, err = dbConnect.Exec(inserttoken)

	// if err != nil {
	// 	log.Fatal("Insert token:" + err.Error())
	// }

	return nil

}
