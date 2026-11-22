package vfs

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/jlaffaye/ftp"
)

// TestFTPReadHoldsLockUntilCompletion 用本地协议替身分离数据传输和结束响应，检查整段锁范围。
func TestFTPReadHoldsLockUntilCompletion(t *testing.T) {
	for _, completion := range []string{"226 transfer complete", "451 transfer failed"} {
		t.Run(completion, func(t *testing.T) {
			control, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer control.Close()
			dataListener, err := net.Listen("tcp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}
			defer dataListener.Close()
			reading, closing, releaseData, releaseReply := make(chan struct{}), make(chan struct{}), make(chan struct{}), make(chan struct{})
			go func() {
				conn, err := control.Accept()
				if err != nil {
					return
				}
				defer conn.Close()
				_ = conn.SetDeadline(time.Now().Add(5 * time.Second))
				reader := bufio.NewReader(conn)
				fmt.Fprint(conn, "220 ready\r\n")
				line, _ := reader.ReadString('\n')
				if strings.TrimSpace(line) != "EPSV" {
					return
				}
				fmt.Fprintf(conn, "229 Entering Extended Passive Mode (|||%d|)\r\n", dataListener.Addr().(*net.TCPAddr).Port)
				data, err := dataListener.Accept()
				if err != nil {
					return
				}
				defer data.Close()
				line, _ = reader.ReadString('\n')
				if !strings.HasPrefix(line, "RETR ") {
					return
				}
				fmt.Fprint(conn, "150 opening data\r\n")
				fmt.Fprint(data, "image bytes")
				close(reading)
				<-releaseData
				data.Close()
				close(closing)
				<-releaseReply
				fmt.Fprint(conn, completion+"\r\n")
			}()
			conn, err := ftp.Dial(control.Addr().String(), ftp.DialWithTimeout(3*time.Second))
			if err != nil {
				t.Fatal(err)
			}
			filesystem := &FTPFS{conn: conn, basePath: "/"}
			result := make(chan error, 1)
			go func() {
				data, err := filesystem.ReadFile("image.jpg")
				if err == nil && string(data) != "image bytes" {
					err = fmt.Errorf("unexpected data: %q", data)
				}
				result <- err
			}()
			select {
			case <-reading:
			case <-time.After(3 * time.Second):
				t.Fatal("数据传输未开始")
			}
			if filesystem.mu.TryLock() {
				filesystem.mu.Unlock()
				t.Error("传输尚未结束便释放了锁")
			}
			close(releaseData)
			select {
			case <-closing:
			case <-time.After(3 * time.Second):
				t.Fatal("数据连接未关闭")
			}
			if filesystem.mu.TryLock() {
				filesystem.mu.Unlock()
				t.Error("结束响应到达前便释放了锁")
			}
			// 关闭文件系统也要等待结束响应，不能插入 QUIT 中断 RETR。
			quit := make(chan error, 1)
			go func() { quit <- filesystem.Close() }()
			select {
			case err := <-quit:
				t.Fatalf("传输尚未结束即发出 QUIT: %v", err)
			case <-time.After(20 * time.Millisecond):
			}
			close(releaseReply)
			select {
			case err := <-result:
				if (err != nil) != strings.HasPrefix(completion, "451") {
					t.Fatalf("结束响应错误被忽略: %v", err)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("读取未完成")
			}
			select {
			case <-quit:
			case <-time.After(3 * time.Second):
				t.Fatal("FTP 连接未关闭")
			}
		})
	}
}
