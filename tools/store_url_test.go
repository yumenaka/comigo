package tools

import "testing"

// 验证常见书库地址能解析出正确的后端类型、主机、端口和路径。
func TestParseStoreURLCommon(t *testing.T) {
	tests := []struct {
		name     string
		storeURL string
		wantType StoreBackendType
		wantHost string
		wantPort int
		wantPath string
	}{
		{
			name:     "本地 file URL",
			storeURL: "file:///home/user/books",
			wantType: StoreBackendLocalDisk,
			wantPath: "/home/user/books",
		},
		{
			name:     "SMB 认证信息",
			storeURL: "smb://workgroup;user:password@server/share/folder",
			wantType: StoreBackendSMB,
			wantHost: "server",
			wantPort: 445,
		},
		{
			name:     "SFTP 端口",
			storeURL: "sftp://user:password@host:2222/path",
			wantType: StoreBackendSFTP,
			wantHost: "host",
			wantPort: 2222,
		},
		{
			name:     "Comigo HTTPS",
			storeURL: "https://user:password@host/path",
			wantType: StoreBackendComigo,
			wantHost: "host",
			wantPort: 443,
		},
		{
			name:     "WebDAV 显式协议",
			storeURL: "webdav://user:password@host/path",
			wantType: StoreBackendWebDAV,
			wantHost: "host",
			wantPort: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseStoreURL(tt.storeURL)
			if err != nil {
				t.Fatalf("ParseStoreURL() 返回错误: %v", err)
			}
			if info.Type != tt.wantType {
				t.Fatalf("Type = %v, 期望 %v", info.Type, tt.wantType)
			}
			if tt.wantHost != "" && info.ServerHost != tt.wantHost {
				t.Fatalf("ServerHost = %q, 期望 %q", info.ServerHost, tt.wantHost)
			}
			if tt.wantPort != 0 && info.ServerPort != tt.wantPort {
				t.Fatalf("ServerPort = %d, 期望 %d", info.ServerPort, tt.wantPort)
			}
			if tt.wantPath != "" && info.LocalPath != tt.wantPath {
				t.Fatalf("LocalPath = %q, 期望 %q", info.LocalPath, tt.wantPath)
			}
		})
	}
}

// 验证书库地址归一化会去掉认证信息并统一等价协议。
func TestNormalizeStoreURLKey(t *testing.T) {
	tests := []struct {
		storeURL string
		want     string
	}{
		{"webdav://user:pass@example.com/books", "http://example.com/books"},
		{"davs://example.com/books", "https://example.com/books"},
		{"sftp://user:pass@example.com:2222/books", "sftp://example.com:2222/books"},
		{"C:\\Users\\test\\Books", "C:\\Users\\test\\Books"},
	}

	for _, tt := range tests {
		got := NormalizeStoreURLKey(tt.storeURL)
		if got != tt.want {
			t.Fatalf("NormalizeStoreURLKey(%q) = %q, 期望 %q", tt.storeURL, got, tt.want)
		}
	}
}

// 验证以点开头的普通目录名不会被误判为父目录跳转。
func TestIsSubPathAllowsDotPrefixedNames(t *testing.T) {
	if !IsSubPath("/a/b", "/a/b/..foo") {
		t.Fatal("..foo 是普通子目录名，不应被当成父目录跳转")
	}
}
