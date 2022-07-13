package storage

type StorageType string

const (
	UpYunStorage   StorageType = "upyun"
	QiNiuStorage   StorageType = "qiniu"
	AliYunStorage  StorageType = "aliyun"
	TencentStorage StorageType = "cos"
)
