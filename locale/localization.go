package locale

import (
	_ "embed"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/example_test.go

var (
	Localizer *i18n.Localizer
)

//go:embed en-us.toml
var enBytes []byte

//go:embed zh-cn.toml
var cnBytes []byte

//go:embed ja-jp.toml
var jpBytes []byte

func getLocale() (string, string) {
	osHost := runtime.GOOS
	defaultLang := "en"
	defaultLoc := "US"
	switch osHost {
	case "windows":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			lang := langLoc[0]
			loc := lang
			if len(langLoc) > 1 {
				loc = langLoc[1]
			}
			return lang, loc
		}
	case "darwin":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := lang
			if len(langLoc) > 1 {
				loc = langLoc[1]
			}
			return lang, loc
		}
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envlang)
			langLocRaw = strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := lang
			if len(langLoc) > 1 {
				loc = langLoc[1]
			}
			return lang, loc
		}
	}
	return defaultLang, defaultLoc
}

//func chcpToUTF8() {
//	var cmd *exec.Cmd
//	if runtime.GOOS == "windows" {
//		cmd = exec.Command("CMD", "/C", "chcp.com", "65001")
//		if err := cmd.Start(); err != nil {
//			fmt.Println("设置Windows活动代码页失败")
//			fmt.Println(err.Error())
//		}
//	}
//}

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	bundle.MustParseMessageFileBytes(enBytes, "en-us.toml")

	bundle.MustParseMessageFileBytes(cnBytes, "zh-cn.toml")

	bundle.MustParseMessageFileBytes(jpBytes, "ja-jp.toml")

	lang, _ := getLocale()
	//fmt.Printf("OK: language=%s, locale=%s\n", lang, loc)
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
	//fmt.Println(GetString("init_locale"))
}

func GetString(id string) string {
	return Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: id})
}
