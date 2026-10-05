// upscale 是供 Comigo 调用的独立图片转换实验程序。
package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

//go:embed index.html
var page string

const maxJobs = 16

type model struct {
	Name   string `json:"name"`
	Scales []int  `json:"scales"`
}

type settings struct {
	Engine string  `json:"-"`
	Models string  `json:"-"`
	CWebP  string  `json:"-"`
	AI     []model `json:"models"`
	WebP   bool    `json:"webp"`
}

type status struct {
	ID       string  `json:"id"`
	State    string  `json:"state"`
	Stage    string  `json:"stage"`
	Progress float64 `json:"progress"`
	Elapsed  float64 `json:"elapsed_seconds"`
	Error    string  `json:"error,omitempty"`
	Width    int     `json:"width,omitempty"`
	Height   int     `json:"height,omitempty"`
	Bytes    int64   `json:"bytes,omitempty"`
	Format   string  `json:"format"`
	Result   string  `json:"result_url,omitempty"`
}

type job struct {
	status
	opts     options
	dir      string
	created  time.Time
	finished time.Time
	ctx      context.Context
	cancel   context.CancelFunc
}

type server struct {
	settings settings
	root     string
	mu       sync.Mutex // 保护任务表以及任务的状态；转换参数创建后不再修改。
	jobs     map[string]*job
	queue    chan *job
}

// main 只启动辅助服务，不引入 Comigo 的配置、书库或启动链路。
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run 等待 worker 和 HTTP 服务退出后清理临时文件。
func run() error {
	listen := flag.String("listen", "127.0.0.1:1235", "HTTP 监听地址（无认证）")
	engine := flag.String("engine", "realesrgan-ncnn-vulkan", "Real-ESRGAN ncnn 可执行文件")
	models := flag.String("models", "", "模型目录，默认使用引擎旁的 models")
	cwebp := flag.String("cwebp", "cwebp", "WebP 编码器可执行文件")
	flag.Parse()
	cfg := discover(*engine, *models, *cwebp)
	root, err := os.MkdirTemp("", "comigo-upscale-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(root)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	s := newServer(root, cfg)
	done := make(chan struct{})
	go func() { defer close(done); s.work(ctx) }()
	e := s.routes()
	shutdownDone := make(chan struct{})
	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		s.mu.Lock()
		for _, j := range s.jobs {
			j.cancel()
		}
		s.mu.Unlock()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = e.Shutdown(shutdown)
	}()
	log.Printf("图片辅助程序：http://%s（AI 模型 %d 个，WebP：%t）", *listen, len(cfg.AI), cfg.WebP)
	err = e.Start(*listen)
	stop()
	<-done
	<-shutdownDone
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

// discover 只公布实际安装的模型；缺少外部引擎时普通缩放仍可使用。
func discover(engine, models, cwebp string) settings {
	resolve := func(name string) string {
		p, err := exec.LookPath(name)
		if err != nil {
			return ""
		}
		p, _ = filepath.Abs(p)
		return p
	}
	cfg := settings{Engine: resolve(engine), CWebP: resolve(cwebp), AI: []model{}}
	cfg.WebP = cfg.CWebP != ""
	if models == "" {
		models = filepath.Join(filepath.Dir(cfg.Engine), "models")
	}
	cfg.Models, _ = filepath.Abs(models)
	if cfg.Engine == "" {
		return cfg
	}
	for _, m := range []model{{"realesr-animevideov3", []int{2, 3, 4}}, {"realesrgan-x4plus-anime", []int{4}}, {"realesrgan-x4plus", []int{4}}} {
		available := model{Name: m.Name}
		for _, scale := range m.Scales {
			name := m.Name
			if name == "realesr-animevideov3" {
				name += fmt.Sprintf("-x%d", scale)
			}
			bin, binErr := os.Stat(filepath.Join(cfg.Models, name+".bin"))
			param, paramErr := os.Stat(filepath.Join(cfg.Models, name+".param"))
			if binErr == nil && paramErr == nil && bin.Mode().IsRegular() && param.Mode().IsRegular() {
				available.Scales = append(available.Scales, scale)
			}
		}
		if len(available.Scales) > 0 {
			cfg.AI = append(cfg.AI, available)
		}
	}
	return cfg
}

// newServer 使用有界队列，避免图片任务无上限占用内存和磁盘。
func newServer(root string, cfg settings) *server {
	return &server{settings: cfg, root: root, jobs: map[string]*job{}, queue: make(chan *job, maxJobs)}
}

// routes 复用现有 Echo；页面通过原生 fetch 和 SSE 通信。
func (s *server) routes() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error { return c.HTML(http.StatusOK, page) })
	e.GET("/api/capabilities", func(c echo.Context) error { return c.JSON(http.StatusOK, s.settings) })
	e.POST("/api/jobs", s.submit)
	e.GET("/api/jobs/:id", s.get)
	e.GET("/api/jobs/:id/events", s.events)
	e.GET("/api/jobs/:id/result", s.result)
	e.POST("/api/jobs/:id/cancel", s.cancel)
	e.DELETE("/api/jobs/:id", s.remove)
	return e
}

// submit 接收图片和 JSON 参数；请求中不能指定本地路径或可执行文件。
func (s *server) submit(c echo.Context) error {
	r := c.Request()
	r.Body = http.MaxBytesReader(c.Response(), r.Body, 32<<20)
	if err := r.ParseMultipartForm(1 << 20); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "无法读取上传数据（请求上限 32 MiB）")
	}
	defer r.MultipartForm.RemoveAll()
	opts := options{Algorithm: "lanczos", Scale: 2, Format: "png", Quality: 85, Method: 4, GPU: -1}
	if raw := r.FormValue("options"); raw != "" {
		if err := json.Unmarshal([]byte(raw), &opts); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "options 必须是合法 JSON")
		}
	}
	if err := opts.validate(s.settings); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	upload, _, err := r.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "缺少 file 图片")
	}
	defer upload.Close()
	dir, err := os.MkdirTemp(s.root, "job-")
	if err != nil {
		return err
	}
	accepted := false
	defer func() {
		if !accepted {
			_ = os.RemoveAll(dir)
		}
	}()
	input := filepath.Join(dir, "upload")
	f, err := os.Create(input)
	if err != nil {
		return err
	}
	_, copyErr := io.Copy(f, upload)
	closeErr := f.Close()
	if err := errors.Join(copyErr, closeErr); err != nil {
		return err
	}
	f, err = os.Open(input)
	if err != nil {
		return err
	}
	size, format, err := image.DecodeConfig(f)
	f.Close()
	if err != nil || (format != "png" && format != "jpeg" && format != "webp") {
		return echo.NewHTTPError(http.StatusBadRequest, "仅支持可解码的静态 PNG、JPEG、WebP")
	}
	if _, _, err := opts.dimensions(size.Width, size.Height); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.jobs) >= maxJobs {
		return echo.NewHTTPError(http.StatusTooManyRequests, "最多保留 16 个任务，请先 DELETE 已结束的任务")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	j := &job{status: status{ID: filepath.Base(dir), State: "queued", Stage: "等待转换", Format: opts.Format}, opts: opts, dir: dir, created: time.Now(), ctx: ctx, cancel: cancel}
	select {
	case s.queue <- j:
		s.jobs[j.ID] = j
		accepted = true
		c.Response().Header().Set("Location", "/api/jobs/"+j.ID)
		return c.JSON(http.StatusAccepted, j.status)
	default:
		cancel()
		return echo.NewHTTPError(http.StatusTooManyRequests, "转换队列已满")
	}
}

// snapshot 在锁内复制状态，SSE 和查询请求不会读取正在修改的字段。
func (s *server) snapshot(id string) (status, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[id]
	if !ok {
		return status{}, false
	}
	state := j.status
	end := j.finished
	if end.IsZero() {
		end = time.Now()
	}
	state.Elapsed = end.Sub(j.created).Seconds()
	return state, true
}

// get 返回可轮询的任务状态，elapsed_seconds 包含排队时间。
func (s *server) get(c echo.Context) error {
	state, ok := s.snapshot(c.Param("id"))
	if !ok {
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, state)
}

// events 每 250ms 推送快照；重连时直接发送当前状态，无需保存事件历史。
func (s *server) events(c echo.Context) error {
	state, ok := s.snapshot(c.Param("id"))
	if !ok {
		return echo.ErrNotFound
	}
	r := c.Response()
	r.Header().Set("Content-Type", "text/event-stream")
	r.Header().Set("Cache-Control", "no-cache")
	r.Header().Set("X-Accel-Buffering", "no")
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()
	for {
		data, _ := json.Marshal(state)
		if _, err := fmt.Fprintf(r, "data: %s\n\n", data); err != nil {
			return err
		}
		r.Flush()
		if terminal(state.State) {
			return nil
		}
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-tick.C:
			state, ok = s.snapshot(c.Param("id"))
			if !ok {
				return nil
			}
		}
	}
}

// result 只提供完整结果，不暴露正在写入的文件。
func (s *server) result(c echo.Context) error {
	s.mu.Lock()
	j, ok := s.jobs[c.Param("id")]
	if !ok {
		s.mu.Unlock()
		return echo.ErrNotFound
	}
	if j.State != "done" {
		s.mu.Unlock()
		return echo.NewHTTPError(http.StatusConflict, "图片尚未转换成功")
	}
	// 先打开文件再释放锁，避免同时删除任务导致读取到不存在的路径。
	f, err := os.Open(filepath.Join(j.dir, "result."+j.opts.Format))
	name := "result." + j.opts.Format
	s.mu.Unlock()
	if err != nil {
		return err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return err
	}
	c.Response().Header().Set("Cache-Control", "no-store")
	http.ServeContent(c.Response(), c.Request(), name, info.ModTime(), f)
	return nil
}

// cancel 终止排队任务或外部引擎；Go 内部缩放在当前计算步骤结束后响应取消。
func (s *server) cancel(c echo.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[c.Param("id")]
	if !ok {
		return echo.ErrNotFound
	}
	j.cancel()
	if j.State == "queued" {
		j.State, j.Stage, j.Error = "canceled", "已取消", context.Canceled.Error()
		j.finished = time.Now()
	}
	return c.NoContent(http.StatusAccepted)
}

// remove 显式释放已结束任务；运行中的任务须先取消，防止转换与文件删除冲突。
func (s *server) remove(c echo.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[c.Param("id")]
	if !ok {
		return echo.ErrNotFound
	}
	if !terminal(j.State) {
		return echo.NewHTTPError(http.StatusConflict, "请先取消任务并等待结束")
	}
	if err := os.RemoveAll(j.dir); err != nil {
		return err
	}
	delete(s.jobs, j.ID)
	return c.NoContent(http.StatusNoContent)
}

// terminal 判断任务是否已结束。
func terminal(state string) bool {
	return state == "done" || state == "failed" || state == "canceled"
}

// work 串行执行转换，防止多个 AI 任务同时争抢显存。
func (s *server) work(ctx context.Context) {
	// ponytail: 单 worker 限制吞吐；确认需要多 GPU 后再增加按设备分组的队列。
	for {
		select {
		case <-ctx.Done():
			return
		case j := <-s.queue:
			s.mu.Lock()
			if terminal(j.State) {
				s.mu.Unlock()
				continue
			}
			j.State = "running"
			s.mu.Unlock()
			update := func(stage string, progress float64) {
				s.mu.Lock()
				j.Stage = stage
				j.Progress = max(j.Progress, progress)
				s.mu.Unlock()
			}
			width, height, err := convert(j.ctx, s.settings, j.dir, j.opts, update)
			s.mu.Lock()
			j.finished = time.Now()
			switch {
			case j.ctx.Err() != nil:
				j.State, j.Stage, j.Error = "canceled", "已取消或超时", j.ctx.Err().Error()
			case err != nil:
				j.State, j.Stage, j.Error = "failed", "转换失败", err.Error()
			default:
				info, statErr := os.Stat(filepath.Join(j.dir, "result."+j.opts.Format))
				if statErr != nil {
					j.State, j.Stage, j.Error = "failed", "读取结果失败", statErr.Error()
				} else {
					j.State, j.Stage, j.Progress = "done", "已完成", 100
					j.Width, j.Height, j.Bytes = width, height, info.Size()
					j.Result = "/api/jobs/" + j.ID + "/result"
				}
			}
			j.cancel()
			s.mu.Unlock()
		}
	}
}
