package api

import (
	"github.com/gin-gonic/gin"
	wework "github.com/go-laoji/wecom-go-sdk"
	"net/http"
)

type AuditCtl struct {
}

type permitUserForm struct {
	Type int `json:"type" form:"type" binding:"oneof=0 1 2 3"`
}

func (a AuditCtl) GetPermitUserList(c *gin.Context) {
	var form permitUserForm
	if ok := c.Bind(&form); ok != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": ok.Error()})
		c.Abort()
	} else {
		if ww, exists := c.Keys["ww"].(wework.IWeWork); exists {
			resp, err := ww.GetPermitUserList(1, form.Type)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"errno": 501, "errmsg": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"errno": 0, "data": resp})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"errno": 404, "errmsg": "未找到sdk客户端"})
		}
	}
}

func (a AuditCtl) CheckSingleAgree(c *gin.Context) {
	var form wework.CheckSingleAgreeRequest
	if ok := c.Bind(&form); ok != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": ok.Error()})
		c.Abort()
	} else {
		if ww, exists := c.Keys["ww"].(wework.IWeWork); exists {
			resp, err := ww.CheckSingleAgree(1, form)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"errno": 501, "errmsg": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"errno": 0, "data": resp})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"errno": 404, "errmsg": "未找到sdk客户端"})
		}
	}
}

type groupChatForm struct {
	RoomId string `json:"roomid" form:"roomid" binding:"required"`
}

func (a AuditCtl) GroupChat(c *gin.Context) {
	var form groupChatForm
	if ok := c.Bind(&form); ok != nil {
		c.JSON(http.StatusOK, gin.H{"errno": 500, "errmsg": ok.Error()})
	} else {
		if ww, exists := c.Keys["ww"].(wework.IWeWork); exists {
			resp, err := ww.GetAuditGroupChat(1, form.RoomId)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"errno": 501, "errmsg": err.Error()})
			} else {
				c.JSON(http.StatusOK, gin.H{"errno": 0, "data": resp})
			}
		} else {
			c.JSON(http.StatusOK, gin.H{"errno": 404, "errmsg": "未找到sdk客户端"})
		}
	}
}
