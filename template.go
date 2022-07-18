package kmsg

import (
	"errors"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
	"github.com/sebarcode/codekit"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Template struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id,omitempty" json:"_id,omitempty" key:"1" grid:"include" form:"hide"`
	Name              string `form_pos:"1,1" form_required:"1" grid_keyword:"1" grid_sortable:"1"`
	LanguageID        string `form_pos:"1,2" grid_keyword:"1" grid_sortable:"1"`
	Title             string `form_pos:"2,1" form_required:"1" grid_keyword:"1" grid_sortable:"1"`
	Group             string `form_pos:"4,1"`
	Message           string `form_kind:"html" form_multirow:"10" form_pos:"3,1" form_required:"1" grid:"hide"`
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

func (t *Template) BuildMessage(m codekit.M) (*Message, error) {
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
	return msg, nil
}

func (t *Template) Indexes() []dbflex.DbIndex {
	return []dbflex.DbIndex{
		{Name: "Name_Language_Index", IsUnique: true, Fields: []string{"Name", "LanguageID"}},
	}
}
