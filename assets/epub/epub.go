package epub

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/yumenaka/comigo/assets"
	"github.com/yumenaka/comigo/tools/logger"
)

// PageData 单个页面的数据
type PageData struct {
	Index     int    // 页码
	FileName  string // 文件名
	MediaType string // MIME 类型
	Width     int    // 图片宽度（用于 viewport）
	Height    int    // 图片高度（用于 viewport）
}

// BookData EPUB 书籍数据
type BookData struct {
	BookID       string     // 书籍唯一标识
	Title        string     // 书名
	Author       string     // 作者
	Language     string     // 语言
	ModifiedTime string     // 修改时间
	Pages        []PageData // 所有页面
}

// ImageFile 图片文件信息
type ImageFile struct {
	Name string // 文件名
	Data []byte // 文件数据
	Path string // 原始文件路径（用于从磁盘读取）
}

// Generator EPUB 生成器
type Generator struct {
	templates map[string]*template.Template
}

// NewGenerator 创建新的 EPUB 生成器
func NewGenerator() (*Generator, error) {
	g := &Generator{
		templates: make(map[string]*template.Template),
	}

	// 加载模板文件
	templateFiles := []string{
		"epub/OEBPS/content.opf.tmpl",
		"epub/OEBPS/toc.ncx.tmpl",
		"epub/OEBPS/nav.xhtml.tmpl",
		"epub/OEBPS/page.xhtml.tmpl",
	}

	for _, tmplPath := range templateFiles {
		data, err := assets.Epub.ReadFile(tmplPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read template %s: %w", tmplPath, err)
		}

		tmpl, err := template.New(filepath.Base(tmplPath)).Parse(string(data))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", tmplPath, err)
		}

		g.templates[filepath.Base(tmplPath)] = tmpl
	}

	return g, nil
}

// Generate 生成 EPUB 文件，写入到 writer
func (g *Generator) Generate(w io.Writer, bookData BookData, images []ImageFile) error {
	// 创建 zip writer
	zipWriter := zip.NewWriter(w)
	defer zipWriter.Close()

	// 1. 写入 mimetype（必须是第一个文件，不压缩）
	mimetypeData, err := assets.Epub.ReadFile("epub/mimetype")
	if err != nil {
		return fmt.Errorf("failed to read mimetype: %w", err)
	}
	if err := g.writeFileUncompressed(zipWriter, "mimetype", mimetypeData); err != nil {
		return fmt.Errorf("failed to write mimetype: %w", err)
	}

	// 2. 写入 META-INF/container.xml
	containerData, err := assets.Epub.ReadFile("epub/META-INF/container.xml")
	if err != nil {
		return fmt.Errorf("failed to read container.xml: %w", err)
	}
	if err := g.writeFile(zipWriter, "META-INF/container.xml", containerData); err != nil {
		return fmt.Errorf("failed to write container.xml: %w", err)
	}

	// 3. 写入 OEBPS/styles.css
	stylesData, err := assets.Epub.ReadFile("epub/OEBPS/styles.css")
	if err != nil {
		return fmt.Errorf("failed to read styles.css: %w", err)
	}
	if err := g.writeFile(zipWriter, "OEBPS/styles.css", stylesData); err != nil {
		return fmt.Errorf("failed to write styles.css: %w", err)
	}

	// 4. 生成 content.opf
	contentOpf, err := g.renderTemplate("content.opf.tmpl", bookData)
	if err != nil {
		return fmt.Errorf("failed to render content.opf: %w", err)
	}
	if err := g.writeFile(zipWriter, "OEBPS/content.opf", contentOpf); err != nil {
		return fmt.Errorf("failed to write content.opf: %w", err)
	}

	// 5. 生成 toc.ncx
	tocNcx, err := g.renderTemplate("toc.ncx.tmpl", bookData)
	if err != nil {
		return fmt.Errorf("failed to render toc.ncx: %w", err)
	}
	if err := g.writeFile(zipWriter, "OEBPS/toc.ncx", tocNcx); err != nil {
		return fmt.Errorf("failed to write toc.ncx: %w", err)
	}

	// 6. 生成 nav.xhtml
	navXhtml, err := g.renderTemplate("nav.xhtml.tmpl", bookData)
	if err != nil {
		return fmt.Errorf("failed to render nav.xhtml: %w", err)
	}
	if err := g.writeFile(zipWriter, "OEBPS/nav.xhtml", navXhtml); err != nil {
		return fmt.Errorf("failed to write nav.xhtml: %w", err)
	}

	// 7. 生成每个页面的 xhtml 和 写入图片
	for i, img := range images {
		// 获取图片数据
		imgData := img.Data
		if imgData == nil && img.Path != "" {
			// 从磁盘读取图片
			var readErr error
			imgData, readErr = os.ReadFile(img.Path)
			if readErr != nil {
				logger.Infof("Failed to read image %s: %v", img.Path, readErr)
				continue
			}
		}

		// 获取图片尺寸用于 viewport
		width, height := getImageDimensions(imgData)

		pageData := PageData{
			Index:     i + 1,
			FileName:  img.Name,
			MediaType: getMediaType(img.Name),
			Width:     width,
			Height:    height,
		}

		// 生成页面 xhtml
		pageXhtml, err := g.renderTemplate("page.xhtml.tmpl", pageData)
		if err != nil {
			return fmt.Errorf("failed to render page %d: %w", i+1, err)
		}
		pagePath := fmt.Sprintf("OEBPS/text/page_%d.xhtml", i+1)
		if err := g.writeFile(zipWriter, pagePath, pageXhtml); err != nil {
			return fmt.Errorf("failed to write page %d: %w", i+1, err)
		}

		// 写入图片
		imgPath := fmt.Sprintf("OEBPS/images/%s", img.Name)
		if err := g.writeFile(zipWriter, imgPath, imgData); err != nil {
			return fmt.Errorf("failed to write image %s: %w", img.Name, err)
		}
	}

	return nil
}

// renderTemplate 渲染模板
func (g *Generator) renderTemplate(name string, data interface{}) ([]byte, error) {
	tmpl, ok := g.templates[name]
	if !ok {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeFile 写入文件到 zip（使用压缩）
func (g *Generator) writeFile(zw *zip.Writer, path string, data []byte) error {
	header := &zip.FileHeader{
		Name:     path,
		Method:   zip.Deflate,
		Modified: time.Now(),
	}
	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// writeFileUncompressed 写入文件到 zip（不压缩，用于 mimetype）
func (g *Generator) writeFileUncompressed(zw *zip.Writer, path string, data []byte) error {
	header := &zip.FileHeader{
		Name:   path,
		Method: zip.Store, // 不压缩
	}
	w, err := zw.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// getMediaType 根据文件扩展名获取 MIME 类型
func getMediaType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		// 默认图片类型
		switch ext {
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		case ".webp":
			return "image/webp"
		case ".svg":
			return "image/svg+xml"
		default:
			return "application/octet-stream"
		}
	}
	return mimeType
}

// getImageDimensions 从图片数据获取尺寸
func getImageDimensions(data []byte) (width, height int) {
	// 默认尺寸（当无法解析图片时使用）
	defaultWidth, defaultHeight := 1600, 2400

	if len(data) == 0 {
		return defaultWidth, defaultHeight
	}

	reader := bytes.NewReader(data)
	config, _, err := image.DecodeConfig(reader)
	if err != nil {
		logger.Infof("Failed to decode image config: %v, using default dimensions", err)
		return defaultWidth, defaultHeight
	}

	return config.Width, config.Height
}

// CreateBookData 从书籍信息创建 BookData
func CreateBookData(bookID, title, author string, imageFiles []ImageFile) BookData {
	pages := make([]PageData, len(imageFiles))
	for i, img := range imageFiles {
		pages[i] = PageData{
			Index:     i + 1,
			FileName:  img.Name,
			MediaType: getMediaType(img.Name),
		}
	}

	// 设置默认语言
	language := "en"

	return BookData{
		BookID:       bookID,
		Title:        title,
		Author:       author,
		Language:     language,
		ModifiedTime: time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		Pages:        pages,
	}
}
