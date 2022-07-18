package ksmsg

import (
	"errors"
	"strings"
	"time"

	"git.kanosolution.net/kano/kaos"
	. "github.com/ariefdarmawan/kmsg"
	"github.com/sebarcode/codekit"
)

type kx struct {
	senders map[string]Sender
	prefix  string
}

func NewKaosModel() *kx {
	k := new(kx)
	k.senders = make(map[string]Sender)
	return k
}

func (k *kx) SetPrefix(p string) *kx {
	k.prefix = p
	return k
}

func (k *kx) RegisterSender(s Sender, name string) *kx {
	k.senders[name] = s
	return k
}

type SendTemplateRequest struct {
	Message      *Message
	TemplateName string
	LanguageID   string
	Data         codekit.M
}

func (obj *kx) Create(ctx *kaos.Context, payload *Message) (string, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("missingDBConn")
	}

	msg, e := NewMessage(h, payload.Kind, payload.Method, payload.From, payload.To, payload.Title, payload.Messsage)
	if e != nil {
		return "", e
	}

	return msg.ID, nil
}

func (obj *kx) SendTemplate(ctx *kaos.Context, request *SendTemplateRequest) (string, error) {
	var e error
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("invalid hub")
	}

	if e = NewMessageFromTemplate(h, request.Message, request.TemplateName, request.LanguageID, request.Data); e != nil {
		return "", e
	}
	go obj.SendByID(ctx, request.Message.ID)

	return request.Message.ID, nil
}

func (obj *kx) SendMessage(ctx *kaos.Context, request *Message) (string, error) {
	var e error
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("invalid hub")
	}

	if obj.prefix != "" && !strings.HasPrefix(request.Title, obj.prefix) {
		request.Title = obj.prefix + " " + request.Title
	}

	if e = h.Save(request); e != nil {
		return "", e
	}

	go obj.SendByID(ctx, request.ID)

	return request.ID, nil
}

func (k *kx) SendByID(ctx *kaos.Context, id string) (string, error) {
	var e error
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("invalid hub")
	}

	m := new(Message)
	if e = h.GetByID(m, id); e != nil {
		return "", errors.New("invalid message: " + e.Error())
	}

	sender, senderOK := k.senders[m.Method]
	if m.Status == "Sent" || m.Method == "" || !senderOK {
		return "", errors.New("invalid message")
	}

	m.Status = "Sending"
	m.SendingAttempt++
	if e = h.Save(m); e != nil {
		m.Status = "Open"
		return "", errors.New("process error: " + e.Error())
	}

	go func() {
		if e = sender.Send(m); e != nil {
			m.Status = "Fail"
			h.Save(m)

			m.CreateAudit(h, "Fail", m.SendingAttempt, e.Error())
			return
		}

		m.CreateAudit(h, "Success", m.SendingAttempt, "")
		m.Status = "Sent"
		m.Sent = time.Now()
		h.Save(m)
	}()

	return "OK", nil
}

func (k *kx) Close() {
	for _, sender := range k.senders {
		sender.Close()
	}
}
