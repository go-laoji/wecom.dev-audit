package storage

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"io"
	"io/ioutil"
	"time"
	"wecom.dev/audit/logger"
)

type QiNiuOSSStorage struct {
	mac    *qbox.Mac
	bucket string
}

func NewQiNiu(accessKey string, secretKey string, bucket string) (oss QiNiuOSSStorage, err error) {
	oss.bucket = bucket
	oss.mac = qbox.NewMac(accessKey, secretKey)
	return
}

func (oss QiNiuOSSStorage) Put(objectKey string, reader io.Reader) (err error) {
	putPolicy := storage.PutPolicy{
		Scope: oss.bucket,
	}
	upToken := putPolicy.UploadToken(oss.mac)
	cfg := storage.Config{}
	// 空间对应的机房
	//cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	data, _ := ioutil.ReadAll(reader)
	dataLen := int64(len(data))
	err = formUploader.Put(context.Background(), &ret, upToken, objectKey, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		logger.Surgar.Error(err)
		return
	}
	return
}

func (oss QiNiuOSSStorage) SignURL(objectKey string, method string, expiredInSec int64) (signedURL string, err error) {
	deadline := time.Now().Add(time.Second * time.Duration(expiredInSec)).Unix()
	signedURL = storage.MakePrivateURL(oss.mac, "", objectKey, deadline)
	return
}

func (oss QiNiuOSSStorage) Get(objectKey string) (content io.ReadCloser, err error) {
	return nil, err
}

func (oss QiNiuOSSStorage) IsExist(objectKey string) (ok bool, err error) {
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(oss.mac, &cfg)
	_, sErr := bucketManager.Stat(oss.bucket, objectKey)
	if sErr != nil {
		return false, sErr
	}
	return true, err
}

func (oss QiNiuOSSStorage) PutFromFile(objectKey string, filePath string) (err error) {
	return err
}

func (oss QiNiuOSSStorage) Delete(objectKeys ...string) (deletedObjects []string, err error) {
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: true,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	bucketManager := storage.NewBucketManager(oss.mac, &cfg)
	deleteOps := make([]string, 0, len(objectKeys))
	for _, key := range objectKeys {
		deleteOps = append(deleteOps, storage.URIDelete(oss.bucket, key))
	}
	rets, err := bucketManager.Batch(deleteOps)
	if err != nil {
		// 遇到错误
		if _, ok := err.(*storage.ErrorInfo); ok {
			for _, ret := range rets {
				if ret.Code != 200 {
					logger.Surgar.Error(ret.Data.Error)
				} else {
					deletedObjects = append(deletedObjects, ret.Data.Hash)
				}
			}
		} else {
			logger.Surgar.Error(fmt.Printf("batch error, %s", err))
		}
	} else {
		// 完全成功
		for _, ret := range rets {
			deletedObjects = append(deletedObjects, ret.Data.Hash)
			fmt.Printf("%d\n", ret.Code)
		}
	}
	return deletedObjects, err
}
