package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"sync"
)

var (
	Logger *logrus.Logger
)

type LogHook struct {
	logger    *lumberjack.Logger
	queue     chan []byte
	waitGroup sync.WaitGroup
	levels    []logrus.Level
}

func NewLogHook(logger *lumberjack.Logger, levels []logrus.Level, queueSize int) *LogHook {
	if queueSize <= 0 {
		queueSize = 64
	}
	hook := &LogHook{
		logger: logger,
		levels: levels,
		queue:  make(chan []byte, queueSize),
	}

	hook.waitGroup.Add(1)
	go func() {
		defer hook.waitGroup.Done()
		for entry := range hook.queue {
			hook.logger.Write(entry)
		}
	}()
	return hook
}

func (h *LogHook) Levels() []logrus.Level {
	return h.levels
}

func (h *LogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}

	select {
	case h.queue <- []byte(line):
	default:
		// 队列满时的则丢弃
	}
	//_, err = h.logger.Write([]byte(line))
	return err
}

func (h *LogHook) Close() {
	close(h.queue)
	h.waitGroup.Wait()
	h.logger.Close()
}

type BackupConfig struct {
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
	QueueSize  int
}
type LogConfig struct {
	Path, File, ErrorFile string
	Backups               [2]BackupConfig
}

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Path: "logs", File: "rest-server.log", ErrorFile: "error.log",
		Backups: [2]BackupConfig{
			{MaxSize: 500, MaxBackups: 30, MaxAge: 7, LocalTime: true, Compress: true, QueueSize: 2048},
			{MaxSize: 50, MaxBackups: 30, MaxAge: 30, LocalTime: true, Compress: true, QueueSize: 512},
		},
	}
}

func (c *LogConfig) SetLogConfig(cfg map[string]any) {
	if val, ok := cfg["path"]; ok {
		c.Path = val.(string)
	}

	if val, ok := cfg["file"]; ok {
		c.File = val.(string)
	}

	if val, ok := cfg["errorfile"]; ok {
		c.ErrorFile = val.(string)
	}

	if backups, ok := cfg["backups"]; ok {
		for index, bak := range backups.([]interface{}) {
			item := bak.(map[string]any)
			if val, ok := item["maxsize"]; ok {
				c.Backups[index].MaxSize = val.(int)
			}

			if val, ok := item["maxbackups"]; ok {
				c.Backups[index].MaxBackups = val.(int)
			}

			if val, ok := item["maxage"]; ok {
				c.Backups[index].MaxAge = val.(int)
			}

			if val, ok := item["localtime"]; ok {
				c.Backups[index].LocalTime = val.(bool)
			}

			if val, ok := item["compress"]; ok {
				c.Backups[index].Compress = val.(bool)
			}

			if val, ok := item["queuesize"]; ok {
				c.Backups[index].QueueSize = val.(int)
			}
		}
	}
}

func JoinPath(path, file string) string {
	if path == "" {
		return file
	}
	return path + "/" + file
}

func InitLog(cfg *LogConfig) {
	if cfg.Path != "" {
		err := os.MkdirAll(cfg.Path, 0755)
		if err != nil {
			panic(fmt.Sprintf("mkdir logs failed: %s", err.Error()))
		}
	}

	Logger = logrus.New()
	Logger.SetReportCaller(true)
	Logger.SetFormatter(&NewFormatter{})
	debugLog := &lumberjack.Logger{
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

	// 添加debug日志 Hook
	Logger.AddHook(NewLogHook(debugLog, logrus.AllLevels, cfg.Backups[0].QueueSize))

	// 添加错误日志 Hook
	Logger.AddHook(NewLogHook(errorLog,
		[]logrus.Level{logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel},
		cfg.Backups[1].QueueSize,
	))

	Logger.SetOutput(os.Stdout)

	// 测试日志
	//Logger.Info("common info level log")
	//Logger.Warn("common warn level log")
	//Logger.Error("common error level log")
}

func GetLogLevel(level string) logrus.Level {
	level = strings.ToLower(level)
	switch level {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
