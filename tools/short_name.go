package tools

import "regexp"

var (
	// 只删除结尾处的常见扩展名（忽略大小写）。
	shortNameExtReg = regexp.MustCompile(`\.(?i)(zip|rar|cbr|cbz|tar|pdf|mp3|mp4|flv|gz|webm|gif|png|jpg|jpeg|webp|svg|psd|bmp|tif)$`)

	// 去除各种括号及其内容；这些规则保持分开，避免不同括号类型互相吞掉内容。
	shortNameRoundReg         = regexp.MustCompile(`\([^()]*?\)`)  // 匹配 ()
	shortNameSquareReg        = regexp.MustCompile(`\[[^\[\]]*?]`) // 匹配 []
	shortNameChineseRoundReg  = regexp.MustCompile(`（[^（）]*?）`)    // 匹配 （）
	shortNameChineseSquareReg = regexp.MustCompile(`【[^【】]*?】`)    // 匹配 【】

	// 只移除开头域名，保留正文中可能属于标题内容的域名字符串。
	shortNameDomainReg = regexp.MustCompile(`^(((ht|f)tps?)://)?([^!@#$%^&*?.\s-]([^!@#$%^&*?.\s]{0,63}[^!@#$%^&*?.\s])?\.)+[a-zA-Z]{2,6}/?`)

	shortNameLeadingSpaceReg  = regexp.MustCompile(`^\s+`)
	shortNameTrailingSpaceReg = regexp.MustCompile(`\s+$`)

	// 去除开头的一连串标点符号，但不处理标题中间的标点。
	shortNameLeadingPunctuationReg = regexp.MustCompile(`^[\-` + "`" + `~!@#$^&*=|{}':;'@#￥……&*——|{}‘；：”“'。，、？]+`)
)

// ShortName 把文件名或书名清洗成书架上展示用的短标题。
func ShortName(title string) string {
	const (
		minReadableRunes = 2
		maxDisplayRunes  = 15
	)

	shortTitle := title

	// 下面的清洗顺序会影响结果，尤其是先删扩展名、括号内容，再删域名和开头标点。
	// 因此这里保留直线流水线，避免规则调整时不小心改变已有书名显示结果。

	// 1. 移除常见文件扩展名。
	shortTitle = shortNameExtReg.ReplaceAllString(shortTitle, "")

	// 2. 顺序移除所有括号及内部描述。
	shortTitle = shortNameRoundReg.ReplaceAllString(shortTitle, "")
	shortTitle = shortNameSquareReg.ReplaceAllString(shortTitle, "")
	shortTitle = shortNameChineseRoundReg.ReplaceAllString(shortTitle, "")
	shortTitle = shortNameChineseSquareReg.ReplaceAllString(shortTitle, "")

	// 3. 移除开头域名。
	shortTitle = shortNameDomainReg.ReplaceAllString(shortTitle, "")

	// 4. 去除首尾空白。
	shortTitle = shortNameLeadingSpaceReg.ReplaceAllString(shortTitle, "")
	shortTitle = shortNameTrailingSpaceReg.ReplaceAllString(shortTitle, "")

	// 5. 去除开头标点，再次去除可能被留下的首尾空白。
	shortTitle = shortNameLeadingPunctuationReg.ReplaceAllString(shortTitle, "")
	shortTitle = shortNameLeadingSpaceReg.ReplaceAllString(shortTitle, "")
	shortTitle = shortNameTrailingSpaceReg.ReplaceAllString(shortTitle, "")

	// 转成 rune 后再判断长度，避免中文、日文等多字节标题被按 byte 截断。
	cleanRunes := []rune(shortTitle)
	if len(cleanRunes) >= minReadableRunes {
		// 清洗后的标题足够可读时，优先使用清洗结果。
		if len(cleanRunes) <= maxDisplayRunes {
			return shortTitle
		}
		// 清洗标题过长时只截取展示长度，不回退原标题。
		return string(cleanRunes[:maxDisplayRunes]) + "…"
	}

	// 清洗后不足 2 个字符时，说明括号组、域名或标点规则可能删掉了主要内容；
	// 此时回退到原标题，避免书架上出现空白或只有一个字符的标题。
	originalRunes := []rune(title)
	if len(originalRunes) == 0 {
		return ""
	}
	if len(originalRunes) <= maxDisplayRunes {
		return string(originalRunes)
	}

	// 只有原标题超过展示长度时才追加省略号；刚好 15 个 rune 时完整展示。
	return string(originalRunes[:maxDisplayRunes]) + "…"
}
