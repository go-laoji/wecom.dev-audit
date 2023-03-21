package archives

/*
#cgo CFLAGS: -I ./ -I./lib -I../lib
#cgo CXXFLAGS: -I./ -I../
#cgo LDFLAGS:  -L../lib  -lWeWorkFinanceSdk_C -ldl
#include "../lib/WeWorkFinanceSdk_C.h"
#include <stdlib.h>
*/
import "C"
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
	"gorm.io/datatypes"
	"os"
	"sync"
	"unsafe"
	"wecom.dev/audit/archives/internal"
	"wecom.dev/audit/logger"
	"wecom.dev/audit/models"
	"wecom.dev/audit/queue"
	"wecom.dev/audit/storage"
)

type AuditSdk interface {
	Sync() error
	GetMediaData(resp ChataMsg) (body []byte, err error)
	DecryptData(key string, msg string) (string, error)
}

type client struct {
	sdk      *C.WeWorkFinanceSdk_t
	lock     sync.Mutex
	idWorker *internal.Worker
}

func InitSdk(corpId string, secret string) (sdk AuditSdk, err error) {
	var c client
	corpID := C.CString(corpId)
	secretKey := C.CString(secret)
	defer func() {
		defer C.free(unsafe.Pointer(corpID))
		defer C.free(unsafe.Pointer(secretKey))
	}()
	c.sdk = C.NewSdk()
	ret := C.Init(c.sdk, corpID, secretKey)
	if ret != 0 {
		logger.Surgar.Error("sdk init failed", ret)
		os.Exit(0)
	}
	c.idWorker, _ = internal.NewWorker(1)
	return &c, nil
}

func (c *client) fetch(beginAt uint64, batchSize uint32, timeout int, proxy, proxyPwd string) ([]byte, error) {
	chatDataSlice := C.NewSlice()
	ret := C.GetChatData(c.sdk, C.ulonglong(beginAt), C.uint(batchSize), C.CString(proxy), C.CString(proxyPwd), C.int(timeout), chatDataSlice)
	if ret != 0 {
		return nil, errors.New("get chat data failed")
	}
	return []byte(C.GoString(chatDataSlice.buf)), nil
}

func (c *client) fetchAndStore() (err error) {
	c.lock.Lock()
	defer func() {
		c.lock.Unlock()
	}()
	var retry = 0
	corpKeys, err := models.RsaKeys(os.Getenv("CorpId"))
	if err != nil {
		return err
	}
	if len(corpKeys) == 0 {
		return errors.New("未配置对应的密钥")
	}
	for {
		if retry == 3 {
			break
		}
		seq, err := models.GetLatestSeq(os.Getenv("CorpId"))
		if err != nil {
			return err
		}
		chatDataSlice, err := c.fetch(
			seq.Seq,
			uint32(10),
			15,
			os.Getenv("Proxy"),
			os.Getenv("ProxyPwd"),
		)
		if err != nil {
			retry += 1
			continue
		}
		chatData := ChataDataResp{}
		json.Unmarshal(chatDataSlice, &chatData)
		if chatData.ErrCode != 0 {
			logger.Surgar.Error(chatData.ErrMsg)
			continue
		}
		if len(chatData.ChatData) == 0 {
			break
		}
		lo.ForEach(chatData.ChatData, func(item ChatData, i int) {
			key, ok := corpKeys[item.PublicKeyVer]
			if !ok {
				return
			}
			decryptKey, err := internal.RsaDecrypt(key, item.EncryptRandomKey)
			if err != nil {
				logger.Surgar.Error("解密串出错", err)
			} else {
				decryptChatMsg, err := c.DecryptData(decryptKey, item.EncryptChatMsg)
				if err != nil {
					logger.Surgar.Error("解密消息出错", err)
					return
				}
				var resp ChataMsg
				err = json.Unmarshal([]byte(decryptChatMsg), &resp)
				if err != nil {
					logger.Surgar.Error("消息序列化出错", err)
					// 出错消息不再重复获取，需要去日志里看原内容
					logger.Surgar.Info(decryptChatMsg)
					return
				}
				if resp.Action == "switch" {
					// 企业切换事件，不做记录
					return
				}
				chatMsg, err := c.msg2struct(resp, item)
				if err != nil {
					logger.Surgar.Error(err)
					logger.Surgar.Info(decryptChatMsg)
				} else {
					switch os.Getenv("StoreType") {
					case string(DataBase):
						models.OrmEngine.Model(&models.ChatMsg{}).Preload("MsgAttachments").Create(&chatMsg)
						break
					case string(Mq):
						queue.Q.Push(chatMsg)
						break
					case string(Both):
						models.OrmEngine.Model(&models.ChatMsg{}).Preload("MsgAttachments").Create(&chatMsg)
						queue.Q.Push(chatMsg)
						break
					}
				}
			}
			seq.Seq = item.Seq
			models.OrmEngine.Save(&seq)
		})
	}
	return err
}

// 消息转为定义的结构体，以便存入数据库或者是写到mq
func (c *client) msg2struct(resp ChataMsg, item ChatData) (msg *models.ChatMsg, err error) {
	msg = new(models.ChatMsg)
	msg.ID = c.idWorker.Next()
	err = copier.CopyWithOption(&msg, resp, copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	})
	if err != nil {
		logger.Surgar.Error("拷贝消息出错", err)
		return
	}
	msg.Seq = item.Seq
	msg.ExtCorpID = os.Getenv("CorpId")
	attach := models.MsgAttachments{}
	attach.ID = c.idWorker.Next()
	attach.ExtCorpID = msg.ExtCorpID
	switch resp.MsgType {
	case string(Text):
		msg.ContentText = resp.Text.Content
		_, jsonBuf, _ := c.contentJsonString(resp)
		attach.Content = datatypes.JSON(jsonBuf)
	case string(Image), string(Voice), string(Video), string(Emotion), string(File), string(MeetingVoiceCall), string(VoipDocShare):
		fileName, jsonBuf, err := c.contentJsonString(resp)
		body, err := c.GetMediaData(resp)
		if err != nil {
			logger.Surgar.Error("拉取附件出错")
		} else {
			key := fmt.Sprintf("%s/%s/%s_%s", msg.ExtCorpID, resp.MsgType, resp.MsgId, fileName)
			err = storage.FileStorage.Put(key, bytes.NewReader(body))
			if err != nil {
				logger.Surgar.Error("上传附件出错", err)
			} else {
				attach.FileURL = key
				attach.FileName = fileName
			}
		}
		attach.Content = datatypes.JSON(jsonBuf)
	default:
		_, jsonBuf, _ := c.contentJsonString(resp)
		attach.Content = datatypes.JSON(jsonBuf)
	}
	msg.MsgAttaches = attach
	return
}

func (c *client) contentJsonString(resp ChataMsg) (fileExt string, body []byte, err error) {
	switch resp.MsgType {
	case string(Text):
		body, err = json.Marshal(resp.Text)
		return "", body, err
	case string(Revoke):
		body, err = json.Marshal(resp.Revoke)
		return "", body, err
	case string(Image):
		body, err = json.Marshal(resp.Image)
		return fmt.Sprintf("%s.jpg", resp.Image.Md5sum), body, err
	case string(Agree):
		body, err = json.Marshal(resp.Agree)
		return "", body, err
	case string(DisAgree):
		body, err = json.Marshal(resp.DisAgree)
		return "", body, err
	case string(Voice):
		body, err = json.Marshal(resp.Voice)
		return fmt.Sprintf("%s.amr", resp.Voice.Md5sum), body, err
	case string(Video):
		body, err = json.Marshal(resp.Video)
		return fmt.Sprintf("%s.mp4", resp.Video.Md5Sum), body, err
	case string(Card):
		body, err = json.Marshal(resp.Card)
		return "", body, err
	case string(Location):
		body, err = json.Marshal(resp.Location)
		return "", body, err
	case string(Emotion):
		body, err = json.Marshal(resp.Emotion)
		return fmt.Sprintf("%s.jpg", resp.Emotion.Md5Sum), body, err
	case string(File):
		body, err = json.Marshal(resp.File)
		return resp.File.Filename, body, err
	case string(Link):
		body, err = json.Marshal(resp.Link)
		return "", body, err
	case string(Weapp):
		body, err = json.Marshal(resp.Weapp)
		return "", body, err
	case string(Chatrecord):
		body, err = json.Marshal(resp.Chatrecord)
		return "", body, err
	case string(Collect):
		body, err = json.Marshal(resp.Collect)
		return "", body, err
	case string(Redpacket):
		body, err = json.Marshal(resp.Redpacket)
		return "", body, err
	case string(Meeting):
		body, err = json.Marshal(resp.Meeting)
		return "", body, err
	case string(Doc):
		body, err = json.Marshal(resp.Doc)
		return "", body, err
	case string(MarkDown):
	case string(News):
		body, err = json.Marshal(resp.Info)
		return "", body, err
	case string(Calendar):
		body, err = json.Marshal(resp.Calendar)
		return "", body, err
	case string(Mixed):
		body, err = json.Marshal(resp.Mixed)
		return "", body, err
	case string(MeetingVoiceCall):
		body, err = json.Marshal(resp.MeetingVoiceCall)
		return fmt.Sprintf("%s.amr", resp.MsgId), body, err
	case string(VoipDocShare):
		body, err = json.Marshal(resp.VoipDocShare)
		return fmt.Sprintf("%s.amr", resp.MsgId), body, err
	case string(Sphfeed):
		body, err = json.Marshal(resp.Sphfeed)
		return "", body, err
	default:
		err = errors.New("未定义的消息类型")
	}
	return
}
func (c *client) Sync() error {
	return c.fetchAndStore()
}
func (c *client) GetMediaData(resp ChataMsg) (body []byte, err error) {
	var fieldId = ""
	var index = ""
	var buf bytes.Buffer
	switch resp.MsgType {
	case string(Image):
		fieldId = resp.Image.Sdkfileid
	case string(Voice):
		fieldId = resp.Voice.Sdkfileid
	case string(Video):
		fieldId = resp.Video.Sdkfileid
	case string(Emotion):
		fieldId = resp.Emotion.Sdkfileid
	case string(File):
		fieldId = resp.File.Sdkfileid
	case string(MeetingVoiceCall):
		fieldId = resp.MeetingVoiceCall.Sdkfileid
	case string(VoipDocShare):
		fieldId = resp.VoipDocShare.Sdkfileid
	default:
		err = errors.New("unsupport msg type to down")
		return nil, err
	}
	for {
		mediaData := C.NewMediaData()
		ret := C.GetMediaData(
			c.sdk, C.CString(index),
			C.CString(fieldId),
			C.CString(os.Getenv("Proxy")),
			C.CString(os.Getenv("ProxyPwd")),
			C.int(15), mediaData)
		if ret != 0 {
			return nil, errors.New(fmt.Sprintf("GetMediaData Error %v", ret))
		}
		_, err = buf.Write(C.GoBytes(unsafe.Pointer(C.GetData(mediaData)), C.int(C.GetDataLen(mediaData))))
		if err != nil {
			return nil, err
		}
		if C.IsMediaDataFinish(mediaData) == 1 {
			C.FreeMediaData(mediaData)
			break
		} else {
			index = C.GoString(C.GetOutIndexBuf(mediaData))
			C.FreeMediaData(mediaData)
		}
	}
	return buf.Bytes(), err
}

func (c *client) DecryptData(key string, msg string) (string, error) {
	decryptSlice := C.NewSlice()
	defer C.FreeSlice(decryptSlice)
	ret := C.DecryptData(C.CString(key), C.CString(msg), decryptSlice)
	if ret != 0 {
		return "", errors.New(fmt.Sprintf("DecryptData error %v", ret))
	}
	msgStr := C.GoString(decryptSlice.buf)
	return string(msgStr), nil
}
