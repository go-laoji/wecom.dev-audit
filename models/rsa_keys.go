package models

type RsaKey struct {
	Id         uint   `gorm:"primaryKey;autoIncrement;column:fi_rsa_id"`
	ExtCorpId  string `json:"ext_corp_id" gorm:"column:fs_ext_corp_id"`
	PrivateKey string `json:"private_key" gorm:"fs_private_key"`
	PublicKey  string `json:"public_key" gorm:"fs_public_key"`
	Ver        uint   `json:"ver" gorm:"fi_version"`
	BizModel
}

func (RsaKey) TableName() string {
	return "ta_rsa_key"
}
