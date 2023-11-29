package logger

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comi/config"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", // 使用完整时间戳
	})
}

// SetLogger 设置日志中间件
func SetLogger(engine *gin.Engine) {
	//禁止控制台输出
	gin.DefaultWriter = io.Discard
	//设置log中间件 TODO:输出到tui界面。
	engine.Use(HandlerLog(config.Config.LogToFile, config.Config.LogFilePath, config.Config.LogFileName))
	if config.Config.LogToFile {
		// 关闭 log 打印的字体颜色。输出到文件不需要颜色
		gin.DisableConsoleColor()
	}
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func DebugWithFields(fields logrus.Fields, args ...interface{}) {
	Log.WithFields(fields).Debug(args...)
}

type MyFormatter struct {
	FullTimestamp   bool
	TimestampFormat string // 使用完整时间戳
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	var newLog string
	newLog = fmt.Sprintf("[%s] [%s] %s ", timestamp, entry.Level, entry.Message)
	b.WriteString(newLog)
	return b.Bytes(), nil
}
