module github.com/yumenaka/comigo/htmx

go 1.24

// 从 Go 1.24 开始，模块可以声明工具依赖项
// https://go.dev/doc/modules/managing-dependencies#tools
// go get -tool github.com/a-h/templ/cmd/templ
//tool (
//  golang.org/x/tools/cmd/stringer
//  ./cmd/migrate
//)

tool (
	github.com/a-h/templ/cmd/templ
	github.com/air-verse/air
)

require (
	github.com/a-h/templ v0.3.833
	github.com/angelofallars/htmx-go v0.5.0
	github.com/yumenaka/comigo v0.9.14-0.20241229132926-7d88d9e595c0
)

// 使用最新的 master 提交 然后，在项目根目录下运行以下命令来下载最新的依赖 go mod tidy && go mod download
// replace github.com/yumenaka/comigo => github.com/yumenaka/comigo master
//replace github.com/yumenaka/comigo => github.com/yumenaka/comigo <commit-hash>
// replace github.com/yumenaka/comigo => github.com/yumenaka/comigo v0.9.12-0.20241002234112-6c9ad7a143e4
replace github.com/yumenaka/comigo v0.9.14-0.20241229132926-7d88d9e595c0 => ../

require (
	github.com/charmbracelet/bubbletea v1.3.3
	github.com/labstack/echo/v4 v4.13.3
	github.com/labstack/gommon v0.4.2
	github.com/spf13/viper v1.19.0
)

require (
	aead.dev/minisign v0.3.0 // indirect
	ariga.io/atlas v0.32.0 // indirect
	dario.cat/mergo v1.0.1 // indirect
	entgo.io/ent v0.14.3 // indirect
	github.com/Baozisoftware/qrcode-terminal-go v0.0.0-20170407111555-c0650d8dff0f // indirect
	github.com/BurntSushi/toml v1.4.0 // indirect
	github.com/PuerkitoBio/goquery v1.10.1 // indirect
	github.com/STARRY-S/zip v0.2.2 // indirect
	github.com/VividCortex/ewma v1.2.0 // indirect
	github.com/a-h/parse v0.0.0-20250122154542-74294addb73e // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/air-verse/air v1.61.7 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bbrks/go-blurhash v1.1.1 // indirect
	github.com/bep/godartsass v1.2.0 // indirect
	github.com/bep/godartsass/v2 v2.1.0 // indirect
	github.com/bep/golibsass v1.2.0 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/bodgit/plumbing v1.3.0 // indirect
	github.com/bodgit/sevenzip v1.6.0 // indirect
	github.com/bodgit/windows v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/charmbracelet/lipgloss v1.0.0 // indirect
	github.com/charmbracelet/x/ansi v0.8.0 // indirect
	github.com/charmbracelet/x/term v0.2.1 // indirect
	github.com/cheggaaa/pb/v3 v3.1.7 // indirect
	github.com/cli/browser v1.3.0 // indirect
	github.com/cli/safeexec v1.0.1 // indirect
	github.com/cloudsoda/go-smb2 v0.0.0-20250228001242-d4c70e6251cc // indirect
	github.com/cloudsoda/sddl v0.0.0-20250224235906-926454e91efc // indirect
	github.com/creack/pty v1.1.23 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/dsnet/compress v0.0.2-0.20230904184137-39efe44ab707 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/erikgeiser/coninput v0.0.0-20211004153227-1c3628e74d0f // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/geoffgarside/ber v1.1.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-openapi/inflect v0.21.0 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gohugoio/hugo v0.134.3 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/hcl/v2 v2.23.0 // indirect
	github.com/hhrutter/lzw v1.0.0 // indirect
	github.com/hhrutter/tiff v1.0.2 // indirect
	github.com/hirochachacha/go-smb2 v1.1.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/goidentity/v6 v6.0.1 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/jxskiss/base62 v1.1.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible // indirect
	github.com/lestrrat-go/strftime v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20250303091104-876f3ea5145d // indirect
	github.com/magiconair/properties v1.8.9 // indirect
	github.com/mandykoh/autocrop v0.4.7 // indirect
	github.com/mandykoh/go-parallel v0.1.0 // indirect
	github.com/mandykoh/prism v0.35.3 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/minio/selfupdate v0.6.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.1-0.20231216201459-8508981c8b6c // indirect
	github.com/muesli/ansi v0.0.0-20230316100256-276c6243b2f6 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/termenv v0.15.2 // indirect
	github.com/natefinch/atomic v1.0.1 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/nicksnyder/go-i18n/v2 v2.5.1 // indirect
	github.com/nwaples/rardecode/v2 v2.1.1 // indirect
	github.com/pdfcpu/pdfcpu v0.9.1 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sanity-io/litter v1.5.8 // indirect
	github.com/shirou/gopsutil/v3 v3.24.5 // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e // indirect
	github.com/sorairolake/lzip-go v0.3.5 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/spf13/cast v1.7.1 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tdewolff/parse/v2 v2.7.15 // indirect
	github.com/therootcompany/xz v1.0.1 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/ulikunitz/xz v0.5.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xxjwxc/gowp v0.0.0-20240929033016-5be68d222389 // indirect
	github.com/xxjwxc/public v0.0.0-20250124084855-a0e71c7e0764 // indirect
	github.com/yumenaka/archives v0.0.0-20241222141717-1f2992f21052 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zclconf/go-cty v1.16.2 // indirect
	github.com/zclconf/go-cty-yaml v1.1.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go4.org v0.0.0-20230225012048-214862532bf5 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/image v0.25.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	golang.org/x/tools v0.31.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/eapache/queue.v1 v1.1.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.61.13 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.8.2 // indirect
	modernc.org/sqlite v1.36.0 // indirect
)
