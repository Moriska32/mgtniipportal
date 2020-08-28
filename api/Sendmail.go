package api

import (
	"PortalMGTNIIP/config"
	"crypto/tls"
	"fmt"
	"net/http"

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

//SendMailITJSON json for BD
type SendMailITJSON struct {
	To       []string `json:"To"`
	Subject  string   `json:"Subject"`
	TextHTML string   `json:"text/html"`
	Author   string   `json:"author "`
	Type     string   `json:"type"`
}

//SendMailIT отправка сообщений с записью в БД
func SendMailIT(c *gin.Context) {

	var json SendMailITJSON
	c.BindJSON(&json)

	m := gomail.NewMessage()

	m.SetHeader("From", "portal@mgtniip.ru")
	m.SetHeader("To", json.To...)
	m.SetHeader("Subject", json.Subject)

	m.SetBody("text/html", json.TextHTML)

	from := "portal@mgtniip.ru"
	password := "London106446"

	emailDialer := gomail.NewDialer("exchange.mgtniip.ru", 25, from, password)
	emailDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := emailDialer.DialAndSend(m); err != nil {
		panic(err)
	}

	dbConnect := config.Connect()
	defer dbConnect.Close()

	todo := fmt.Sprintf(`INSERT INTO public.mails
	(author, "type", subject, "text","to")
	VALUES('%s', '%s', '%s', '%s','%s');`, json.Author, json.Type, json.Subject, json.TextHTML, json.To)

	_, err := dbConnect.Exec(todo)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("insert: %s", err.Error()))
	}

}
