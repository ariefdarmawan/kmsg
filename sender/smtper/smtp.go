package smtper

import (
	"github.com/ariefdarmawan/kmsg"
	"gopkg.in/gomail.v2"
	"crypto/tls"
)

type Options struct {
	Server       string
	Port         int
	TLS          bool
	SkipVerify bool
	Certificates []string
	UID          string
	Password     string
}

type smtp struct {
	opts Options
}

func NewSender(opts Options) *smtp {
	s := new(smtp)
	s.opts = opts
	return s
}

func (s *smtp) Send(msg *kmsg.Message) error {
	d := gomail.NewDialer(s.opts.Server, s.opts.Port, s.opts.UID, s.opts.Password)
	
	// todo - tls and certificates
	if s.opts.TLS && s.opts.SkipVerify  {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	

	m := gomail.NewMessage()
	if msg.From == "" {
		m.SetHeader("From", s.opts.UID)
	} else {
		m.SetHeader("From", msg.From)
	}
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Title)
	m.SetBody("text/html", msg.Message)

	// todo - attachment

	if e := d.DialAndSend(m); e != nil {
		return e
	}
	return nil
}

func (s *smtp) Close() {
	//panic("not implemented") // TODO: Implement
}
