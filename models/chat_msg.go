package models

import (
	"encoding/json"
	"gorm.io/datatypes"
)

type ChatMsg struct {
	ID int64 `json:"id" gorm:"primaryKey;type:bigint;comment:'ID'"`
	// ExtCorpID 外部企业ID
	ExtCorpID string `json:"ext_corp_id" gorm:"column:fs_ext_corp_id;index;type:char(18);comment:外部企业ID"`
	//消息id，消息的唯一标识，企业可以使用此字段进行消息去重。String类型
	MsgId string `gorm:"type:char(128);unique" json:"msgid"`
	//消息动作，目前有send(发送消息)/recall(撤回消息)/switch(切换企业日志)三种类型。String类型
	Action string `gorm:"type:char(8)" json:"action"`
	//消息发送方id。同一企业内容为userid，非相同企业为external_userid。消息如果是机器人发出，也为external_userid。String类型
	From string `gorm:"type:char(32)" json:"from"`
	//消息接收方列表，可能是多个，同一个企业内容为userid，非相同企业为external_userid。数组，内容为string类型
	ToList StringArrayField `gorm:"type:json" json:"tolist"`
	//群聊消息的群id。如果是单聊则为空。String类型
	RoomId string `gorm:"type:char(128)" json:"roomid"`
	//消息发送时间戳，utc时间，ms单位。
	MsgTime int64 `gorm:"type:bigint(64)" json:"msgtime"`
	//文本消息为：text。String类型
	MsgType string `gorm:"type:varchar(32)" json:"msgtype"`
	// 聊天的文本内容
	ContentText string `gorm:"type:text;class:FULLTEXT,option:WITH PARSER ngram INVISIBLE" json:"content_text"`
	//消息的seq值，标识消息的序号。再次拉取需要带上上次回包中最大的seq。Uint64类型，范围0-pow(2,64)-1
	Seq         uint64         `gorm:"column:fi_seq;type:bigint unsigned" json:"seq"`
	MsgAttaches MsgAttachments `gorm:"foreignKey:chat_msg_id;reference:ID" json:"chat_msg_content"`
	// 客户名/员工名/群名
	Keywords string `gorm:"comment:用于搜索会话" json:"keywords"`
	BizModel
}

func (ChatMsg) TableName() string {
	return "ta_chat_msg"
}

func (i ChatMsg) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}

type MsgAttachments struct {
	ID int64 `json:"id" gorm:"primaryKey;type:bigint;comment:'ID'" validate:"int64"`
	// ExtCorpID 外部企业ID
	ExtCorpID string `json:"ext_corp_id" gorm:"index;type:char(18);comment:外部企业ID" validate:"ext_corp_id"`
	ChatMsgID string `json:"chat_msg_id"`
	// 文件类型
	ContentType string `json:"content_type"`
	// 聊天的非文字内容
	Content  datatypes.JSON `gorm:"type:json" json:"content"`
	FileURL  string         `json:"file_url"`
	FileName string         `json:"file_name"`
	BizModel
}

func (MsgAttachments) TableName() string {
	return "ta_msg_attachment"
}
