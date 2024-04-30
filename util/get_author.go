package util

import (
	"strings"
)

func GetAuthor(input string) string {
	pairs := map[rune]rune{
		'[': ']', //这里的中括号是半角的
		'［': '］', //这里的中括号是全角的
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
				return input[1:closeIndex]
			}
		}
	}

	pairsError := map[rune]rune{
		'[': '］', //半角——全角
		'［': ']', //全角——半角
	}
	for open, closed := range pairsError {
		if strings.HasPrefix(input, string(open)) {
			closeIndex := strings.IndexRune(input, closed)
			if closeIndex != -1 {
				return input[1:closeIndex]
			}
		}
	}
	return ""
}
