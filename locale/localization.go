package locale

import (
	_ "embed"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/example_test.go

var Localizer *i18n.Localizer

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
			loc := langLoc[1]
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
			loc := langLoc[1]
			return lang, loc
		}
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envlang)
			langLocRaw = strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			lang := langLoc[0]
			loc := langLoc[1]
			return lang, loc
		}
	}
	return defaultLang, defaultLoc
}

func chcpToUTF8() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("CMD", "/C", "chcp.com", "65001")
		if err := cmd.Start(); err != nil {
			fmt.Println("设置Windows活动代码页失败")
			fmt.Println(err.Error())
		}
	}
}

func InitLocale() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	//go:embed en-us.toml
	var enBytes []byte
	bundle.MustParseMessageFileBytes(enBytes, "en-us.toml")

	//go:embed zh-cn.toml
	var cnBytes []byte
	bundle.MustParseMessageFileBytes(cnBytes, "zh-cn.toml")
	lang, loc := getLocale()
	fmt.Printf("OK: language=%s, locale=%s\n", lang, loc)
	switch lang {
	case "zh":
		Localizer = i18n.NewLocalizer(bundle, "zh-CN")
		fmt.Println(Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "init_locale"}))
	case "en":
		Localizer = i18n.NewLocalizer(bundle, "en-US")
		fmt.Println(Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "init_locale"}))
	case "ja":
		Localizer = i18n.NewLocalizer(bundle, "en-US")
		fmt.Println(Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID: "init_locale"}))
	}
}

func GetString(id string) string {
	return Localizer.MustLocalize(&i18n.LocalizeConfig{MessageID:id})
}
