package archivebook

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"net/url"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/encoding"
)

// Page 描述压缩包内可阅读的单个媒体文件。
// 该结构保持轻量，方便服务端模板与 wasm 前端共同使用。
type Page struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
	URL     string    `json:"url"`
	PageNum int       `json:"page_num"`
}

// Options 控制压缩包扫描和读取行为。
type Options struct {
	BookID       string
	TextEncoding string
	MediaTypes   []string
	ExcludeDirs  []string
	SortBy       string
}

// DefaultMediaTypes 是 Comigo 默认展示的压缩包内媒体类型。
var DefaultMediaTypes = []string{
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".svg", ".avif",
}

// ListPages 从压缩包流中扫描可阅读页面。
// source 需要支持 ReaderAt 与 Seeker，bytes.Reader、os.File、远程缓存 Reader 均可适配。
func ListPages(ctx context.Context, filename string, source io.ReaderAt, size int64, opt Options) ([]Page, error) {
	reader := io.NewSectionReader(source, 0, size)
	fsys, err := archives.FileSystem(ctx, filename, reader)
	if err != nil {
		return nil, err
	}

	pages := make([]Page, 0)
	pageNum := 1
	err = fs.WalkDir(fsys, ".", func(walkPath string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			if isSkipDir(walkPath, opt.ExcludeDirs) {
				return fs.SkipDir
			}
			return nil
		}
		if !IsSupportMedia(walkPath, opt.MediaTypes) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}
		nameInArchive := strings.TrimPrefix(path.Clean(walkPath), "./")
		if archivedFile, ok := info.(archives.FileInfo); ok && archivedFile.NameInArchive != "" {
			nameInArchive = archivedFile.NameInArchive
		}
		pages = append(pages, Page{
			Name:    nameInArchive,
			Path:    "",
			Size:    info.Size(),
			ModTime: info.ModTime(),
			URL:     pageURL(opt.BookID, nameInArchive),
			PageNum: pageNum,
		})
		pageNum++
		return nil
	})
	if err != nil {
		return nil, err
	}
	SortPages(pages, opt.SortBy)
	return pages, nil
}

// ReadPage 读取压缩包内指定文件的原始字节。
func ReadPage(ctx context.Context, filename string, source io.Reader, nameInArchive string, opt Options) ([]byte, error) {
	if nameInArchive == "" {
		return nil, errors.New("name in archive is empty")
	}
	format, stream, err := archives.Identify(ctx, filename, source)
	if err != nil {
		return nil, err
	}
	if zipFormat, ok := format.(archives.Zip); ok {
		if opt.TextEncoding != "" {
			zipFormat.TextEncoding = encoding.ByName(opt.TextEncoding)
		}
		return extractFile(ctx, zipFormat, stream, nameInArchive)
	}
	if extractor, ok := format.(archives.Extractor); ok {
		return extractFile(ctx, extractor, stream, nameInArchive)
	}
	return nil, errors.New("unsupported archive format")
}

// SortPages 按 Comigo 的页面排序规则排序。
func SortPages(pages []Page, sortBy string) {
	if sortBy == "" {
		sortBy = "default"
	}
	less := func(i, j int) bool {
		return tools.Compare(pages[i].Name, pages[j].Name)
	}
	switch sortBy {
	case "filename_reverse":
		less = func(i, j int) bool { return !tools.Compare(pages[i].Name, pages[j].Name) }
	case "filesize":
		less = func(i, j int) bool { return pages[i].Size > pages[j].Size }
	case "filesize_reverse":
		less = func(i, j int) bool { return pages[i].Size <= pages[j].Size }
	case "modify_time":
		less = func(i, j int) bool { return pages[i].ModTime.After(pages[j].ModTime) }
	case "modify_time_reverse":
		less = func(i, j int) bool { return pages[i].ModTime.Before(pages[j].ModTime) }
	}
	sort.SliceStable(pages, less)
}

// IsSupportMedia 判断压缩包内文件是否应作为页面展示。
func IsSupportMedia(checkPath string, mediaTypes []string) bool {
	base := path.Base(checkPath)
	if strings.HasPrefix(base, ".") {
		return false
	}
	if len(mediaTypes) == 0 {
		mediaTypes = DefaultMediaTypes
	}
	suffix := strings.ToLower(path.Ext(checkPath))
	for _, ext := range mediaTypes {
		if strings.ToLower(ext) == suffix {
			return true
		}
	}
	return false
}

func extractFile(ctx context.Context, extractor archives.Extractor, sourceArchive io.Reader, nameInArchive string) ([]byte, error) {
	var data []byte
	err := extractor.Extract(ctx, sourceArchive, func(ctx context.Context, f archives.FileInfo) error {
		if f.NameInArchive != nameInArchive {
			return nil
		}
		readCloser, err := f.Open()
		if err != nil {
			return err
		}
		defer readCloser.Close()
		data, err = io.ReadAll(readCloser)
		return err
	})
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, errors.New("file not found in archive")
	}
	return data, nil
}

func pageURL(bookID string, nameInArchive string) string {
	if bookID == "" {
		return ""
	}
	return "/api/get-file?id=" + bookID + "&filename=" + url.QueryEscape(nameInArchive)
}

func isSkipDir(checkPath string, excludeDirs []string) bool {
	base := filepath.Base(filepath.Clean(checkPath))
	switch base {
	case ".comigo", "flutter_ui", "node_modules":
		return true
	}
	for _, name := range excludeDirs {
		if name != "" && base == name {
			return true
		}
	}
	return false
}
