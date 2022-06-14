package mailing

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

func SendMail(recipient string, link string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", "fontan2312@gmail.com")
	m.SetHeader("To", recipient)
	m.SetBody("text/plain", link)

	d := gomail.NewDialer("smtp.gmail.com", 587, "fontan2312@gmail.com", "lciuouwwbhchcmxf")

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil

}
