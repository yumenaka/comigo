package tools

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gofrs/flock"
	"github.com/yumenaka/comigo/tools/logger"
)

// SingleInstance 单实例管理器
type SingleInstance struct {
	lockFile   *flock.Flock
	socketPath string
	listener   net.Listener
	onNewArgs  func(args []string) error
}

// Message 进程间通信消息
type Message struct {
	Args []string `json:"args"`
}

var (
	globalInstance *SingleInstance
)

// NewSingleInstance 创建单实例管理器
func NewSingleInstance(onNewArgs func(args []string) error) (*SingleInstance, error) {
	// 获取配置目录或临时目录
	var configDir string
	home, err := os.UserHomeDir()
	if err != nil {
		// 如果获取 Home 失败，使用临时目录
		configDir = os.TempDir()
	} else {
		// 使用 ~/.config/comigo 作为配置目录
		configDir = filepath.Join(home, ".config", "comigo")
		// 确保目录存在
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			logger.Infof("Failed to create config dir, using temp dir: %v", err)
			configDir = os.TempDir()
		}
	}

	// 创建锁文件路径
	lockFilePath := filepath.Join(configDir, "comigo.lock")

	// 创建 socket 路径（跨平台）
	var socketPath string
	if runtime.GOOS == "windows" {
		// Windows 使用 TCP localhost（固定端口）
		socketPath = "127.0.0.1:12346"
	} else {
		// Unix/Linux/macOS 使用 Unix Domain Socket
		socketPath = filepath.Join(configDir, "comigo.sock")
	}

	si := &SingleInstance{
		lockFile:   flock.New(lockFilePath),
		socketPath: socketPath,
		onNewArgs:  onNewArgs,
	}

	return si, nil
}

// TryLock 尝试获取锁，如果已有实例运行则返回 false
func (si *SingleInstance) TryLock() (bool, error) {
	locked, err := si.lockFile.TryLock()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}
	if !locked {
		return false, nil
	}
	return true, nil
}

// Unlock 释放锁
func (si *SingleInstance) Unlock() error {
	if si.lockFile != nil {
		return si.lockFile.Unlock()
	}
	return nil
}

// StartServer 启动服务器，监听来自其他实例的消息
func (si *SingleInstance) StartServer() error {
	var err error
	if runtime.GOOS == "windows" {
		// Windows 使用 TCP localhost
		si.listener, err = net.Listen("tcp", si.socketPath)
		if err != nil {
			return fmt.Errorf("failed to create TCP listener: %w", err)
		}
	} else {
		// Unix/Linux/macOS 使用 Unix Domain Socket
		// 确保 socket 文件不存在
		if _, err := os.Stat(si.socketPath); err == nil {
			os.Remove(si.socketPath)
		}
		si.listener, err = net.Listen("unix", si.socketPath)
		if err != nil {
			return fmt.Errorf("failed to create unix socket: %w", err)
		}
		// 设置 socket 文件权限
		if err := os.Chmod(si.socketPath, 0600); err != nil {
			logger.Infof("Warning: failed to set socket permissions: %v", err)
		}
	}

	logger.Infof("Single instance server started on: %s", si.socketPath)

	// 在 goroutine 中处理连接
	go func() {
		for {
			conn, err := si.listener.Accept()
			if err != nil {
				// 如果 listener 已关闭，正常退出
				if opErr, ok := err.(*net.OpError); ok {
					if opErr.Err.Error() == "use of closed network connection" || opErr.Err.Error() == "operation on closed network connection" {
						return
					}
				}
				logger.Infof("Failed to accept connection: %v", err)
				continue
			}
			go si.handleConnection(conn)
		}
	}()

	return nil
}

// handleConnection 处理来自其他实例的连接
func (si *SingleInstance) handleConnection(conn net.Conn) {
	defer conn.Close()

	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// 读取消息
	var msg Message
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&msg); err != nil {
		logger.Infof("Failed to decode message: %v", err)
		return
	}

	// 处理新参数
	if si.onNewArgs != nil && len(msg.Args) > 0 {
		if err := si.onNewArgs(msg.Args); err != nil {
			logger.Infof("Failed to handle new args: %v", err)
			// 发送错误响应
			encoder := json.NewEncoder(conn)
			encoder.Encode(map[string]string{"error": err.Error()})
			return
		}
		// 发送成功响应
		encoder := json.NewEncoder(conn)
		encoder.Encode(map[string]string{"status": "ok"})
		logger.Infof("Received and processed new args: %v", msg.Args)
	}
}

// SendArgs 向已运行的实例发送参数
func (si *SingleInstance) SendArgs(args []string) error {
	var conn net.Conn
	var err error

	if runtime.GOOS == "windows" {
		// Windows: 使用 TCP localhost
		conn, err = net.Dial("tcp", si.socketPath)
		if err != nil {
			return fmt.Errorf("failed to connect to existing instance: %w", err)
		}
	} else {
		// Unix/Linux/macOS: 连接 Unix Domain Socket
		conn, err = net.Dial("unix", si.socketPath)
		if err != nil {
			return fmt.Errorf("failed to connect to existing instance: %w", err)
		}
	}
	defer conn.Close()

	// 设置写入超时
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))

	// 发送消息
	msg := Message{Args: args}
	encoder := json.NewEncoder(conn)
	if err := encoder.Encode(&msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	// 读取响应
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var response map[string]string
	decoder := json.NewDecoder(conn)
	if err := decoder.Decode(&response); err != nil {
		// 即使读取响应失败，消息可能已经发送成功
		logger.Infof("Failed to read response, but message may have been sent: %v", err)
		return nil
	}

	if errMsg, ok := response["error"]; ok {
		return fmt.Errorf("instance returned error: %s", errMsg)
	}

	logger.Infof("Successfully sent args to existing instance: %v", args)
	return nil
}

// Stop 停止服务器并清理资源
func (si *SingleInstance) Stop() error {
	if si.listener != nil {
		if err := si.listener.Close(); err != nil {
			logger.Infof("Error closing listener: %v", err)
		}
		si.listener = nil
	}

	// 清理 socket 文件（Unix 系统）
	if runtime.GOOS != "windows" && si.socketPath != "" {
		if _, err := os.Stat(si.socketPath); err == nil {
			os.Remove(si.socketPath)
		}
	}

	return si.Unlock()
}

// EnsureSingleInstance 确保程序以单实例模式运行
// 如果已有实例运行，将 args 发送给已运行的实例并返回 false
// 如果是第一个实例，返回 true
func EnsureSingleInstance(args []string, onNewArgs func(args []string) error) (bool, error) {
	si, err := NewSingleInstance(onNewArgs)
	if err != nil {
		return false, fmt.Errorf("failed to create single instance manager: %w", err)
	}

	// 尝试获取锁
	locked, err := si.TryLock()
	if err != nil {
		return false, err
	}

	if !locked {
		// 已有实例运行，发送参数
		logger.Infof("Another instance is already running, sending args to it...")
		if err := si.SendArgs(args); err != nil {
			return false, fmt.Errorf("failed to send args to existing instance: %w", err)
		}
		return false, nil
	}

	// 第一个实例，启动服务器
	globalInstance = si
	if err := si.StartServer(); err != nil {
		si.Unlock()
		return false, fmt.Errorf("failed to start single instance server: %w", err)
	}

	return true, nil
}

// CleanupSingleInstance 清理单实例资源（在程序退出时调用）
func CleanupSingleInstance() {
	if globalInstance != nil {
		globalInstance.Stop()
		globalInstance = nil
	}
}
