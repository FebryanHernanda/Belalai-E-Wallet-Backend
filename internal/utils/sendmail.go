package utils

import (
	"fmt"
	"os"

	"gopkg.in/mail.v2"
)

type SendOptions struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	BodyIsHTML  bool
	Attachments []string
}

func Send(opt SendOptions) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")

	if host == "" || port == "" || user == "" || pass == "" || from == "" {
		return fmt.Errorf("missing SMTP env config")
	}

	var portNum int
	fmt.Sscanf(port, "%d", &portNum)

	m := mail.NewMessage()
	m.SetHeader("From", from)
	if len(opt.To) > 0 {
		m.SetHeader("To", opt.To...)
	}
	if len(opt.Cc) > 0 {
		m.SetHeader("Cc", opt.Cc...)
	}
	if len(opt.Bcc) > 0 {
		m.SetHeader("Bcc", opt.Bcc...)
	}
	m.SetHeader("Subject", opt.Subject)

	if opt.BodyIsHTML {
		m.SetBody("text/html", opt.Body)
	} else {
		m.SetBody("text/plain", opt.Body)
	}

	for _, f := range opt.Attachments {
		m.Attach(f)
	}

	d := mail.NewDialer(host, portNum, user, pass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}
	return nil
}
