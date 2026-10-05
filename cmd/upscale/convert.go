package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp" // 注册 WebP 解码；编码复用官方 cwebp。
)

type options struct {
	Algorithm string  `json:"algorithm"`
	Model     string  `json:"model"`
	Scale     float64 `json:"scale"`
	Format    string  `json:"format"`
	Quality   int     `json:"quality"`
	Lossless  bool    `json:"lossless"`
	Method    int     `json:"method"`
	Tile      int     `json:"tile"`
	GPU       int     `json:"gpu"`
	TTA       bool    `json:"tta"`
}

// validate 在启动外部进程前校验参数；只接受固定算法及已安装模型。
func (o options) validate(cfg settings) error {
	if math.IsNaN(o.Scale) || math.IsInf(o.Scale, 0) || o.Scale < 0.1 || o.Scale > 4 {
		return errors.New("scale 范围为 0.1–4")
	}
	switch o.Algorithm {
	case "none":
		if o.Scale != 1 {
			return errors.New("仅转换格式时 scale 必须为 1")
		}
	case "lanczos", "cubic", "nearest":
	case "ai":
		found := false
		for _, m := range cfg.AI {
			for _, scale := range m.Scales {
				found = found || (m.Name == o.Model && float64(scale) == o.Scale)
			}
		}
		if !found {
			return errors.New("AI 模型未安装或不支持所选倍率")
		}
	default:
		return errors.New("algorithm 必须为 ai、lanczos、cubic、nearest 或 none")
	}
	if o.Format != "png" && o.Format != "jpeg" && o.Format != "webp" {
		return errors.New("format 必须为 png、jpeg 或 webp")
	}
	if o.Format == "webp" && !cfg.WebP {
		return errors.New("WebP 编码需要安装 cwebp，或通过 -cwebp 指定路径")
	}
	if o.Quality < 0 || o.Quality > 100 || (o.Format == "jpeg" && o.Quality == 0) {
		return errors.New("quality 范围为 0–100（JPEG 最低 1）")
	}
	if o.Method < 0 || o.Method > 6 || o.GPU < -1 || o.GPU > 31 || o.Tile < 0 || o.Tile > 4096 || (o.Tile > 0 && o.Tile < 32) {
		return errors.New("method 范围 0–6；gpu 为 -1（自动）或 0–31；tile 为 0（自动）或 32–4096")
	}
	return nil
}

// dimensions 在解码像素前限制工作量，WebP 的单边尺寸最多为 16383。
func (o options) dimensions(width, height int) (int, int, error) {
	w, h := max(1, int(math.Round(float64(width)*o.Scale))), max(1, int(math.Round(float64(height)*o.Scale)))
	if width <= 0 || height <= 0 || int64(width)*int64(height) > 16_000_000 || int64(w)*int64(h) > 64_000_000 || w > 32768 || h > 32768 {
		return 0, 0, errors.New("图片过大：输入最多 1600 万像素，输出最多 6400 万像素、单边 32768")
	}
	if o.Format == "webp" && (w > 16383 || h > 16383) {
		return 0, 0, errors.New("WebP 输出的宽、高不能超过 16383")
	}
	return w, h, nil
}

// convert 统一 EXIF 朝向并移除元数据；原始上传文件始终保留不变。
func convert(ctx context.Context, cfg settings, dir string, o options, update func(string, float64)) (int, int, error) {
	if err := ctx.Err(); err != nil {
		return 0, 0, err
	}
	update("解码图片", 2)
	src, err := imaging.Open(filepath.Join(dir, "upload"), imaging.AutoOrientation(true))
	if err != nil {
		return 0, 0, err
	}
	w, h, err := o.dimensions(src.Bounds().Dx(), src.Bounds().Dy())
	if err != nil {
		return 0, 0, err
	}
	var output image.Image = src
	if o.Algorithm == "ai" {
		input, target := filepath.Join(dir, "input.png"), filepath.Join(dir, "scaled.png")
		if err := imaging.Save(src, input); err != nil {
			return 0, 0, err
		}
		update("AI 推理（首次加载模型可能较慢）", 5)
		args := []string{"-i", input, "-o", target, "-m", cfg.Models, "-n", o.Model, "-s", strconv.Itoa(int(o.Scale)), "-t", strconv.Itoa(o.Tile), "-j", "1:1:1"}
		if o.GPU >= 0 {
			args = append(args, "-g", strconv.Itoa(o.GPU))
		}
		if o.TTA {
			args = append(args, "-x")
		}
		if err := runEngine(ctx, cfg.Engine, args, func(p float64) { update("AI 推理", 5+p*0.75) }); err != nil {
			return 0, 0, err
		}
		output, err = imaging.Open(target)
		if err != nil {
			return 0, 0, err
		}
		if output.Bounds().Dx() != w || output.Bounds().Dy() != h {
			return 0, 0, errors.New("AI 引擎输出尺寸不符合所选倍率")
		}
	} else if o.Algorithm != "none" {
		update("普通缩放", 10)
		filters := map[string]imaging.ResampleFilter{"lanczos": imaging.Lanczos, "cubic": imaging.CatmullRom, "nearest": imaging.NearestNeighbor}
		output = imaging.Resize(src, w, h, filters[o.Algorithm])
	}
	if err := ctx.Err(); err != nil {
		return 0, 0, err
	}
	update("编码输出", 82)
	result := filepath.Join(dir, "result."+o.Format)
	switch o.Format {
	case "webp":
		intermediate := filepath.Join(dir, "encode.png")
		if err := imaging.Save(output, intermediate); err != nil {
			return 0, 0, err
		}
		args := []string{"-q", strconv.Itoa(o.Quality), "-m", strconv.Itoa(o.Method), "-mt", "-progress", "-o", result}
		if o.Lossless {
			args = append(args, "-lossless", "-exact")
		}
		args = append(args, "--", intermediate)
		err = runEngine(ctx, cfg.CWebP, args, func(p float64) { update("WebP 编码", 82+p*0.17) })
	case "jpeg":
		// JPEG 无透明通道，明确使用白底，避免透明漫画页变成黑底。
		matte := image.NewRGBA(output.Bounds())
		draw.Draw(matte, matte.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
		draw.Draw(matte, matte.Bounds(), output, output.Bounds().Min, draw.Over)
		err = imaging.Save(matte, result, imaging.JPEGQuality(o.Quality))
	default:
		err = imaging.Save(output, result)
	}
	return w, h, err
}

var percentage = regexp.MustCompile(`(\d+(?:\.\d+)?)\s*%`)

// splitProgress 同时支持引擎逐行输出和 cwebp 的回车刷新。
func splitProgress(data []byte, atEOF bool) (int, []byte, error) {
	if i := bytes.IndexAny(data, "\r\n"); i >= 0 {
		return i + 1, data[:i], nil
	}
	if atEOF && len(data) > 0 {
		return len(data), data, nil
	}
	return 0, nil, nil
}

// runEngine 不经过 shell；只解析真实百分比，保留末尾诊断用于报告失败。
func runEngine(ctx context.Context, executable string, args []string, update func(float64)) error {
	cmd := exec.CommandContext(ctx, executable, args...)
	cmd.Stdout = io.Discard
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	scan := bufio.NewScanner(stderr)
	scan.Split(splitProgress)
	var tail string
	for scan.Scan() {
		line := scan.Text()
		if match := percentage.FindStringSubmatch(line); len(match) > 1 {
			p, _ := strconv.ParseFloat(match[1], 64)
			update(min(100, max(0, p)))
		} else {
			tail += line + "\n"
			if len(tail) > 4096 {
				tail = tail[len(tail)-4096:]
			}
		}
	}
	if scan.Err() != nil {
		_ = cmd.Process.Kill()
	}
	if err := errors.Join(cmd.Wait(), scan.Err()); err != nil {
		return fmt.Errorf("外部转换失败：%w\n%s", err, strings.TrimSpace(tail))
	}
	return nil
}
