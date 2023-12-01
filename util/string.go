package util

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/yumenaka/comi/logger"
)

//// 一个语言检测包，它告诉您某些提供的文本数据是用哪种（人类）语言编写的。 需要导入：
////go get github.com/pemistahl/lingua-go@v1.0.5
//func CheckStringLanguage(s string) string {
//	languages := []lingua.Language{
//		lingua.English,
//		lingua.Japanese,
//		lingua.Chinese,
//		lingua.French,
//		lingua.German,
//		lingua.Spanish,
//	}
//	detector := lingua.NewLanguageDetectorBuilder().
//		FromLanguages(languages...).
//		Build()
//	if language, exists := detector.DetectLanguageOf("languages are awesome"); exists {
//		logger.Info(language)
//		return language
//	}
//	return ""
//}

// DetectUTF8 检测 s 是否为有效的 UTF-8 字符串，以及该字符串是否必须被视为 UTF-8 编码（即，不兼容CP-437、ASCII 或任何其他常见编码）。
// 来自： go\src\archive\zip\reader.go
func DetectUTF8(s string) (valid, require bool) {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		// Officially, ZIP uses CP-437, but many readers use the system's
		// local character encoding. Most encoding are compatible with a large
		// subset of CP-437, which itself is ASCII-like.
		//
		// Forbid 0x7e and 0x5c since EUC-KR and Shift-JIS replace those
		// characters with localized currency and overline characters.
		if r < 0x20 || r > 0x7d || r == 0x5c {
			if !utf8.ValidRune(r) || (r == utf8.RuneError && size == 1) {
				return false, false
			}
			require = true
		}
	}
	return true, require
}

// MD5file 计算字符串MD5
func MD5file(fName string) string {
	f, e := os.Open(fName)
	if e != nil {
		logger.Info(e)
		//log.Fatal(e)
	}
	h := md5.New()
	_, e = io.Copy(h, f)
	if e != nil {
		logger.Info(e)
		//log.Fatal(e)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// 从字符串中提取数字,如果有几个数字，就简单地加起来
func getNumberFromString(s string) (int, error) {
	var err error
	num := 0
	//同时设定倒计时秒数
	valid := regexp.MustCompile("\\d+")
	numbers := valid.FindAllStringSubmatch(s, -1)
	if len(numbers) > 0 {
		//循环取出多维数组
		for _, value := range numbers {
			for _, v := range value {
				temp, errTemp := strconv.Atoi(v)
				if errTemp != nil {
					logger.Info("error num value:" + v)
				} else {
					num = num + temp
				}
			}
		}
		//logger.Info("get Number:",num," form string:",s,"numbers[]=",numbers)
	} else {
		err = errors.New("number not found")
		return 0, err
	}
	return num, err
}

// 检测字符串中是否有关键字
func haveKeyWord(checkString string, list []string) bool {
	//转换为小写，使Sketch、DOUBLE也生效
	checkString = strings.ToLower(checkString)
	for _, key := range list {
		if strings.Contains(checkString, key) {
			return true
		}
	}
	return false
}
