package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	logger *zap.Logger
	Surgar *zap.SugaredLogger
)

func init() {
	hook := &lumberjack.Logger{
		Filename:   "./logs/audit.log", // 日志文件路径
		MaxSize:    500,                // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 3,                  // 日志文件最多保存多少个备份
		MaxAge:     28,                 // 文件最多保存多少天
		Compress:   true,               // 是否压缩
	}
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(hook)), // 打印到控制台和文件
		zap.InfoLevel, // 日志级别
	)
	logger = zap.New(core, zap.AddCaller())
	Surgar = logger.Sugar()
	defer logger.Sync()
}
