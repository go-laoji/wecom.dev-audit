package models

import (
	"github.com/samber/lo"
)

func GetLatestSeq(corpId string) (seq *MsgSeq, err error) {
	err = OrmEngine.Model(MsgSeq{}).Where(ChatMsg{ExtCorpID: corpId}).Limit(1).
		FirstOrCreate(&seq).Error
	return seq, err
}

type rsakeymap struct {
	PrivateKey string `json:"private_key"`
	Ver        uint32 `json:"ver"`
}

func RsaKeys(corpId string) (result map[uint32]string, err error) {
	var keys = []rsakeymap{}
	err = OrmEngine.Model(RsaKey{}).Where(RsaKey{ExtCorpId: corpId}).Find(&keys).Error
	result = make(map[uint32]string)
	if err != nil {
		return nil, err
	}
	lo.ForEach(keys, func(key rsakeymap, i int) {
		result[key.Ver] = key.PrivateKey
	})
	return
}
