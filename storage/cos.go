package storage

import (
	"context"
	"github.com/samber/lo"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
	"wecom.dev/audit/logger"
)

type COSStorage struct {
	client *cos.Client
}

func NewCos() (storage COSStorage, err error) {
	u, err := url.Parse(os.Getenv("CosBaseUrl"))
	if err != nil {
		logger.Surgar.Fatal("cos bucket url is invalid", err)
		return
	}
	b := &cos.BaseURL{BucketURL: u}
	storage.client = cos.NewClient(b, &http.Client{
		Timeout: 120 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv("CosSecretId"),
			SecretKey: os.Getenv("CosSecretKey"),
		},
	})
	return
}

func (o COSStorage) SignURL(objectKey string, method string, expiredInSec int64) (signedURL string, err error) {

	u, err := o.client.Object.GetPresignedURL(
		context.Background(),
		method,
		objectKey,
		os.Getenv("CosSecretId"),
		os.Getenv("CosSecretKey"),
		time.Duration(expiredInSec)*time.Second,
		nil,
	)
	signedURL = u.String()
	return
}

func (o COSStorage) Get(objectKey string) (content io.ReadCloser, err error) {
	resp, err := o.client.Object.Get(context.Background(), objectKey, nil)
	if err != nil {
		logger.Surgar.Error("Get Cos Object Failed", err)
		return
	}
	return resp.Body, nil
}
func (o COSStorage) Put(objectKey string, reader io.Reader) (err error) {
	opt := &cos.ObjectPutOptions{
		ACLHeaderOptions: &cos.ACLHeaderOptions{
			XCosACL: "private",
		},
	}
	_, err = o.client.Object.Put(context.Background(), objectKey, reader, opt)
	if err != nil {
		logger.Surgar.Error("Put Cos Object Failed", err)
		return
	}
	return
}

func (o COSStorage) IsExist(objectKey string) (ok bool, err error) {
	_, err = o.client.Object.Head(context.Background(), objectKey, nil)
	if err != nil {
		logger.Surgar.Error("Object IsExist Failed", err)
		return false, err
	}
	return true, nil
}

func (o COSStorage) PutFromFile(objectKey string, filePath string) (err error) {
	return
}

func (o COSStorage) Delete(objectKeys ...string) (deletedObjects []string, err error) {
	objects := lo.Map(objectKeys, func(t string, i int) cos.Object {
		return cos.Object{Key: t}
	})
	opt := &cos.ObjectDeleteMultiOptions{Objects: objects}
	r, _, err := o.client.Object.DeleteMulti(context.Background(), opt)
	if err != nil {
		logger.Surgar.Error("Delete Object Failed", err)
		return
	}
	deletedObjects = lo.Map(r.DeletedObjects, func(t cos.Object, i int) string {
		return t.Key
	})
	return
}
