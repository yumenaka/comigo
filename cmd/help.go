package cmd

import "github.com/yumenaka/comigo/assets/locale"

// refreshCLIText 在参数解析后更新帮助文案，不读取用户配置或初始化服务。
func refreshCLIText() {
	RootCmd.Use = locale.GetString("comigo_use")
	RootCmd.Short = locale.GetString("short_description")
	RootCmd.Long = locale.GetString("long_description") + "\n\n" + locale.GetString("cli_flags_scope")
	RootCmd.Example = locale.GetString("comigo_example")
	for name, key := range map[string]string{
		"config":                  "config",
		"username":                "username",
		"password":                "password",
		"timeout":                 "timeout",
		"auto-rescan-min":         "auto_rescan_interval_minutes",
		"database":                "enable_database",
		"db-type":                 "db_type",
		"db-dsn":                  "db_dsn",
		"port":                    "port",
		"host":                    "local_host",
		"base-path":               "base_path_description",
		"tls":                     "tls_enable",
		"auto-tls":                "auto_https_cert",
		"tls-crt":                 "tls_crt",
		"tls-key":                 "tls_key",
		"enable-upload":           "enable_file_upload",
		"open-browser":            "open_browser",
		"no-default-library":      "no_default_library",
		"no-tui":                  "no_tui",
		"temp":                    "temp_reader_mode",
		"local":                   "disable_lan",
		"max-depth":               "max_depth",
		"min-image":               "min_media_num",
		"log-file":                "log_to_file",
		"use-cache":               "cache_file_enable",
		"cache-dir":               "cache_file_dir",
		"cache-clean":             "cache_file_clean",
		"zip-encode":              "zip_encode",
		"tailscale":               "enable_tailscale",
		"tailscale-funnel":        "funnel_tunnel_label",
		"funnel-password-check":   "funnel_login_check",
		"tailscale-hostname":      "tailscale_hostname",
		"tailscale-port":          "tailscale_port",
		"tailscale-authKey":       "tailscale_auth_key",
		"read-only":               "read_only_mode",
		"lang":                    "lang",
		"register-context-menu":   "register_context_menu",
		"unregister-context-menu": "unregister_context_menu",
		"plugin":                  "plugin_enable",
		"debug":                   "debug_mode",
		"upgrade":                 "self_upgrade_flag",
	} {
		if flag := RootCmd.PersistentFlags().Lookup(name); flag != nil {
			flag.Usage = locale.GetString(key)
		}
	}
	for _, command := range RootCmd.Commands() {
		switch command.Name() {
		case "run", "start", "status", "stop", "reload", "upgrade", "version":
			command.Short = locale.GetString("cli_" + command.Name())
			if command.Name() == "status" {
				command.Long = command.Short + "\n\n" + locale.GetString("cli_status_details")
			}
			if command.Name() == "stop" || command.Name() == "reload" {
				command.Long = command.Short + "\n\n" + locale.GetString("cli_flags_scope")
			}
		}
	}
}
