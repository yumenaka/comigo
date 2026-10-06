package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
)

// 验证文档中的环境变量覆盖 TOML，并使用真实的配置解码流程校验类型。
func TestEnvironmentVariables(t *testing.T) {
	for _, tc := range []struct {
		name, value, field string
		want               any
	}{
		{"PORT", "2345", "Port", 2345},
		{"ENABLE_UPLOAD", "false", "EnableUpload", false},
		{"MAX_DEPTH", "8", "MaxScanDepth", 8},
		{"MIN_IMAGE", "3", "MinImageNum", 3},
		{"BASE_PATH", "/comics", "BasePath", "/comics"},
		{"LOCAL", "true", "DisableLAN", true},
		{"LANGUAGE", "zh", "Language", "zh"},
		{"DEBUG", "true", "Debug", true},
		{"USERNAME", "fixture-user", "Username", "fixture-user"},
		{"PASSWORD", "fixture-password", "Password", "fixture-password"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := processTestCommand(t)
			t.Setenv("COMIGO_"+tc.name, tc.value)
			file := filepath.Join(t.TempDir(), "config.toml")
			if err := os.WriteFile(file, []byte("Port=1234\nEnableUpload=true\nMaxScanDepth=2\nMinImageNum=1\n"), 0600); err != nil {
				t.Fatal(err)
			}
			if err := root.PersistentFlags().Set("config", file); err != nil {
				t.Fatal(err)
			}
			if err := LoadConfigFile(); err != nil {
				t.Fatal(err)
			}
			got := reflect.ValueOf(config.CopyCfg()).FieldByName(tc.field).Interface()
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("COMIGO_%s 未正确覆盖配置", tc.name)
			}
		})
	}
}

// 环境变量撤销、TOML 删除配置项后必须回到原始默认值，不能固化此前的环境值。
func TestEnvironmentReloadPriority(t *testing.T) {
	root := processTestCommand(t)
	file := filepath.Join(t.TempDir(), "config.toml")
	if err := root.PersistentFlags().Set("config", file); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		env, body, flag string
		want            int
	}{
		{"8", "MaxScanDepth=7", "", 8},
		{"", "MaxScanDepth=7", "", 7},
		{"", "", "", 5},
		{"8", "MaxScanDepth=7", "11", 11},
	} {
		t.Setenv("COMIGO_MAX_DEPTH", tc.env)
		if err := os.WriteFile(file, []byte(tc.body), 0600); err != nil {
			t.Fatal(err)
		}
		if tc.flag != "" {
			if err := root.PersistentFlags().Set("max-depth", tc.flag); err != nil {
				t.Fatal(err)
			}
		}
		if err := LoadConfigFile(); err != nil {
			t.Fatal(err)
		}
		candidate, err := readReloadConfig(file)
		if err != nil || candidate.MaxScanDepth != tc.want {
			t.Fatalf("优先级或默认值回退错误：深度=%d，期望=%d，错误=%v", candidate.MaxScanDepth, tc.want, err)
		}
	}
}

// 只允许显式绑定的参数变量，配置路径和运行时行为仍由命令行参数控制。
func TestRemovedEnvironmentVariables(t *testing.T) {
	root := processTestCommand(t)
	for key, value := range map[string]string{
		"COMIGO_ENABLEUPLOAD": "false", "COMIGO_CONFIGFILE": "missing.toml",
		"COMIGO_TIMEOUTLIMITFORSCAN": "9", "COMIGO_NO_TUI": "true",
		"COMIGO_TEMP": "true", "COMIGO_UPGRADE": "true",
	} {
		t.Setenv(key, value)
	}
	if runtimeViper.GetString("ConfigFile") != "" {
		t.Fatal("无效配置路径变量仍被读取")
	}
	file := filepath.Join(t.TempDir(), "config.toml")
	if err := os.WriteFile(file, []byte("EnableUpload=true\nTimeoutLimitForScan=7"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := root.PersistentFlags().Set("config", file); err != nil {
		t.Fatal(err)
	}
	if err := LoadConfigFile(); err != nil {
		t.Fatal(err)
	}
	cfg := config.GetCfg()
	if !cfg.EnableUpload || cfg.TimeoutLimitForScan != 7 || cfg.NoTUI || cfg.TemporaryReaderMode || cfg.SelfUpgrade {
		t.Fatal("未绑定的变量影响了配置或运行方式")
	}
}

// 自定义配置目录的发现、读取、保存和进程状态定位必须一致。
func TestConfigDirectoryEnvironment(t *testing.T) {
	processTestCommand(t)
	dir := t.TempDir()
	t.Setenv("COMIGO_CONFIG_DIR", dir)
	file := filepath.Join(dir, config.PlatformConfigFilename())
	if err := os.WriteFile(file, []byte("Port=2345"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := LoadConfigFile(); err != nil {
		t.Fatal(err)
	}
	if config.GetCfg().ConfigFile != file || config.GetCfg().Port != 2345 {
		t.Fatal("未读取环境变量指定目录中的配置")
	}
	base, err := processStatePath()
	if err != nil || filepath.Dir(base) != dir {
		t.Fatal("进程状态目录不一致", err)
	}
	if err := config.SaveConfig(config.HomeDirectory); err != nil {
		t.Fatal(err)
	}
	if config.GetCfg().ConfigFile != file {
		t.Fatal("保存位置不一致")
	}
}

// 帮助无需加载 TOML，但仍遵守显式语言参数优先于环境变量的规则。
func TestEnvironmentHelpLanguage(t *testing.T) {
	t.Cleanup(func() { locale.InitLanguageFromConfig("auto") })
	for _, lang := range []string{"zh", "en"} {
		processTestCommand(t)
		t.Setenv("COMIGO_LANGUAGE", "zh")
		args := []string{"status", "--help"}
		if lang == "en" {
			args = append(args, "--lang", "en")
		}
		var out bytes.Buffer
		handled, err := RunProcessCommand(args, &out)
		if err != nil || !handled || config.GetCfg().Language != lang || !strings.Contains(out.String(), locale.GetString("cli_status")) {
			t.Fatalf("帮助语言未生效：%s，错误=%v", lang, err)
		}
	}
}
