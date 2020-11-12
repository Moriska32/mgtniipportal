package api

import (
	"PortalMGTNIIP/config"
	"crypto/tls"
	"encoding/json"
	js "encoding/json"
	"fmt"
	"log"
	"net/http"

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

	m.SetBody("text/html", json.Text)

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
	todo := fmt.Sprintf(`INSERT INTO public.mail
	("json", type_id)
	VALUES('%s'::json, %s);
	`, sql, typeid)

	_, err = dbConnect.Exec(todo)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

}

//GetRequest mail
func GetRequest(c *gin.Context) {

	typeid := c.PostForm("type_id")

	dbConnect := config.Connect()
	defer dbConnect.Close()
	todo := fmt.Sprintf(`SELECT mail.json
	FROM public.mail mail, public.mail_type mail_type
	WHERE 
		mail_type.type_id = mail.type_id and mail.type_id = %s order by mail.date desc;`, typeid)
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
	todo := fmt.Sprintf(`SELECT mail.json
	FROM public.mail mail, public.mail_type mail_type
	WHERE 
		mail_type.type_id = mail.type_id and mail.type_id = %s order by mail.date desc limit %s offset %s;`, t, limit, offset)
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

	todo = fmt.Sprintf(`SELECT ceil(count(*)::real/%s::real) as pages_length from
	(SELECT mail.json
	FROM public.mail mail, public.mail_type mail_type
	WHERE 
		mail_type.type_id = mail.type_id and mail.type_id = %s) a;`, t, limit)

	theCase := "lower"
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
		"data":   result,
		"count":  count,
	})

	return

}
