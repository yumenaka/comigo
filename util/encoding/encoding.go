package encoding

//source: github.com/mholt/archiver/pull/149/files/92cf5d0fb45d7fa943e25fc83fc71cd2e734a4fb
import (
	"errors"
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

var encodings = map[string]encoding.Encoding{
	"ibm866":            charmap.CodePage866,
	"iso8859_2":         charmap.ISO8859_2,
	"iso8859_3":         charmap.ISO8859_3,
	"iso8859_4":         charmap.ISO8859_4,
	"iso8859_5":         charmap.ISO8859_5,
	"iso8859_6":         charmap.ISO8859_6,
	"iso8859_7":         charmap.ISO8859_7,
	"iso8859_8":         charmap.ISO8859_8,
	"iso8859_8I":        charmap.ISO8859_8I,
	"iso8859_10":        charmap.ISO8859_10,
	"iso8859_13":        charmap.ISO8859_13,
	"iso8859_14":        charmap.ISO8859_14,
	"iso8859_15":        charmap.ISO8859_15,
	"iso8859_16":        charmap.ISO8859_16,
	"koi8r":             charmap.KOI8R,
	"koi8u":             charmap.KOI8U,
	"macintosh":         charmap.Macintosh,
	"windows874":        charmap.Windows874,
	"windows1250":       charmap.Windows1250,
	"windows1251":       charmap.Windows1251,
	"windows1252":       charmap.Windows1252,
	"windows1253":       charmap.Windows1253,
	"windows1254":       charmap.Windows1254,
	"windows1255":       charmap.Windows1255,
	"windows1256":       charmap.Windows1256,
	"windows1257":       charmap.Windows1257,
	"windows1258":       charmap.Windows1258,
	"macintoshcyrillic": charmap.MacintoshCyrillic,
	"gbk":               simplifiedchinese.GBK,
	"gb18030":           simplifiedchinese.GB18030,
	"big5":              traditionalchinese.Big5,
	"eucjp":             japanese.EUCJP,
	"iso2022jp":         japanese.ISO2022JP,
	"shiftjis":          japanese.ShiftJIS,
	"euckr":             korean.EUCKR,
	"utf16be":           unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM),
	"utf16le":           unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM),
}

func GuessText(unknowString string) (string, error) {
	//var GuestList = []string{"gbk", "gb18030", "big5", "eucjp", "iso2022jp", "shiftjis", "euckr", "utf16be", "utf16le"}
	if isGBK([]byte(unknowString)) {
		utfString, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(unknowString)) //将gbk转换为utf-8
		return string(utfString), err
	} else {
		utfString, err := japanese.ShiftJIS.NewDecoder().Bytes([]byte(unknowString)) //将ShiftJIS转换为utf-8
		return string(utfString), err
	}
}

func ShiftjisToUtf8(unknowString string) (string, error) {
	utfString, err := japanese.ShiftJIS.NewDecoder().Bytes([]byte(unknowString)) //将ShiftJIS转换为utf-8
	return string(utfString), err
}

func GbkToUtf8(unknowString string) (string, error) {
	utfString, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(unknowString)) //将gbk转换为utf-8
	return string(utfString), err
}

// ToShiftJIS Convert a string encoding from UTF-8 to ShiftJIS
func ToShiftJIS(str string) (string, error) {
	return transformEncoding(strings.NewReader(str), japanese.ShiftJIS.NewEncoder())
}

// ToGBK Convert a string encoding from UTF-8 to ShiftJIS
func ToGBK(str string) (string, error) {
	return transformEncoding(strings.NewReader(str), simplifiedchinese.GBK.NewEncoder())
}

func transformEncoding(rawReader io.Reader, trans transform.Transformer) (string, error) {
	ret, err := io.ReadAll(transform.NewReader(rawReader, trans))
	if err == nil {
		return string(ret), nil
	} else {
		return "", err
	}
}

func isGBK(data []byte) bool {
	length := len(data)
	var i = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有一个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//大于127的使用双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

func GetEncoding(charset string) (encoding.Encoding, bool) {
	charset = strings.ToLower(charset)
	enc, ok := encodings[charset]
	return enc, ok
}

func Decode(in []byte, charset string) ([]byte, error) {
	if enc, ok := GetEncoding(charset); ok {
		return enc.NewDecoder().Bytes(in)
	}
	return nil, errors.New("charset not found")
}

func DecodeFileName(headerName string, ZipFilenameEncoding string) string {
	if ZipFilenameEncoding != "" { //common.Config.ZipFileTextEncoding
		if filename, err := Decode([]byte(headerName), ZipFilenameEncoding); err == nil {
			return string(filename)
		}
	}
	return headerName
}

func GetEncodingByName(charset string) encoding.Encoding {
	charset = strings.ToLower(charset)
	enc, ok := encodings[charset]
	if !ok {
		return charmap.Windows1252
	} else {
		return enc
	}
}
