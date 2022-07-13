package api

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-laoji/wecom-go-sdk/pkg/svr/logic"
	"github.com/go-laoji/wxbizmsgcrypt"
	"wecom.dev/audit/archives"
	"wecom.dev/audit/logger"
)

type CallBackCtl struct {
	Audit archives.AuditSdk
}

func (c CallBackCtl) Get(ctx *gin.Context) {

	var params logic.EventPushQueryBinding
	if ok := ctx.ShouldBindQuery(&params); ok == nil {
		receiveId := params.CorpId
		if receiveId == "" {
			receiveId = os.Getenv("CorpId")
		}
		wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(os.Getenv("Token"), os.Getenv("EncodingAESKey"),
			receiveId, wxbizmsgcrypt.XmlType)
		echoStr, cryptErr := wxcpt.VerifyURL(params.MsgSign, params.Timestamp, params.Nonce, params.EchoStr)
		if nil != cryptErr {
			logger.Surgar.Error(cryptErr)
			ctx.JSON(http.StatusOK, gin.H{"err": cryptErr, "echoStr": echoStr})
		} else {
			ctx.Writer.Write(echoStr)
		}
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"errno": 500, "errmsg": "no echostr"})
	}

}

type BizEvent struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
}
type MsgAuditNotifyEvent struct {
	BizEvent
	AgentId uint32 `xml:"AgentID"`
}

func (c CallBackCtl) Post(ctx *gin.Context) {

	var params logic.EventPushQueryBinding
	if ok := ctx.ShouldBindQuery(&params); ok == nil {
		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			logger.Surgar.Error(err)
			ctx.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.Error()})
			return
		} else {
			wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(os.Getenv("Token"), os.Getenv("EncodingAESKey"),
				os.Getenv("CorpId"), wxbizmsgcrypt.XmlType)
			if msg, err := wxcpt.DecryptMsg(params.MsgSign, params.Timestamp, params.Nonce, body); err != nil {
				logger.Surgar.Error(err)
				ctx.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.ErrMsg})
				return
			} else {
				var bizEvent BizEvent
				if e := xml.Unmarshal(msg, &bizEvent); e != nil {
					ctx.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": err.ErrMsg})
					return
				}
				switch bizEvent.Event {
				case "msgaudit_notify":
					go c.Audit.Sync()
				default:
					logger.Surgar.Error("un support event", bizEvent)

				}
				ctx.Writer.WriteString("success")
			}
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{"errno": 400, "errmsg": ok.Error()})
	}

}
