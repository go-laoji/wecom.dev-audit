package main

import (
	"context"
	"github.com/gin-gonic/gin"
	wework "github.com/go-laoji/wecom-go-sdk"
	"github.com/go-laoji/wecom-go-sdk/pkg/svr/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"wecom.dev/audit/api"
	"wecom.dev/audit/archives"
	"wecom.dev/audit/logger"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ww := wework.NewWeWork(wework.WeWorkConfig{})
	ww.SetAppSecretFunc(func(corpId uint) (corpid string, secret string, customizedApp bool) {
		return os.Getenv("CorpId"), os.Getenv("Secret"), true
	})
	client, err := archives.InitSdk(os.Getenv("CorpId"), os.Getenv("Secret"))
	err = client.Sync()
	logger.Surgar.Info("初始化启动时同步完成，开启web服务监听")

	router := gin.Default()
	callbackCtl := api.CallBackCtl{Audit: client}
	callback := router.Group("/callback")
	{
		callback.GET("", callbackCtl.Get)
		callback.POST("", callbackCtl.Post)
	}
	apiAudit := router.Group("/api/audit")
	apiAudit.Use(middleware.InjectSdk(ww))
	auditCtl := api.AuditCtl{}
	{
		apiAudit.POST("groupchat", auditCtl.GroupChat)
		apiAudit.POST("checkagree", auditCtl.CheckSingleAgree)
		apiAudit.POST("permituser", auditCtl.GetPermitUserList)
	}

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "ok"})
	})
	srv01 := &http.Server{
		Addr:           "127.0.0.1:8080",
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := srv01.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Surgar.Fatal(err)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv01.Shutdown(ctx); err != nil {
		logger.Surgar.Fatal("server shutdown:", err)
	}
}
