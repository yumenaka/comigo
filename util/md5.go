package util

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5string 计算字符串的 MD5 值
func Md5string(s string) string {
	r := md5.Sum([]byte(s))
	return hex.EncodeToString(r[:])
}
