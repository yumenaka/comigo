package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", // 使用完整时间戳
	})
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
