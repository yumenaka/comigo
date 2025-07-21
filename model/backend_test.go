package model

import (
	"testing"
)

func TestParseFileBackendURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *FileBackend
		wantErr bool
	}{
		{
			name: "本地文件路径 - Unix风格",
			url:  "/home/user/books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "/home/user/books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - Windows风格",
			url:  "C:\\Users\\user\\books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "C:\\Users\\user\\books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - Windows风格 D盘",
			url:  "D:\\Documents\\books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "D:\\Documents\\books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - Windows风格 E盘",
			url:  "E:\\Data\\books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "E:\\Data\\books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - Windows风格 正斜杠",
			url:  "C:/Users/user/books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "C:/Users/user/books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - Windows风格 网络路径",
			url:  "\\\\server\\share\\books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "\\\\server\\share\\books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - file://scheme",
			url:  "file:///home/user/books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "/home/user/books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - file://scheme Windows",
			url:  "file:///C:/Users/user/books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "C:/Users/user/books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - file://scheme Windows 反斜杠",
			url:  "file:///C:\\Users\\user\\books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "C:\\Users\\user\\books",
			},
			wantErr: false,
		},
		{
			name: "本地文件路径 - file://scheme 相对路径",
			url:  "file://./books",
			want: &FileBackend{
				Type: LocalDisk,
				URL:  "./books",
			},
			wantErr: false,
		},
		{
			name: "SMB URL - 无认证",
			url:  "smb://server/share/folder",
			want: &FileBackend{
				Type:         SMB,
				URL:          "smb://server/share/folder",
				ServerHost:   "server",
				ServerPort:   445,
				NeedAuth:     false,
				SMBShareName: "share",
				SMBPath:      "folder",
			},
			wantErr: false,
		},
		{
			name: "SMB URL - 带认证",
			url:  "smb://workgroup;user:password@server/share/folder",
			want: &FileBackend{
				Type:         SMB,
				URL:          "smb://workgroup;user:password@server/share/folder",
				ServerHost:   "server",
				ServerPort:   445,
				NeedAuth:     true,
				AuthUsername: "user",
				AuthPassword: "password",
				SMBShareName: "share",
				SMBPath:      "folder",
			},
			wantErr: false,
		},
		{
			name: "SFTP URL - 无认证",
			url:  "sftp://host/path",
			want: &FileBackend{
				Type:       SFTP,
				URL:        "sftp://host/path",
				ServerHost: "host",
				ServerPort: 22,
				NeedAuth:   false,
			},
			wantErr: false,
		},
		{
			name: "SFTP URL - 带认证和端口",
			url:  "sftp://user:password@host:2222/path",
			want: &FileBackend{
				Type:         SFTP,
				URL:          "sftp://user:password@host:2222/path",
				ServerHost:   "host",
				ServerPort:   2222,
				NeedAuth:     true,
				AuthUsername: "user",
				AuthPassword: "password",
			},
			wantErr: false,
		},
		{
			name: "WebDAV URL - HTTP",
			url:  "http://host/path",
			want: &FileBackend{
				Type:       WebDAV,
				URL:        "http://host/path",
				ServerHost: "host",
				ServerPort: 80,
				NeedAuth:   false,
			},
			wantErr: false,
		},
		{
			name: "WebDAV URL - HTTPS带认证",
			url:  "https://user:password@host:8443/path",
			want: &FileBackend{
				Type:         WebDAV,
				URL:          "https://user:password@host:8443/path",
				ServerHost:   "host",
				ServerPort:   8443,
				NeedAuth:     true,
				AuthUsername: "user",
				AuthPassword: "password",
			},
			wantErr: false,
		},
		{
			name: "FTP URL - 无认证",
			url:  "ftp://host/dir",
			want: &FileBackend{
				Type:       FTP,
				URL:        "ftp://host/dir",
				ServerHost: "host",
				ServerPort: 21,
				NeedAuth:   false,
			},
			wantErr: false,
		},
		{
			name: "FTPS URL - 带认证和端口",
			url:  "ftps://user:password@host:990/dir",
			want: &FileBackend{
				Type:         FTP,
				URL:          "ftps://user:password@host:990/dir",
				ServerHost:   "host",
				ServerPort:   990,
				NeedAuth:     true,
				AuthUsername: "user",
				AuthPassword: "password",
			},
			wantErr: false,
		},
		{
			name: "S3 URL",
			url:  "s3://endpoint/bucket/prefix",
			want: &FileBackend{
				Type:       S3,
				URL:        "s3://endpoint/bucket/prefix",
				ServerHost: "endpoint",
				ServerPort: 443,
				NeedAuth:   false,
			},
			wantErr: false,
		},
		{
			name: "S3 URL - 带端口",
			url:  "s3://endpoint:9000/bucket",
			want: &FileBackend{
				Type:       S3,
				URL:        "s3://endpoint:9000/bucket",
				ServerHost: "endpoint",
				ServerPort: 9000,
				NeedAuth:   false,
			},
			wantErr: false,
		},
		{
			name:    "空URL",
			url:     "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "无效的URL",
			url:     "invalid://url",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "不支持的scheme",
			url:     "ssh://host/path",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFileBackendURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFileBackendURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// 比较关键字段
			if got.Type != tt.want.Type {
				t.Errorf("ParseFileBackendURL() Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.ServerHost != tt.want.ServerHost {
				t.Errorf("ParseFileBackendURL() ServerHost = %v, want %v", got.ServerHost, tt.want.ServerHost)
			}
			if got.ServerPort != tt.want.ServerPort {
				t.Errorf("ParseFileBackendURL() ServerPort = %v, want %v", got.ServerPort, tt.want.ServerPort)
			}
			if got.NeedAuth != tt.want.NeedAuth {
				t.Errorf("ParseFileBackendURL() NeedAuth = %v, want %v", got.NeedAuth, tt.want.NeedAuth)
			}
			if got.AuthUsername != tt.want.AuthUsername {
				t.Errorf("ParseFileBackendURL() AuthUsername = %v, want %v", got.AuthUsername, tt.want.AuthUsername)
			}
			if got.AuthPassword != tt.want.AuthPassword {
				t.Errorf("ParseFileBackendURL() AuthPassword = %v, want %v", got.AuthPassword, tt.want.AuthPassword)
			}
			if got.SMBShareName != tt.want.SMBShareName {
				t.Errorf("ParseFileBackendURL() SMBShareName = %v, want %v", got.SMBShareName, tt.want.SMBShareName)
			}
			if got.SMBPath != tt.want.SMBPath {
				t.Errorf("ParseFileBackendURL() SMBPath = %v, want %v", got.SMBPath, tt.want.SMBPath)
			}
		})
	}
}
