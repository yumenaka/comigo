//go:build js && wasm

package main

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"path/filepath"
	"strings"
	"syscall/js"
	"time"

	"github.com/yumenaka/comigo/tools/archivebook"
)

type archiveState struct {
	Filename string
	Data     []byte
	Pages    []archivebook.Page
	Book     map[string]any
}

var current archiveState

func main() {
	api := map[string]any{
		"open":     js.FuncOf(openArchive),
		"readPage": js.FuncOf(readPage),
		"close":    js.FuncOf(closeArchive),
	}
	js.Global().Set("ComiGoArchive", js.ValueOf(api))
	select {}
}

func openArchive(this js.Value, args []js.Value) any {
	if len(args) < 2 {
		return reject("ComiGoArchive.open requires arrayBuffer and filename")
	}
	arrayBuffer := args[0]
	filename := args[1].String()
	options := parseOptions(args)

	return promise(func(resolve js.Value, rejectFn js.Value) {
		data := make([]byte, arrayBuffer.Get("byteLength").Int())
		js.CopyBytesToGo(data, js.Global().Get("Uint8Array").New(arrayBuffer))

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		pages, err := archivebook.ListPages(ctx, filename, bytes.NewReader(data), int64(len(data)), options)
		if err != nil {
			rejectFn.Invoke(err.Error())
			return
		}
		bookID := buildBookID(filename, data)
		for i := range pages {
			pages[i].URL = "comigo-reader://" + pages[i].Name
		}
		current = archiveState{
			Filename: filename,
			Data:     data,
			Pages:    pages,
			Book: map[string]any{
				"id":          bookID,
				"title":       filename,
				"type":        strings.ToLower(filepath.Ext(filename)),
				"page_count":  len(pages),
				"PageInfos":   pages,
				"reader_only": true,
			},
		}
		resolve.Invoke(toJS(current.Book))
	})
}

func readPage(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return reject("ComiGoArchive.readPage requires page index or name")
	}
	target := args[0]
	name := ""
	if target.Type() == js.TypeNumber {
		index := target.Int()
		if index < 0 || index >= len(current.Pages) {
			return reject("page index out of range")
		}
		name = current.Pages[index].Name
	} else {
		name = target.String()
	}
	options := archivebook.Options{}
	if len(args) >= 2 && args[1].Type() == js.TypeObject {
		options.TextEncoding = args[1].Get("textEncoding").String()
	}

	return promise(func(resolve js.Value, rejectFn js.Value) {
		if len(current.Data) == 0 {
			rejectFn.Invoke("no archive opened")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		data, err := archivebook.ReadPage(ctx, current.Filename, bytes.NewReader(current.Data), name, options)
		if err != nil {
			rejectFn.Invoke(err.Error())
			return
		}
		uint8Array := js.Global().Get("Uint8Array").New(len(data))
		js.CopyBytesToJS(uint8Array, data)
		resolve.Invoke(uint8Array)
	})
}

func closeArchive(this js.Value, args []js.Value) any {
	current = archiveState{}
	return nil
}

func parseOptions(args []js.Value) archivebook.Options {
	opt := archivebook.Options{SortBy: "default"}
	if len(args) < 3 || args[2].Type() != js.TypeObject {
		return opt
	}
	raw := args[2]
	opt.TextEncoding = raw.Get("textEncoding").String()
	opt.SortBy = raw.Get("sortBy").String()
	return opt
}

func promise(run func(resolve js.Value, reject js.Value)) js.Value {
	handler := js.FuncOf(func(this js.Value, args []js.Value) any {
		resolve := args[0]
		rejectFn := args[1]
		go run(resolve, rejectFn)
		return nil
	})
	return js.Global().Get("Promise").New(handler)
}

func reject(message string) js.Value {
	return promise(func(resolve js.Value, reject js.Value) {
		reject.Invoke(message)
	})
}

func toJS(v any) js.Value {
	data, _ := json.Marshal(v)
	return js.Global().Get("JSON").Call("parse", string(data))
}

func buildBookID(filename string, data []byte) string {
	hash := sha1.New()
	hash.Write([]byte(filename))
	hash.Write(data[:min(len(data), 65536)])
	return "reader-" + hex.EncodeToString(hash.Sum(nil))[:12]
}
