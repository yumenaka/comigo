package logger

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/yumenaka/comigo/tools/sse_hub"
)

// 对外暴露的 log 接口
var (
	SetOutput          func(output io.Writer)
	SetMirrorOutput    func(output io.Writer)
	SetSuppressStdout  func(disable bool)
	IsStdoutSuppressed func() bool
	Info               func(args ...interface{})
	Infof              func(format string, args ...interface{})
	Error              func(args ...interface{})
	Errorf             func(format string, args ...interface{})
	Fatal              func(args ...interface{})
	Fatalf             func(format string, args ...interface{})
)

// 全局 logger
var (
	logger         *logrus.Logger
	baseOutput     io.Writer
	mirrorOutput   io.Writer
	suppressStdout bool
	outputMu       sync.Mutex
	fileOutput     *rotatelogs.RotateLogs
)

func init() {
	logger = logrus.New()
	baseOutput = os.Stderr

	// 设置自定义 formatter
	logger.SetFormatter(&EchoLogFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// 将 logger 的输出函数赋值给外部调用
	// 每个包的包级变量初始化是在 init 函数执行前完成的。所以不可以用 var Info = logger.Info 这种写法，会因为 logger 未初始化，导致空指针异常
	SetOutput = func(output io.Writer) {
		outputMu.Lock()
		defer outputMu.Unlock()
		baseOutput = output
		applyOutput()
	}
	SetMirrorOutput = func(output io.Writer) {
		outputMu.Lock()
		defer outputMu.Unlock()
		mirrorOutput = output
		applyOutput()
	}
	SetSuppressStdout = func(disable bool) {
		outputMu.Lock()
		defer outputMu.Unlock()
		suppressStdout = disable
		applyOutput()
	}
	IsStdoutSuppressed = func() bool {
		outputMu.Lock()
		defer outputMu.Unlock()
		return suppressStdout
	}
	Info = logger.Info
	Infof = logger.Infof
	Error = logger.Error
	Errorf = logger.Errorf
	Fatal = logger.Fatal
	Fatalf = logger.Fatalf
	applyOutput()
}

// logOutputs 将同一条日志交给每个目的地，终端或镜像写入失败仍可写文件。
type logOutputs []io.Writer

func (outputs logOutputs) Write(data []byte) (int, error) {
	var writeErr error
	for _, output := range outputs {
		n, err := output.Write(data)
		if err == nil && n != len(data) {
			err = io.ErrShortWrite
		}
		if err != nil {
			writeErr = errors.Join(writeErr, err)
		}
	}
	if writeErr != nil {
		return 0, writeErr
	}
	return len(data), nil
}

// applyOutput 配置终端、TUI 镜像和轮转文件，切换文件模式不改写原终端设置。
func applyOutput() {
	var outputs logOutputs
	if !suppressStdout {
		console := baseOutput
		if console == os.Stderr && fileOutput != nil {
			console = os.Stdout
		}
		if console == nil && mirrorOutput == nil {
			console = os.Stderr
		}
		if console != nil {
			outputs = append(outputs, console)
		}
	}
	if mirrorOutput != nil {
		outputs = append(outputs, mirrorOutput)
	}
	if fileOutput != nil {
		outputs = append(outputs, fileOutput)
	}
	logger.SetOutput(outputs)
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
	if entry.HasCaller() {
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
		callerInfo = fmt.Sprintf("%s:%d %s() ", fileName, line, shortFuncName)
	} else {
		callerInfo = ""
	}

	// 拼装日志内容
	logLine := fmt.Sprintf("[%s][%s]%s%s\n",
		strings.ToUpper(entry.Level.String()), timestamp, callerInfo, entry.Message)
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

// EchoLogHandler 配置日志输出并记录 HTTP 请求。
func EchoLogHandler(LogToFile bool, LogFilePath string, LogFileName string, Debug bool) echo.MiddlewareFunc {
	logger.SetLevel(logrus.DebugLevel)
	logger.SetReportCaller(Debug)
	logger.SetFormatter(&EchoLogFormatter{
		FullTimestamp:   false,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	// 先验证新轮转文件可写，再替换旧输出；失败时仍保留终端和 TUI 镜像。
	var logWriter *rotatelogs.RotateLogs
	var fileErr error
	if LogToFile {
		filename := path.Join(LogFilePath, LogFileName)
		logWriter, fileErr = rotatelogs.New(filename+"%Y%m%d.log",
			rotatelogs.WithLinkName(filename),
			rotatelogs.WithMaxAge(7*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour))
		if fileErr == nil {
			_, fileErr = logWriter.Write(nil)
			if fileErr != nil {
				_ = logWriter.Close()
				logWriter = nil
			}
		}
	}
	outputMu.Lock()
	oldFile := fileOutput
	fileOutput = logWriter
	applyOutput()
	// logrus.SetOutput 已等待旧输出的写入完成，此时可安全关闭旧文件。
	if oldFile != nil {
		_ = oldFile.Close()
	}
	outputMu.Unlock()
	if fileErr != nil {
		logger.Info("logger err:", fileErr)
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
			// 错误响应可能由外层 Echo 处理，日志先按返回错误确定状态。
			if err != nil && !c.Response().Committed {
				statusCode = http.StatusInternalServerError
				var httpErr *echo.HTTPError
				if errors.As(err, &httpErr) {
					statusCode = httpErr.Code
				}
			}
			// 高频状态轮询成功时保持安静；失败和其他接口仍输出到所有日志通道。
			route := c.Path()
			polling := false
			for _, endpoint := range []string{"/api/server", "/api/server/traffic", "/api/connections", "/api/configs", "/api/configs/status"} {
				if strings.HasSuffix(route, endpoint) {
					polling = true
					break
				}
			}
			// 需要查看成功轮询日志时，临时取消下一行注释，让 Debug 模式放行。
			// polling = polling && !Debug
			if reqMethod == http.MethodGet && polling && err == nil && statusCode >= 200 && statusCode < 300 {
				return nil
			}
			logMsg := fmt.Sprintf("[%s:%d][%6.2fms][%s]%s",
				reqMethod,
				statusCode,
				latencyTime,
				c.RealIP(),
				reqURI,
			)
			// RealIP 与 RequestURI 都来自请求，写入日志面板 HTML 前必须转义。
			logMsgWeb := fmt.Sprintf("[%s:%d]<span style=\"color:#d08700\">[%6.2fms]</span><span style=\"color:#0084d1\">[%s]</span>%s",
				reqMethod,
				statusCode,
				latencyTime,
				html.EscapeString(c.RealIP()),
				html.EscapeString(reqURI),
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
