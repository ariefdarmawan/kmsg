package kmsg

import (
	"errors"
	"time"

	"git.kanosolution.net/kano/dbflex/orm"
	"git.kanosolution.net/kano/kaos"
)

type kx struct {
	senders map[string]Sender
}

func NewKaosEngine() *kx {
	k := new(kx)
	k.senders = make(map[string]Sender)
	return k
}

func (k *kx) RegisterSender(s Sender, name string) *kx {
	k.senders[name] = s
	return k
}

func (k *kx) Send(ctx *kaos.Context, id string) (string, error) {
	var e error
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("invalid hub")
	}

	m := orm.NewDataModel(new(Message)).SetObjectID(id).(*Message)
	if e = h.Get(m); e != nil {
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

			m.CreateSendAudit(h, "Fail", m.SendingAttempt, e.Error())
			return
			//return "", errors.New("process error: " + e.Error())
		}
		m.CreateSendAudit(h, "Success", m.SendingAttempt, "")

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
