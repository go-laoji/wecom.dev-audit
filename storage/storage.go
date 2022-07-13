package storage

import (
	"github.com/joho/godotenv"
	"io"
	"os"
	"wecom.dev/audit/logger"
)

type FileStorageInterface interface {
	SignURL(objectKey string, method string, expiredInSec int64) (signedURL string, err error)
	Get(objectKey string) (content io.ReadCloser, err error)
	Put(objectKey string, reader io.Reader) (err error)
	IsExist(objectKey string) (ok bool, err error)
	PutFromFile(objectKey string, filePath string) (err error)
	Delete(objectKeys ...string) (deletedObjects []string, err error)
}

var FileStorage FileStorageInterface

func init() {
	godotenv.Load()
	var err error
	switch os.Getenv("Storage") {
	case string(QiNiuStorage):
		FileStorage, err = NewQiNiu(os.Getenv("QiNiuAccessKey"), os.Getenv("QiNiuSecretKey"), os.Getenv("Bucket"))
		if err != nil {
			logger.Surgar.Error(err)
			os.Exit(1)
		}
		break
	case string(TencentStorage):
		FileStorage, err = NewCos()
		if err != nil {
			logger.Surgar.Error(err)
			os.Exit(1)
		}
		break
	default:
		logger.Surgar.Error("unsupport file storage")
	}
}
