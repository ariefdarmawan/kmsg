package kmsg

import (
	"time"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/dbflex/orm"
)

type sendAudit struct {
	orm.DataModelBase `bson:"-" json:"-"`
	ID                string `bson:"_id" json:"_id"`
	MessageID         string
	Attempt           int
	Method            string
	ErrorTxt          string
	Status            string
	Created           time.Time
}

func (m *sendAudit) TableName() string {
	return "KNotifMsgSendAudits"
}

func (m *sendAudit) GetID(_ dbflex.IConnection) ([]string, []interface{}) {
	return []string{"_id"}, []interface{}{m.ID}
}

func (m *sendAudit) SetID(keys ...interface{}) {
	if len(keys) > 0 {
		m.ID = keys[0].(string)
	}
}
