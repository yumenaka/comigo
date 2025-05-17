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

var (
	Localizer *i18n.Localizer
)

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
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp title")
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
	case "linux", "darwin":
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

func init() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustParseMessageFileBytes(enBytes, "en_US.json")
	bundle.MustParseMessageFileBytes(cnBytes, "zh_CN.json")
	bundle.MustParseMessageFileBytes(jpBytes, "ja_JP.json")

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

func GetString(id string) string {
	return Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: id})
}
