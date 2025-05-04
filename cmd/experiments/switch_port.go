package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

var (
	server *http.Server
	mutex  sync.Mutex
)

// startServer 启动一个在指定端口监听的HTTP服务器
func startServer(port string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Listening on port "+port)
	})

	mutex.Lock()
	server = &http.Server{
		Addr:    ":" + port,
		Handler: e,
	}
	mutex.Unlock()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}

// stopServer 停止当前的HTTP服务器
func stopServer() error {
	mutex.Lock()
	defer mutex.Unlock()
	if server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	server = nil
	return nil
}

// switchPort 停止当前的服务器并在新的端口上重新启动
func switchPort(newPort string) {
	if err := stopServer(); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server Shutdown Successfully", "Starting Server...", "on port", newPort, "...")
	go startServer(newPort)
}

func main() {
	// 在端口8080上启动服务器
	go startServer("18080")

	// 假设在一段时间后需要切换到端口8081
	time.Sleep(20 * time.Second)
	switchPort("18081")

	// 为了示例，我们在这里阻塞主线程
	// 在实际应用中，需要一个更复杂的逻辑来决定何时停止程序
	select {}
}
