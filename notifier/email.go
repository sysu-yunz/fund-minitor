package notifier

import (
	"bytes"
	"crypto/tls"
	"fund/config"
	"fund/log"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Email struct {
	To      string
	Subject string
	Body    string
}

func (e *Email) Send() {
	pwd := config.EnvVariable("EMAIL_PWD")
	d := gomail.NewDialer("smtp.qq.com", 25, "dukeyunz@foxmail.com", pwd)

	msg := gomail.NewMessage()
	msg.SetHeader("From", "dukeyunz@foxmail.com")
	msg.SetHeader("To", e.To)
	msg.SetHeader("Subject", e.Subject)
	msg.SetBody("text/html", e.Body)

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := d.DialAndSend(msg)
	if err != nil {
		log.Error("Send Error: %+v", err)
	}

	log.Debug("发送完成")
}

func (e *Email) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	e.Body = buf.String()
	return nil
}

// func (e *Email) sendEmail() {
// 	// Choose auth method and set it up
// 	pwd := config.EnvVariable("EMAIL_PWD")
// 	auth := smtp.PlainAuth("", "dukeyunz@foxmail.com", pwd, "smtp.qq.com")

// 	// Here we do it all: connect to our server, set up a message and send it
// 	to := []string{"dukeyunz@hotmail.com"}
// 	msg := []byte(
// 		"From: dukeyunz@foxmail.com\r\n" +
// 			"To: dukeyunz@hotmail.com\r\n" +
// 			"Subject: Why are you not using Mailtrap yet?\r\n" +
// 			"\r\n" +
// 			"Here’s the space for our great sales pitch\r\n")
// 	err := smtp.SendMail("smtp.qq.com:25", auth, "dukeyunz@foxmail.com", to, msg)
// 	if err != nil {
// 		log.Fatal("err")
// 	}
// 	fmt.Println("Sent!")
// }
