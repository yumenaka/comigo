//go:build js

package main

import (
	"archive/zip"
	"bytes"
	"io"
	"strings"
	"syscall/js"
)

// cp $(go env GOROOT)/lib/wasm/wasm_exec.js .
// env GOOS=js GOARCH=wasm go build -o main.wasm cmd/wasm/main.go

// -- 关键解释 & 常见坑 --
// 1. JS 与 Go 互拷字节	js.CopyBytesToGo 与 js.CopyBytesToJS 零拷贝；速度比逐字节 Get("getUint8") 高一个量级。
// 2. 内存限制	浏览器 WebAssembly 默认堆上限 ~4 GB，但单次 make([]byte) 仍要注意；大档建议分块或用 JS stream。
// 3. 解码格式	我们没有在 Go 里解码图片，只是原样回传给浏览器。这样能保持 zip 内原始文件的体积与格式，并由 <img> 自行解码。
// 4. TinyGo	若追求 100 KB 级别的产物，可用 tinygo build -target wasm；但需要 image/* 时请确认 TinyGo 对 archive/zip/image 的支持程度。
// 5. 安全	仅在浏览器本地运行，无后端上传，满足普通隐私诉求；如果放到服务器仍走 HTTPS，浏览器同源策略即可保护 Blob URL。

// 仅接受常见静态图格式
func isImage(name string) bool {
	n := strings.ToLower(name)
	return strings.HasSuffix(n, ".png") ||
		strings.HasSuffix(n, ".jpg") || strings.HasSuffix(n, ".jpeg") ||
		strings.HasSuffix(n, ".gif") || strings.HasSuffix(n, ".webp") ||
		strings.HasSuffix(n, ".avif")
}

// 由 JS 传入 Uint8Array，返回 JS 数组，每项含 {name, data(Uint8Array)}
func unzipImages(_ js.Value, args []js.Value) any {
	if len(args) == 0 {
		js.Global().Call("console.error", "missing zip bytes")
		return nil
	}

	zipBytesJS := args[0]
	size := zipBytesJS.Get("length").Int()
	data := make([]byte, size)
	js.CopyBytesToGo(data, zipBytesJS)

	zr, err := zip.NewReader(bytes.NewReader(data), int64(size))
	if err != nil {
		js.Global().Call("console.error", err.Error())
		return nil
	}

	array := js.Global().Get("Array").New()
	for _, f := range zr.File {
		if !isImage(f.Name) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			continue
		}
		buf, _ := io.ReadAll(rc)
		rc.Close()

		u8 := js.Global().Get("Uint8Array").New(len(buf))
		js.CopyBytesToJS(u8, buf)

		obj := js.Global().Get("Object").New()
		obj.Set("name", f.Name)
		obj.Set("data", u8)
		array.Call("push", obj)
	}
	return array
}

func main() {
	js.Global().Set("unzipImages", js.FuncOf(unzipImages))
	select {} // 保持运行
}
