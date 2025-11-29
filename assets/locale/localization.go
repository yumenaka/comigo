package locale

import (
	_ "embed"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/example_test.go

var Localizer *i18n.Localizer

//go:embed en_US.json
var enBytes []byte

//go:embed zh_CN.json
var cnBytes []byte

//go:embed ja_JP.json
var jpBytes []byte

func getLocale() (string, string) {
	osHost := runtime.GOOS
	defaultLang := "en"
	defaultLoc := "US"
	switch osHost {
	case "windows":
		// 在 Windows 上通过 PowerShell 获取当前 UI 语言
		// 使用 Culture 的 Name（如 zh-CN、en-US、ja-JP），而不是 Title（如 "Chinese (Simplified, China)"）
		// 避免在 Windows 11 等系统上解析失败
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "[System.Globalization.CultureInfo]::CurrentUICulture.Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			if langLocRaw != "" {
				// 统一处理可能出现的后缀和分隔符：
				// 例如：zh-CN, en-US, ja-JP, zh-Hans-CN
				// 1. 去掉编码后缀（如果有的话），例如 zh-CN.UTF-8
				langLocRaw = strings.Split(langLocRaw, ".")[0]
				langLocRaw = strings.TrimSpace(langLocRaw)
				// 2. 将 - 统一替换为 _，方便按 "_" 切分
				langLocRaw = strings.ReplaceAll(langLocRaw, "-", "_")

				langLoc := strings.Split(langLocRaw, "_")
				lang := langLoc[0]
				loc := lang
				// Windows 上可能出现 zh_CN、zh_Hans_CN 等形式，这里统一取最后一段作为地区信息
				if len(langLoc) > 1 {
					loc = langLoc[len(langLoc)-1]
				}
				return lang, loc
			}
		}
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.Split(envlang, ".")[0]
			langLocRaw = strings.TrimSpace(langLocRaw)
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := lang
			if len(langLoc) > 1 {
				loc = langLoc[1]
			}
			return lang, loc
		}
	case "darwin":
		// 首先尝试从 LANG 环境变量获取（保持向后兼容）
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.Split(envlang, ".")[0]
			langLocRaw = strings.TrimSpace(langLocRaw)
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := lang
			if len(langLoc) > 1 {
				loc = langLoc[1]
			}
			return lang, loc
		}
		// 如果 LANG 不存在，使用 defaults 命令获取 macOS 系统语言设置
		cmd := exec.Command("defaults", "read", "-g", "AppleLanguages")
		output, err := cmd.Output()
		if err == nil {
			// defaults 返回的格式通常是: ("zh-Hans-CN", "en", ...)
			// 或者: (zh-Hans-CN, en, ...)
			outputStr := strings.TrimSpace(string(output))
			// 移除括号和引号
			outputStr = strings.Trim(outputStr, "()")
			// 按逗号分割
			languages := strings.Split(outputStr, ",")
			if len(languages) > 0 {
				// 获取第一个语言代码
				firstLang := strings.TrimSpace(languages[0])
				// 移除引号
				firstLang = strings.Trim(firstLang, "\"'")
				// 解析语言代码，格式可能是: zh-Hans-CN, en, ja-JP 等
				langParts := strings.Split(firstLang, "-")
				lang := langParts[0]
				loc := lang
				// 尝试从语言代码中提取地区信息
				// 例如: zh-Hans-CN -> lang=zh, loc=CN
				// 例如: en-US -> lang=en, loc=US
				if len(langParts) > 1 {
					// 检查最后一部分是否是地区代码（通常是2-3个大写字母）
					lastPart := langParts[len(langParts)-1]
					if len(lastPart) >= 2 && len(lastPart) <= 3 && strings.ToUpper(lastPart) == lastPart {
						loc = lastPart
					} else if len(langParts) >= 2 {
						// 如果最后一部分不是地区代码，使用第二个部分
						loc = langParts[1]
					}
				}
				return lang, loc
			}
		}
	}
	return defaultLang, defaultLoc
}

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustParseMessageFileBytes(enBytes, "en_US.json")
	bundle.MustParseMessageFileBytes(cnBytes, "zh_CN.json")
	bundle.MustParseMessageFileBytes(jpBytes, "ja_JP.json")

	// 默认使用自动检测的语言
	lang, _ := getLocale()
	// logger.Infof("OK: language=%s, locale=%s\n", lang, loc)
	switch lang {
	case "zh":
		Localizer = i18n.NewLocalizer(bundle, "zh-CN")
	case "en":
		Localizer = i18n.NewLocalizer(bundle, "en-US")
	case "ja":
		Localizer = i18n.NewLocalizer(bundle, "ja-JP")
	default:
		Localizer = i18n.NewLocalizer(bundle, "en-US")
	}
}

// InitLanguageFromConfig 根据配置初始化语言设置
// 如果配置的 Language 为 "auto"，则使用自动检测的语言
// 否则使用配置指定的语言
func InitLanguageFromConfig(configLanguage string) {
	if configLanguage == "" || configLanguage == "auto" {
		// 使用自动检测的语言
		lang, _ := getLocale()
		switch lang {
		case "zh":
			Localizer = i18n.NewLocalizer(bundle, "zh-CN")
		case "en":
			Localizer = i18n.NewLocalizer(bundle, "en-US")
		case "ja":
			Localizer = i18n.NewLocalizer(bundle, "ja-JP")
		default:
			Localizer = i18n.NewLocalizer(bundle, "en-US")
		}
	} else {
		// 使用配置指定的语言
		switch configLanguage {
		case "zh", "zh-CN":
			Localizer = i18n.NewLocalizer(bundle, "zh-CN")
		case "en", "en-US":
			Localizer = i18n.NewLocalizer(bundle, "en-US")
		case "ja", "ja-JP":
			Localizer = i18n.NewLocalizer(bundle, "ja-JP")
		default:
			// 如果配置的语言不支持，使用自动检测的语言
			lang, _ := getLocale()
			switch lang {
			case "zh":
				Localizer = i18n.NewLocalizer(bundle, "zh-CN")
			case "en":
				Localizer = i18n.NewLocalizer(bundle, "en-US")
			case "ja":
				Localizer = i18n.NewLocalizer(bundle, "ja-JP")
			default:
				Localizer = i18n.NewLocalizer(bundle, "en-US")
			}
		}
	}
}

func GetString(id string) string {
	return Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: id})
}

func GetStringByLocal(id string, local string) string {
	return i18n.NewLocalizer(bundle, local).MustLocalize(&i18n.LocalizeConfig{MessageID: id})
}

// SetLanguage 设置当前语言
func SetLanguage(lang string) error {
	switch lang {
	case "zh-CN", "zh":
		Localizer = i18n.NewLocalizer(bundle, "zh-CN")
	case "en-US", "en":
		Localizer = i18n.NewLocalizer(bundle, "en-US")
	case "ja-JP", "ja":
		Localizer = i18n.NewLocalizer(bundle, "ja-JP")
	default:
		Localizer = i18n.NewLocalizer(bundle, "en-US")
	}
	return nil
}
