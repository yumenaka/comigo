package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", // 使用完整时间戳
	})
}

func SetOutput(output io.Writer) {
	logger.SetOutput(output)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func DebugWithFields(fields logrus.Fields, args ...interface{}) {
	logger.WithFields(fields).Debug(args...)
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
