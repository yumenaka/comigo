package tools

import (
	"strings"
)

// GetAuthor 提取书名开头括号中的作者，按括号的 UTF-8 字节长度截取。
func GetAuthor(input string) string {
	pairs := map[rune]rune{
		'[': ']', // 这里的中括号是半角的
		'［': '］', // 这里的中括号是全角的
		'【': '】',
		'「': '」',
		'『': '』',
		'(': ')',
		'（': '）',
		'{': '}',
		'＜': '＞',
		'<': '>',
		'《': '》',
		'〈': '〉',
		'〔': '〕',
		'〖': '〗',
	}

	for open, closed := range pairs {
		if strings.HasPrefix(input, string(open)) {
			closeIndex := strings.IndexRune(input, closed)
			if closeIndex != -1 {
				return input[len(string(open)):closeIndex]
			}
		}
	}

	pairsError := map[rune]rune{
		'[': '］', // 半角——全角
		'［': ']', // 全角——半角
	}
	for open, closed := range pairsError {
		if strings.HasPrefix(input, string(open)) {
			closeIndex := strings.IndexRune(input, closed)
			if closeIndex != -1 {
				return input[len(string(open)):closeIndex]
			}
		}
	}
	return ""
}
