package file

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yumenaka/archives"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/tools/logger"
	"golang.org/x/net/html"
)

// 获取epub数据（html、xml等）
func getDataFromEpub(epubPath string, needFile string) (data []byte, err error) {
	// 必须传值
	if needFile == "" {
		return nil, errors.New(locale.GetString("err_needfile_empty"))
	}
	// 打开文件，只读模式
	file, err := os.OpenFile(epubPath, os.O_RDONLY, 0o400) // Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Infof(locale.GetString("log_file_close_error"), err)
		}
	}(file)
	// 是否是压缩包
	format, _, err := archives.Identify(context.Background(), epubPath, file)
	if err != nil {
		return nil, err
	}
	// 如果是epub文件,文件编码为UTF-8，不需要特殊处理。
	if ex, ok := format.(archives.Zip); ok {
		// 特殊编码
		// ex.TextEncoding = tools.GetEncodingByName(textEncoding)
		ctx := context.Background()
		found := false
		// 这这个用法，只根据file名获取，不加文件夹内二级路径
		err := ex.Extract(ctx, file, func(ctx context.Context, f archives.FileInfo) error {
			// 检查是否是需要的文件
			if f.Name() != needFile {
				return nil
			}
			// logger.Info("file.Name():"+f.Name(), "needFile:", needFile)
			found = true
			// 取得特定压缩文件
			file, err := f.Open()
			if err != nil {
				return err
			}
			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			data = content
			file.Close()
			return nil
		})
		if err != nil {
			return nil, err
		}
		if !found {
			return nil, fmt.Errorf("file not found in epub: %s", needFile)
		}
		return data, nil
	}
	return nil, fmt.Errorf(locale.GetString("err_getdata_from_epub_error"), epubPath, needFile)
}

// Container was generated 2025-04-15 23:51:38 by https://xml-to-go.github.io/ in Ukraine.
type Container struct {
	XMLName   xml.Name `xml:"container"`
	Text      string   `xml:",chardata"`
	Version   string   `xml:"version,attr"`
	Xmlns     string   `xml:"xmlns,attr"`
	Rootfiles struct {
		Text     string `xml:",chardata"`
		Rootfile struct {
			Text      string `xml:",chardata"`
			FullPath  string `xml:"full-path,attr"`
			MediaType string `xml:"media-type,attr"`
		} `xml:"rootfile"`
	} `xml:"rootfiles"`
}

// Package 定义结构体、映射OPF文件（xml）结构用。
// was generated 2022-12-09 00:47:41 by https://xml-to-go.github.io/ in Ukraine.
type Package struct {
	XMLName          xml.Name `xml:"package"`
	Text             string   `xml:",chardata"`
	Version          string   `xml:"version,attr"`
	UniqueIdentifier string   `xml:"unique-identifier,attr"`
	Xmlns            string   `xml:"xmlns,attr"`
	Metadata         struct {
		Text string `xml:",chardata"`
		Dc   string `xml:"dc,attr"`
		Opf  string `xml:"opf,attr"`
		Meta []struct {
			Text    string `xml:",chardata"`
			Name    string `xml:"name,attr"`
			Content string `xml:"content,attr"`
		} `xml:"meta"`
		Identifier struct {
			Text   string `xml:",chardata"`
			ID     string `xml:"id,attr"`
			Scheme string `xml:"scheme,attr"`
		} `xml:"identifier"`
		Title     string `xml:"title"`
		Language  string `xml:"language"`
		Creator   string `xml:"creator"`
		Publisher string `xml:"publisher"`
		Date      string `xml:"date"`
		Rights    string `xml:"rights"`
		Series    string `xml:"series"`
	} `xml:"metadata"`
	Manifest struct {
		Text string `xml:",chardata"`
		Item []struct {
			Text       string `xml:",chardata"`
			ID         string `xml:"id,attr"`
			Href       string `xml:"href,attr"`
			MediaType  string `xml:"media-type,attr"`
			Properties string `xml:"properties,attr"`
		} `xml:"item"`
	} `xml:"manifest"`
	Spine struct {
		Text    string `xml:",chardata"`
		Toc     string `xml:"toc,attr"`
		Itemref []struct {
			Text  string `xml:",chardata"`
			Idref string `xml:"idref,attr"`
		} `xml:"itemref"`
	} `xml:"spine"`
	Guide struct {
		Text      string `xml:",chardata"`
		Reference struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Href  string `xml:"href,attr"`
			Title string `xml:"title,attr"`
		} `xml:"reference"`
	} `xml:"guide"`
}

// 获取OPF文件的路径
func getOPFPath(epubPath string) (opfPath string, err error) {
	data, err := getDataFromEpub(epubPath, filepath.Base("META-INF/container.xml"))
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_container_xml"), err)
		return "", fmt.Errorf("getOPFPath Error: %w", err)
	}
	if len(data) == 0 {
		return "", errors.New(locale.GetString("err_container_xml_empty"))
	}
	con := new(Container)
	err = xml.Unmarshal(data, con)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_parse_container_xml"), err)
		return "", fmt.Errorf("XML Unmarshal Error: %w", err)
	}
	opfPath = con.Rootfiles.Rootfile.FullPath
	if opfPath == "" {
		return "", errors.New(locale.GetString("err_no_valid_opf_path"))
	}
	return
}

// 获取html文件里的第一个img标签
func findAttrValue(r io.Reader, imgKey string) (value string) {
	// NewTokenizer 为给定的 Reader 返回一个新的 HTML 分词器（Tokenizer）。假定输入是 UTF-8 编码
	tokenizer := html.NewTokenizer(r)
	for tokenizer.Token().Data != "html" {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return
			}
			logger.Infof(locale.GetString("log_html_tokenizer_error"), tokenizer.Err())
			return
		}
		tagName, _ := tokenizer.TagName()
		// 第一个img标签
		if string(tagName) == "img" {
			attrKey, attrValue, _ := tokenizer.TagAttr()
			if string(attrKey) == imgKey {
				return string(attrValue)
			}
		}
	}
	return
}

// 返回image src相对于根路径的地址
func absUrl(currUrl, baseUrl string) string {
	urlInfo, err := url.Parse(currUrl)
	if err != nil {
		return ""
	}
	if urlInfo.Scheme != "" {
		return currUrl
	}
	baseInfo, err := url.Parse(baseUrl)
	if err != nil {
		return ""
	}
	u := ""
	var path string
	if strings.Index(urlInfo.Path, "/") == 0 {
		path = urlInfo.Path
	} else {
		path = filepath.Dir(baseInfo.Path) + "/" + urlInfo.Path
	}
	rst := make([]string, 0)
	pathArr := strings.Split(path, "/")

	// 如果path是已/开头，那在rst加入一个空元素
	if pathArr[0] == "" {
		rst = append(rst, "")
	}
	for _, p := range pathArr {
		if p == ".." {
			if rst[len(rst)-1] == ".." {
				rst = append(rst, "..")
			} else {
				rst = rst[:len(rst)-1]
			}
		} else if p != "" && p != "." {
			rst = append(rst, p)
		}
	}
	return u + strings.Join(rst, "/")
}

// GetImageListFromEpubFile 根据Epub信息，获取有序的imgSrc列表
func GetImageListFromEpubFile(epubPath string) (imageList []string, err error) {
	pack := new(Package)
	opfPath, err := getOPFPath(epubPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get OPF file path: %w", err)
	}
	b, err := getDataFromEpub(epubPath, filepath.Base(opfPath))
	if err != nil {
		return nil, fmt.Errorf("failed to read opf file: %w", err)
	}
	err = xml.Unmarshal(b, pack)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OPF files: %w", err)
	}
	// 顺序信息
	itemRef := pack.Spine.Itemref
	// 资源列表
	item := pack.Manifest.Item
	// 获取Spine里面排好序的HTML文件ID，并生成Html文件列表
	var htmlList []string
	for i := 0; i < len(itemRef); i++ {
		id := itemRef[i].Idref
		for j := 0; j < len(item); j++ {
			if item[j].ID == id {
				if item[j].MediaType == "application/xhtml+xml" {
					htmlList = append(htmlList, item[j].Href)
				}
			}
		}
	}
	// 根据有序的html列表，解析html，生成有序图片列表，供其他模块排序用
	for i := 0; i < len(htmlList); i++ {
		data, err := getDataFromEpub(epubPath, filepath.Base(htmlList[i]))
		if err != nil {
			logger.Infof("%s", err)
			continue
		}
		reader := bytes.NewReader(data)
		tempSrc := findAttrValue(reader, "src")
		src := absUrl(tempSrc, htmlList[i])
		imageList = append(imageList, src)
		// logger.Infof(src)
	}
	return imageList, err
}

type EpubMetadata struct {
	Title     string `xml:"title"`
	Language  string `xml:"language"`
	Creator   string `xml:"creator"`
	Publisher string `xml:"publisher"`
	Date      string `xml:"date"`
	Rights    string `xml:"rights"`
	Series    string `xml:"series"`
}

// GetEpubMetadata 根据Epub信息，获取书籍详情
func GetEpubMetadata(epubPath string) (metadata EpubMetadata, err error) {
	pack := new(Package)
	opfPath, err := getOPFPath(epubPath)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_get_opf_file_path"), err)
		return EpubMetadata{}, fmt.Errorf("获取元数据失败: %w", err)
	}
	b, err := getDataFromEpub(epubPath, filepath.Base(opfPath))
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_read_opf_file"), err)
		return EpubMetadata{}, fmt.Errorf("读取OPF文件失败: %w", err)
	}
	err = xml.Unmarshal(b, pack)
	if err != nil {
		logger.Infof(locale.GetString("log_failed_to_parse_opf_file"), err)
		return EpubMetadata{}, fmt.Errorf("解析OPF文件失败: %w", err)
	}
	return EpubMetadata{
		Title:     pack.Metadata.Title,
		Language:  pack.Metadata.Language,
		Creator:   pack.Metadata.Creator,
		Publisher: pack.Metadata.Publisher,
		Date:      pack.Metadata.Date,
		Rights:    pack.Metadata.Rights,
		Series:    pack.Metadata.Series,
	}, nil
}
