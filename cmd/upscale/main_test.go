package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// testServer 启动真实 HTTP 服务和可选的转换 worker，退出时等待计算结束。
func testServer(t *testing.T, cfg settings, worker bool) (*server, string) {
	t.Helper()
	s := newServer(t.TempDir(), cfg)
	httpServer := httptest.NewServer(s.routes())
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	if worker {
		go func() { defer close(done); s.work(ctx) }()
	} else {
		close(done)
	}
	t.Cleanup(func() {
		cancel()
		s.mu.Lock()
		for _, j := range s.jobs {
			j.cancel()
		}
		s.mu.Unlock()
		<-done
		httpServer.Close()
	})
	return s, httpServer.URL
}

// uploadTest 提交合成透明图片，验证转换不依赖仓库外的测试素材。
func uploadTest(t *testing.T, base, opts string) (int, status) {
	t.Helper()
	img := image.NewNRGBA(image.Rect(0, 0, 64, 48))
	for y := 0; y < 48; y++ {
		for x := 8; x < 64; x++ {
			img.SetNRGBA(x, y, color.NRGBA{R: uint8(x * 4), G: uint8(y * 5), B: 120, A: 255})
		}
	}
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	file, err := w.CreateFormFile("file", "test.png")
	if err != nil {
		t.Fatal(err)
	}
	if err := png.Encode(file, img); err != nil {
		t.Fatal(err)
	}
	if err := w.WriteField("options", opts); err != nil {
		t.Fatal(err)
	}
	w.Close()
	response, err := http.Post(base+"/api/jobs", w.FormDataContentType(), &body)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var state status
	if err := json.NewDecoder(response.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	return response.StatusCode, state
}

// awaitTest 从 SSE 读取直到终态，同时检查进度不会倒退或提前完成。
func awaitTest(t *testing.T, base, id string) status {
	t.Helper()
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Get(base + "/api/jobs/" + id + "/events")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != 200 || !strings.Contains(response.Header.Get("Content-Type"), "text/event-stream") {
		t.Fatalf("SSE: %s", response.Status)
	}
	var state status
	last := 0.0
	scan := bufio.NewScanner(response.Body)
	for scan.Scan() {
		if !strings.HasPrefix(scan.Text(), "data: ") {
			continue
		}
		if err := json.Unmarshal([]byte(strings.TrimPrefix(scan.Text(), "data: ")), &state); err != nil {
			t.Fatal(err)
		}
		if state.Progress < last || (state.Progress == 100 && state.State != "done") {
			t.Fatalf("无效进度：%+v", state)
		}
		last = state.Progress
	}
	if err := scan.Err(); err != nil {
		t.Fatal(err)
	}
	if !terminal(state.State) {
		t.Fatalf("SSE 未报告终态：%+v", state)
	}
	return state
}

// TestConversionAPI 覆盖 HTTP 上传、SSE、尺寸、透明度、编码及结果清理。
func TestConversionAPI(t *testing.T) {
	cfg := discover("", "", "cwebp")
	_, base := testServer(t, cfg, true)
	for _, tc := range []struct {
		algorithm, format string
		scale             float64
		width, height     int
	}{
		{"lanczos", "png", 0.5, 32, 24},
		{"cubic", "png", 2, 128, 96},
		{"nearest", "jpeg", 2, 128, 96},
		{"none", "webp", 1, 64, 48},
	} {
		t.Run(tc.algorithm+"-"+tc.format, func(t *testing.T) {
			if tc.format == "webp" && !cfg.WebP {
				t.Skip("未安装 cwebp；安装后会验证真实无损 WebP 编码")
			}
			data, _ := json.Marshal(map[string]any{"algorithm": tc.algorithm, "scale": tc.scale, "format": tc.format, "lossless": true, "quality": 100})
			code, state := uploadTest(t, base, string(data))
			if code != 202 {
				t.Fatalf("提交状态：%d", code)
			}
			state = awaitTest(t, base, state.ID)
			if state.State != "done" || state.Width != tc.width || state.Height != tc.height || state.Bytes == 0 {
				t.Fatalf("转换失败：%+v", state)
			}
			response, err := http.Get(base + state.Result)
			if err != nil {
				t.Fatal(err)
			}
			decoded, format, err := image.Decode(response.Body)
			response.Body.Close()
			if err != nil || format != tc.format || decoded.Bounds().Dx() != tc.width || decoded.Bounds().Dy() != tc.height {
				t.Fatalf("结果不符：format=%s, err=%v", format, err)
			}
			r, g, b, a := decoded.At(0, 0).RGBA()
			if tc.format == "jpeg" {
				if r < 60000 || g < 60000 || b < 60000 {
					t.Fatal("JPEG 透明处应为白色")
				}
			} else if a != 0 {
				t.Fatal("透明通道丢失")
			}
			if tc.format == "webp" {
				pixel := color.NRGBAModel.Convert(decoded.At(20, 20)).(color.NRGBA)
				if pixel != (color.NRGBA{R: 80, G: 100, B: 120, A: 255}) {
					t.Fatalf("无损 WebP 像素改变：%v", pixel)
				}
			}
			req, _ := http.NewRequest(http.MethodDelete, base+"/api/jobs/"+state.ID, nil)
			response, err = http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			response.Body.Close()
			if response.StatusCode != 204 {
				t.Fatalf("删除任务失败：%s", response.Status)
			}
		})
	}
}

// TestValidationAndCancel 确认错误参数不启动任务，排队任务可立即取消并删除。
func TestValidationAndCancel(t *testing.T) {
	s, base := testServer(t, settings{}, false)
	for _, opts := range []string{`{"scale":0}`, `{"tile":1}`, `{"algorithm":"ai","model":"../model"}`, `{"format":"webp"}`, `{"algorithm":"none","scale":2}`, `{"quality":101}`, `{"gpu":-2}`} {
		if code, _ := uploadTest(t, base, opts); code != 400 {
			t.Fatalf("应拒绝 %s，实际状态 %d", opts, code)
		}
	}
	if _, _, err := (options{Scale: 4, Format: "webp"}).dimensions(5000, 2); err == nil {
		t.Fatal("未限制 WebP 单边尺寸")
	}
	if _, _, err := (options{Scale: 4}).dimensions(4000, 3000); err == nil {
		t.Fatal("未限制输出像素总量")
	}
	_, state := uploadTest(t, base, `{}`)
	response, err := http.Post(base+"/api/jobs/"+state.ID+"/cancel", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	state = awaitTest(t, base, state.ID)
	if state.State != "canceled" {
		t.Fatalf("取消失败：%+v", state)
	}
	s.mu.Lock()
	dir := s.jobs[state.ID].dir
	s.mu.Unlock()
	req, _ := http.NewRequest(http.MethodDelete, base+"/api/jobs/"+state.ID, nil)
	response, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Fatal("任务文件未删除")
	}
}

// TestFailedJob 确认外部工具故障会成为失败状态，不会提供半成品或堵塞后续任务。
func TestFailedJob(t *testing.T) {
	_, base := testServer(t, settings{WebP: true, CWebP: filepath.Join(t.TempDir(), "missing-cwebp")}, true)
	_, state := uploadTest(t, base, `{"algorithm":"none","scale":1,"format":"webp"}`)
	state = awaitTest(t, base, state.ID)
	if state.State != "failed" || state.Error == "" || state.Result != "" {
		t.Fatalf("错误状态不符：%+v", state)
	}
	response, err := http.Get(base + "/api/jobs/" + state.ID + "/result")
	if err != nil {
		t.Fatal(err)
	}
	response.Body.Close()
	if response.StatusCode != 409 {
		t.Fatal("失败任务不应提供结果")
	}
	_, state = uploadTest(t, base, `{}`)
	if state = awaitTest(t, base, state.ID); state.State != "done" {
		t.Fatalf("失败阻塞了后续任务：%+v", state)
	}
}

// TestRealAI 通过显式环境变量启用真实 GPU 实验，不把模型或下载动作放进常规测试。
func TestRealAI(t *testing.T) {
	engine := os.Getenv("UPSCALE_TEST_ENGINE")
	if engine == "" {
		t.Skip("设置 UPSCALE_TEST_ENGINE 可运行真实 ncnn 测试")
	}
	cfg := discover(engine, os.Getenv("UPSCALE_TEST_MODELS"), "cwebp")
	if len(cfg.AI) == 0 {
		t.Fatal("未找到引擎和模型")
	}
	_, base := testServer(t, cfg, true)
	for _, m := range cfg.AI {
		t.Run(m.Name, func(t *testing.T) {
			for _, scale := range m.Scales {
				data, _ := json.Marshal(map[string]any{"algorithm": "ai", "model": m.Name, "scale": scale, "tile": 32, "format": "png"})
				code, state := uploadTest(t, base, string(data))
				if code != 202 {
					t.Fatalf("提交失败：%d", code)
				}
				state = awaitTest(t, base, state.ID)
				if state.State != "done" || state.Width != 64*scale || state.Height != 48*scale {
					t.Fatalf("AI 输出不符：%+v", state)
				}
				response, err := http.Get(base + state.Result)
				if err != nil {
					t.Fatal(err)
				}
				result, _, err := image.Decode(response.Body)
				response.Body.Close()
				if err != nil {
					t.Fatal(err)
				}
				_, _, _, alpha := result.At(0, 0).RGBA()
				if alpha != 0 {
					t.Fatal("AI 丢失透明通道")
				}
				t.Logf("%dx %dx%d %.3fs", scale, state.Width, state.Height, state.Elapsed)
			}
		})
	}
}

// TestProgressAndFailure 验证回车进度解析和真实子进程错误传播。
func TestProgressAndFailure(t *testing.T) {
	scan := bufio.NewScanner(strings.NewReader("engine\r12.50%\n[encode] 99 %\r"))
	scan.Split(splitProgress)
	var values []string
	for scan.Scan() {
		if match := percentage.FindStringSubmatch(scan.Text()); len(match) > 1 {
			values = append(values, match[1])
		}
	}
	if strings.Join(values, ",") != "12.50,99" {
		t.Fatalf("进度解析错误：%v", values)
	}
	if err := runEngine(context.Background(), filepath.Join(t.TempDir(), "missing-engine"), nil, func(float64) {}); err == nil {
		t.Fatal("未报告缺失引擎")
	}
	// 使用 Go 测试进程自身，跨平台验证非零退出及 stderr 诊断。
	if err := runEngine(context.Background(), os.Args[0], []string{"-invalid-upscale-test-flag"}, func(float64) {}); err == nil || !strings.Contains(err.Error(), "flag") {
		t.Fatalf("未报告子进程错误：%v", err)
	}
}
