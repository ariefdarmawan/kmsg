package kmsg

import (
	"errors"
	"fmt"
	"io"
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"github.com/ariefdarmawan/datahub"
	"github.com/sebarcode/codekit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id"`
	Kind              string
	From              string
	To                string
	Cc                []string
	Bcc               []string
	Title             string
	Message           string
	SendingAttempt    int
	Status            string
	Method            string
	Created           time.Time
	Sent              time.Time
}

func (m *Message) TableName() string {
	return "KNotifMessages"
}

func (m *Message) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{m.ID}
}

func (m *Message) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		m.ID = keys[0].(string)
	}
}

func (m *Message) PreSave(conn dbflex.IConnection) error {
	if m.ID == "" {
		m.ID = primitive.NewObjectID().Hex()
	}

	if m.Status == "" {
		m.Status = "Open"
		m.Created = time.Now()
	}
	return nil
}

func (m *Message) CreateAudit(h *datahub.Hub, status string, attempt int, errTxt string) error {
	sa := new(sendAudit)
	sa.ID = primitive.NewObjectID().Hex()
	sa.MessageID = m.ID
	sa.Method = m.Method
	sa.Attempt = attempt
	sa.Created = time.Now()
	sa.Status = status
	sa.ErrorTxt = errTxt

	e := h.Save(sa)
	if e != nil {
		return fmt.Errorf("fail create message send audit recor: %s", e.Error())
	}
	return nil
}

func (m *Message) Audits(h *datahub.Hub) ([]sendAudit, error) {
	res := []sendAudit{}
	w := dbflex.Eq("MessageID", m.ID)
	cmd := dbflex.From(new(Message).TableName()).Where(w)
	if _, e := h.Populate(cmd, &res); e != nil {
		return res, e
	}
	return res, nil
}

func NewMessage(h *datahub.Hub, kind, method, from, to, title, message string) (*Message, error) {
	msg := new(Message)
	msg.ID = primitive.NewObjectID().Hex()
	msg.From = from
	msg.Kind = kind
	msg.Method = method
	msg.To = to
	msg.Title = title
	msg.Message = message
	if e := h.Save(msg); e != nil {
		return nil, errors.New("system error when create message: " + e.Error())
	}
	return msg, nil
}

func NewMessageFromTemplate(h *datahub.Hub, msg *Message, templateName string, langID string, data codekit.M) error {
	var e error

	t := new(Template)
	if e = h.GetByFilter(t, dbflex.And(dbflex.Eq("Name", templateName), dbflex.Eq("LanguageID", langID))); e != nil {
		if e != io.EOF {
			return e
		}

		t.Title = templateName
		t.Message = "{{.}}"
	}

	tm, e := t.BuildMessage(data)
	if e != nil {
		return errors.New("template-error: " + e.Error())
	}
	msg.Title = tm.Title
	msg.Message = tm.Message

	if e = h.Save(msg); e != nil {
		return e
	}

	return nil
}
