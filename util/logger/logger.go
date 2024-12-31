package logger

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
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
var logger *logrus.Logger
var logLevel = logrus.DebugLevel
var ReportCaller bool

func init() {
	logger = logrus.New()
	// 开启收集调用信息
	logger.SetReportCaller(ReportCaller)

	// 设置自定义 formatter
	logger.SetFormatter(&CustomFormatter{
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

// CustomFormatter 可以根据自己需求自由命名
type CustomFormatter struct {
	FullTimestamp   bool
	TimestampFormat string
}

// Format 实现 logrus.Formatter 接口
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
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

// 根据你的需求对函数名进行简单截断，比如取最后一个“.”之后的部分
func trimFuncName(funcName string) string {
	// 例子：funcName = "github.com/xxx/xxx/pkg.Foo"
	// 我们想得到 "Foo"
	parts := strings.Split(funcName, ".")
	if len(parts) == 0 {
		return funcName
	}
	return parts[len(parts)-1]
}

func GinLogHandler(LogToFile bool, LogFilePath string, LogFileName string, Debug bool) gin.HandlerFunc {
	logger.SetLevel(logrus.DebugLevel) //设置最低的日志级别
	logger.SetReportCaller(Debug)      //显示函数名和行号
	logger.SetFormatter(&CustomFormatter{
		FullTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	//设置输出
	if LogToFile {
		//日志文件路径
		filename := path.Join(LogFilePath, LogFileName)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("logger err:", err)
		}
		//设置多种输出类型(默认值“os.Stderr”)
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
		//设置rotatelogs
		logWriter, err := rotatelogs.New(
			//分割后的文件名
			filename+"%Y%m%d.log",
			//指向最新日志文件的软链接
			rotatelogs.WithLinkName(filename),
			//最长保存时间
			rotatelogs.WithMaxAge(7*24*time.Hour),
			//切割间隔时间
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			fmt.Println("logger err:", err)
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
	//自定义gin处理函数
	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		// 先让请求继续处理
		c.Next()
		// 请求结束后，计算所需的时间并记录日志
		endTime := time.Now()
		latencyTime := float64(endTime.Sub(startTime).Microseconds()) / 1000
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()
		// 日志格式 https://www.runoob.com/go/go-fmt-printf.html
		// 格式化占位符的结构为：%[flags][width][.precision]verb
		//     flags：用于控制格式化输出的标志（可选）。
		//        -：左对齐。
		//        +：始终显示数值的符号。
		//        0：用零填充。
		//        #：为二进制、八进制、十六进制等加上前缀。
		//        空格：正数前加空格，负数前加 -。
		//    width：输出宽度（可选）。
		//    .precision：浮点数小数点后的位数（可选）。
		//    verb：用于指定数据的格式化方式。
		// %v	以默认格式输出变量
		// %f	十进制浮点数
		// %6.2f 表示输出浮点数，宽度为6，小数点后保留2位
		logger.WithFields(logrus.Fields{
			//"status_code":  statusCode,
			//"client_ip":    c.ClientIP(),
		}).Info(fmt.Sprintf("[%s:%d][%6.2fms][%s]%s", reqMethod, statusCode, latencyTime, c.ClientIP(), reqURI))
	}
}
