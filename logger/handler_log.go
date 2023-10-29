package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"time"
)

// HandlerLog 默认log
func HandlerLog(LogToFile bool, LogFilePath string, LogFileName string) gin.HandlerFunc {
	Log.SetLevel(logrus.DebugLevel) //设置最低的日志级别
	//Log.SetReportCaller(true)           //开启返回函数名和行号
	Log.SetFormatter(&MyFormatter{}) // 使用自定义格式化程序

	//设置输出
	if LogToFile {
		//日志文件路径
		filename := path.Join(LogFilePath, LogFileName)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("Log err:", err)
		}
		//设置多种输出类型(默认值“os.Stderr”)
		Log.SetOutput(io.MultiWriter(os.Stdout, file))
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
			fmt.Println("Log err:", err)
		}
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}
		Log.AddHook(lfshook.NewHook(writeMap, &logrus.TextFormatter{
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
		latencyTime := endTime.Sub(startTime).String()
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()
		//if len(reqURI) > 40 {
		//	reqURI = reqURI[:40] + "..."
		//}
		// 日志格式
		//%15s 是一个格式说明符，用于格式化字符串。这里的 15 表示输出的字符串宽度应该为至少15个字符。
		//使用 %12s 使得字符串 "test" 右对齐，并在其左侧添加了空格来达到15个字符的宽度。
		//使用 %-12s 使得字符串左对齐，并在其右侧添加了空格来达到15个字符的宽度
		Log.WithFields(logrus.Fields{
			//"status_code":  statusCode,
			//"client_ip":    c.ClientIP(),
		}).Info(fmt.Sprintf("%s:%3d|%13v|%s%s\n", reqMethod, statusCode, latencyTime, c.ClientIP(), reqURI))
	}
}

func LoggerToStdout() gin.HandlerFunc {
	//设置日志级别
	Log.SetLevel(logrus.DebugLevel)
	//设置输出
	Log.SetOutput(os.Stdout)
	//自定义gin处理函数
	return func(c *gin.Context) {
		//开始时间
		startTime := time.Now()
		//处理请求
		c.Next()
		//结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime).String()
		//请求方式
		reqMethod := c.Request.Method
		//请求路由
		reqUrl := c.Request.RequestURI
		//状态码
		statusCode := c.Writer.Status()
		//请求ip
		clientIP := c.ClientIP()
		// 日志格式
		Log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqUrl,
		}).Info()
	}
}
