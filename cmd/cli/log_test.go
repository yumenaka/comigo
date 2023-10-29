package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 自定义logrus输出格式
func setupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", // 使用完整时间戳
	})
}

// 自定义Gin的日志中间件，使用logrus
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := c.Request.Method
		reqURI := c.Request.RequestURI
		statusCode := c.Writer.Status()

		logrus.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    c.ClientIP(),
			"req_method":   reqMethod,
			"req_uri":      reqURI,
		}).Info()
	}
}

func main() {
	setupLogger()

	r := gin.New()
	r.Use(LoggerMiddleware())

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	r.Run(":8080")
}
