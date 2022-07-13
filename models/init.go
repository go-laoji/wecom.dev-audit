package models

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"os"
	"time"
	"wecom.dev/audit/logger"
)

type BizModel struct {
	CreatedAt time.Time      `gorm:"column:ft_create_time"`
	UpdatedAt time.Time      `gorm:"column:ft_update_time"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:ft_delete_time"`
}

var OrmEngine *gorm.DB

func init() {
	godotenv.Load()
	var err error
	switch os.Getenv("DataBase") {
	case "mysql":
		OrmEngine, err = gorm.Open(mysql.Open(os.Getenv("DSN")), &gorm.Config{})
		break
	case "sqlserver":
		OrmEngine, err = gorm.Open(sqlserver.Open(os.Getenv("DSN")), &gorm.Config{})
		break
	default:
		logger.Surgar.Error("不支持的数据库")
		break
	}

	if err != nil {
		logger.Surgar.Error(err)
		os.Exit(1)
	}
	OrmEngine.AutoMigrate(&RsaKey{}, &MsgSeq{})
	switch os.Getenv("StoreType") {
	case "1", "3":
		OrmEngine.AutoMigrate(&ChatMsg{}, &MsgAttachments{})
		break
	}
}
