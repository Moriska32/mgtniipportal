package api

import (
	"crypto/tls"

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
