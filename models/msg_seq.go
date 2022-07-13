package models

type MsgSeq struct {
	Id        uint   `gorm:"primaryKey;autoIncrement;column:fi_rsa_id"`
	ExtCorpId string `json:"ext_corp_id" gorm:"column:fs_ext_corp_id"`
	Seq       uint64 `json:"seq"`
	BizModel
}

func (MsgSeq) TableName() string {
	return "ta_msg_seq"
}
