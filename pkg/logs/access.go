package logs

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"runtime/debug"
)

var (
	AccessLogger *logrus.Logger
)

func InitAccessLog(cfg *LogConfig) {
	if cfg.Path != "" {
		err := os.MkdirAll(cfg.Path, 0755)
		if err != nil {
			panic(fmt.Sprintf("mkdir logs failed: %s", err.Error()))
		}
	}

	AccessLogger = logrus.New()
	infoLog := &lumberjack.Logger{
		Filename:   JoinPath(cfg.Path, cfg.File),
		MaxSize:    cfg.Backups[0].MaxSize,    // 缺省500MB
		MaxBackups: cfg.Backups[0].MaxBackups, // 缺省保留30个备份
		MaxAge:     cfg.Backups[0].MaxAge,     // 缺省保留7天
		LocalTime:  cfg.Backups[0].LocalTime,  // 缺省使用本地时间命名备份文件
		Compress:   cfg.Backups[0].Compress,   // 缺省压缩旧日志
	}

	errorLog := &lumberjack.Logger{
		Filename:   JoinPath(cfg.Path, cfg.ErrorFile),
		MaxSize:    cfg.Backups[1].MaxSize,    // 缺省50MB
		MaxBackups: cfg.Backups[1].MaxBackups, // 缺省保留30个备份
		MaxAge:     cfg.Backups[1].MaxAge,     // 缺省保留30天
		LocalTime:  cfg.Backups[1].LocalTime,  // 缺省使用本地时间命名备份文件
		Compress:   cfg.Backups[1].Compress,   // 缺省压缩旧日志
	}

	// 添加info日志 Hook
	AccessLogger.AddHook(NewLogHook(
		infoLog, []logrus.Level{logrus.InfoLevel}, cfg.Backups[0].QueueSize))

	// 添加错误日志 Hook
	AccessLogger.AddHook(NewLogHook(
		errorLog, []logrus.Level{logrus.ErrorLevel}, cfg.Backups[1].QueueSize))

	AccessLogger.SetOutput(os.Stdout)

	// 测试日志
	//AccessLogger.Info("access info level log")
	//AccessLogger.Warn("access warn level log")
	//AccessLogger.Error("access error level log")
}

func AccessInfo(param gin.LogFormatterParams) string {
	AccessLogger.WithFields(logrus.Fields{
		"status":     param.StatusCode,
		"method":     param.Method,
		"path":       param.Path,
		"client_ip":  param.ClientIP,
		"latency":    param.Latency,
		"user_agent": param.Request.UserAgent(),
	}).Info("HTTP request")
	return ""
}

func RecoveryError(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		AccessLogger.WithFields(logrus.Fields{
			"path":   c.Request.URL.Path,
			"method": c.Request.Method,
			"client": c.ClientIP(),
			"error":  err,
			"stack":  string(debug.Stack()),
		}).Error("HTTP panic recovered")
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
