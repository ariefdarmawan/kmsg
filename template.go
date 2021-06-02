package kmsg

import (
	"bytes"
	"errors"
	"text/template"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"github.com/eaciit/toolkit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id,omitempty" json:"_id,omitempty" key:"1" grid-show:"include" form-show:"hide"`
	Name              string `kf-pos:"1,1" required:"1" search:"1"`
	LanguageID        string `kf-pos:"1,2"`
	Title             string `kf-pos:"2,1" required:"1" search:"1"`
	Group             string `kf-pos:"4,1"`
	Message           string `kf-control:"html" kf-multirow:"10" kf-pos:"3,1" required:"1" grid-show:"hide"`
}

func (m *Template) TableName() string {
	return "KNotifMsgTemplates"
}

func (m *Template) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{m.ID}
}

func (m *Template) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		m.ID = keys[0].(string)
	}
}

func (m *Template) PreSave(_ dbflex.IConnection) error {
	if m.ID == "" {
		m.ID = primitive.NewObjectID().Hex()
	}
	return nil
}

func (t *Template) BuildMessage(m toolkit.M) (*Message, error) {
	msg := new(Message)
	if s, e := translate(t.Title, m); e == nil {
		msg.Title = s
	} else {
		return nil, errors.New("fail to generate subject from template: " + e.Error())
	}

	if s, e := translate(t.Message, m); e == nil {
		msg.Messsage = s
	} else {
		return nil, errors.New("fail to generate message content from template: " + e.Error())
	}
	return nil, nil
}

func translate(source string, data toolkit.M) (string, error) {
	w := bytes.NewBufferString("")
	tt, e := template.New("tmp").Parse(source)
	if e != nil {
		return source, e
	}

	e = tt.Execute(w, data)
	if e != nil {
		return source, e
	}

	return w.String(), nil
}
