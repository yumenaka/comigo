package tools

import (
	"strings"
)

func GetAuthor(input string) (string, bool) {
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

	for open, close := range pairs {
		if strings.HasPrefix(input, string(open)) {
			closeIndex := strings.IndexRune(input, close)
			if closeIndex != -1 {
				return input[1:closeIndex], true
			}
		}
	}

	pairs_error := map[rune]rune{
		'[': '］', //半角——全角
		'［': ']', //全角——半角
	}
	for open, close := range pairs_error {
		if strings.HasPrefix(input, string(open)) {
			closeIndex := strings.IndexRune(input, close)
			if closeIndex != -1 {
				return input[1:closeIndex], true
			}
		}
	}
	return "", false
}
