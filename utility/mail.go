package utility

import (
	"github.com/ywengineer/smart-kit/pkg/logk"
	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

type MailClient struct {
	client gomail.SendCloser
}

func (mc *MailClient) SendMail(from, to, cc, bcc string, subject, bodyType, bodyString string) {
	var rec []string
	if ValidMail(from) == false {
		logk.Error("missing mail's sender")
		return
	}
	if ValidMail(to) == false {
		logk.Error("missing mail's to")
		return
	}
	rec = append(rec, to)
	//
	if len(cc) > 0 {
		if ValidMail(cc) == false {
			logk.Error("unknown mail's cc. %s", zap.String("cc", cc))
			return
		}
		rec = append(rec, to)
	}

	//
	if len(bcc) > 0 {
		if ValidMail(bcc) == false {
			logk.Error("unknown mail's bcc. %s", zap.String("bcc", bcc))
			return
		}
		rec = append(rec, bcc)
	}
	//
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	if len(cc) > 0 {
		m.SetAddressHeader("Cc", cc, cc)
	}
	if len(bcc) > 0 {
		m.SetAddressHeader("Bcc", bcc, bcc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody(bodyType, bodyString)
	// Send the email to Bob, Cora and Dan.
	if err := mc.client.Send(from, rec, m); err != nil {
		logk.Error("send mail failed, %v", zap.Error(err))
	}
}

//var client gomail.SendCloser

func NewMailSender(host string, port int, username, password string) (*MailClient, error) {
	d := gomail.NewDialer(host, port, username, password)
	if c, err := d.Dial(); err != nil {
		return nil, err
	} else {
		return &MailClient{client: c}, nil
	}
}

var client *MailClient

func Dial(host string, port int, username, password string) {
	if client != nil {
		logk.Error("global mail client already exists.")
		return
	}
	//
	if c, e := NewMailSender(host, port, username, password); e != nil {
		logk.Fatal("create global mail client failed. %v", zap.Error(e))
	} else {
		client = c
	}
}

func SendMail(from, to, cc, bcc string, subject, bodyType, bodyString string) {
	if client == nil {
		logk.Error("global mail client has not been created.")
		return
	}
	client.SendMail(from, to, cc, bcc, subject, bodyType, bodyString)
}

func DirectSendMail(host string, port int, username, password string,
	from, to, cc, bcc string, subject, bodyType, bodyString string) {
	if ValidMail(from) == false {
		logk.Error("missing mail's sender")
		return
	}
	if ValidMail(to) == false {
		logk.Error("missing mail's to")
		return
	}
	if len(cc) > 0 && ValidMail(cc) == false {
		logk.Error("unknown mail's cc. %s", zap.String("cc", cc))
		return
	}
	if len(bcc) > 0 && ValidMail(bcc) == false {
		logk.Error("unknown mail's bcc. %s", zap.String("bcc", bcc))
		return
	}
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	if len(cc) > 0 {
		m.SetAddressHeader("Cc", cc, cc)
	}
	if len(bcc) > 0 {
		m.SetAddressHeader("Bcc", bcc, bcc)
	}
	m.SetHeader("Subject", subject)
	m.SetBody(bodyType, bodyString)
	//m.Attach("/home/Alex/lolcat.jpg")
	d := gomail.NewDialer(host, port, username, password)
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		logk.Error("send mail failed, %v", zap.Error(err))
	}
}
