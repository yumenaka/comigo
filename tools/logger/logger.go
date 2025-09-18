package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comigo/tools/sse_hub"
)

// 对外暴露的 log 接口
var (
	SetOutput func(output io.Writer)
	Info      func(args ...interface{})
	Infof     func(format string, args ...interface{})
	Error     func(args ...interface{})
	Errorf    func(format string, args ...interface{})
	Fatal     func(args ...interface{})
	Fatalf    func(format string, args ...interface{})
)

// 全局 logger
var (
	logger       *logrus.Logger
	logLevel     = logrus.DebugLevel
	ReportCaller bool
)

func init() {
	logger = logrus.New()
	// 开启收集调用信息
	logger.SetReportCaller(ReportCaller)

	// 设置自定义 formatter
	logger.SetFormatter(&EchoLogFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 将 logger 的输出函数赋值给外部调用
	// 每个包的包级变量初始化是在 init 函数执行前完成的。所以不可以用 var Info = logger.Info 这种写法，会因为 logger 未初始化，导致空指针异常
	SetOutput = logger.SetOutput
	Info = logger.Info
	Infof = logger.Infof
	Error = logger.Error
	Errorf = logger.Errorf
	Fatal = logger.Fatal
	Fatalf = logger.Fatalf
}

// EchoLogFormatter 可以根据自己需求自由命名
type EchoLogFormatter struct {
	FullTimestamp   bool
	TimestampFormat string
}

// Format 实现 logrus.Formatter 接口
func (f *EchoLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// 如果 entry.Buffer 不为空，就使用原有的 buffer
	// 否则新建一个 bytes.Buffer
	b := entry.Buffer
	if b == nil {
		b = &bytes.Buffer{}
	}

	// 格式化时间
	timestamp := entry.Time.Format(f.TimestampFormat)

	// 默认的调用信息
	var callerInfo string
	if ReportCaller && entry.HasCaller() {
		// 例如：/path/to/pkg/file.go
		file := entry.Caller.File
		// 行号
		line := entry.Caller.Line
		// 完整函数名，例如：github.com/xxx/xxx/pkg.Foo
		funcName := entry.Caller.Function

		// 仅输出文件名
		_, fileName := filepath.Split(file)

		// 有时函数名前会带包名、路径，可以根据需要做简化
		shortFuncName := trimFuncName(funcName)

		// 组合 caller 信息
		callerInfo = fmt.Sprintf("%s:%d %s()", fileName, line, shortFuncName)
	} else {
		callerInfo = ""
	}

	// 拼装日志内容
	var logLine string
	if ReportCaller {
		logLine = fmt.Sprintf("[%s][%s]%s\n",
			strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	} else {
		logLine = fmt.Sprintf("[%s][%s]%s%s\n",
			strings.ToUpper(entry.Level.String()), timestamp, callerInfo, entry.Message)
	}
	b.WriteString(logLine)
	return b.Bytes(), nil
}

// 根据你的需求对函数名进行简单截断，比如取最后一个"."之后的部分
func trimFuncName(funcName string) string {
	// 例子：funcName = "github.com/xxx/xxx/pkg.Foo"
	// 我们想得到 "Foo"
	parts := strings.Split(funcName, ".")
	if len(parts) == 0 {
		return funcName
	}
	return parts[len(parts)-1]
}

func EchoLogHandler(LogToFile bool, LogFilePath string, LogFileName string, Debug bool) echo.MiddlewareFunc {
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(Debug)
	logger.SetFormatter(&EchoLogFormatter{
		FullTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 设置输出
	if LogToFile {
		// 日志文件路径
		filename := path.Join(LogFilePath, LogFileName)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logger.Info("logger err:", err)
		}
		// 设置多种输出类型(默认值“os.Stderr”)
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
		// 设置rotatelogs
		logWriter, err := rotatelogs.New(
			// 分割后的文件名
			filename+"%Y%m%d.log",
			// 指向最新日志文件的软链接
			rotatelogs.WithLinkName(filename),
			// 最长保存时间
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// 切割间隔时间
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			logger.Info("logger err:", err)
		}
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}
		logger.AddHook(lfshook.NewHook(writeMap, &logrus.TextFormatter{
			TimestampFormat: "2006-01-02 03:04:05",
		}))
	}
	// 自定义log处理函数
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			err := next(c)
			endTime := time.Now()
			latencyTime := float64(endTime.Sub(startTime).Microseconds()) / 1000
			reqMethod := c.Request().Method
			reqURI := c.Request().RequestURI
			statusCode := c.Response().Status
			logMsg := fmt.Sprintf("[%s:%d][%6.2fms][%s]%s",
				reqMethod,
				statusCode,
				latencyTime,
				c.RealIP(),
				reqURI,
			)
			logMsgWeb := fmt.Sprintf("[%s:%d]<span style=\"color:#d08700\">[%6.2fms]</span><span style=\"color:#0084d1\">[%s]</span>%s",
				reqMethod,
				statusCode,
				latencyTime,
				c.RealIP(),
				reqURI,
			)
			// 把log发送给所有网页客户端 <span style="color:green">[GET:200]</span>
			nowTimeStr := "<span style=\"color:oklch(62.7% 0.194 149.214)\">[" + time.Now().Format("2006-01-02 15:04:05") + "]</span>"
			sse_hub.MessageHub.Broadcast(sse_hub.Event{
				Name: "log",
				ID:   fmt.Sprintf("%d", time.Now().UnixNano()),
				Data: fmt.Sprintf("%s%s", nowTimeStr, logMsgWeb),
			})
			logger.WithFields(logrus.Fields{}).Info(
				logMsg,
			)
			return err
		}
	}
}
