package arch

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/yumenaka/archiver/v4"
	"github.com/yumenaka/comi/logger"
	"golang.org/x/net/html"
)

// 获取epub数据（html、xml等）
func getDataFromEpub(epubPath string, needFile string) (data []byte, err error) {
	//必须传值
	if needFile == "" {
		return nil, errors.New("needFile is empty")
	}
	//打开文件，只读模式
	file, err := os.OpenFile(epubPath, os.O_RDONLY, 0400) //Use mode 0400 for a read-only // file and 0600 for a readable+writable file.
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Info("file.Close() Error:", err)
		}
	}(file)
	//是否是压缩包
	format, _, err := archiver.Identify(epubPath, file)
	if err != nil {
		return nil, err
	}
	//如果是epub文件,文件编码为UTF-8
	if ex, ok := format.(archiver.Zip); ok {
		//特殊编码
		ex.TextEncoding = "" // “”
		ctx := context.Background()
		//这里是file，而不是sourceArchive，否则会出错。
		err := ex.Extract(ctx, file, []string{needFile}, func(ctx context.Context, f archiver.File) error {
			// 取得特定压缩文件
			file, err := f.Open()
			if err != nil {
				return err
			}
			defer func(file io.ReadCloser) {
				err := file.Close()
				if err != nil {
					logger.Info("file.Close() Error:", err)
				}
			}(file)
			content, err := io.ReadAll(file)
			if err != nil {
				return err
			}
			data = content
			return err
		})
		return data, err
	}
	return nil, errors.New("getDataFromEpub Error. epubPath:" + epubPath + "  needFile:" + needFile)
}

// Container 定义结构体、映射xml结构 was generated 2022-12-09 00:41:31 by https://xml-to-go.github.io/ in Ukraine.
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
	//data, err := os.ReadFile(ContainerXMLPath)
	//if err != nil {
	//	logger.Info("ReadFile Error:", err)
	//}
	data, err := getDataFromEpub(epubPath, "META-INF/container.xml")
	if err != nil {
		logger.Info(err)
		return "", errors.New("getOPFPath Error epubPath:" + epubPath)
	}
	con := new(Container)
	err = xml.Unmarshal(data, con)
	if err != nil {
		logger.Info("XML Unmarshal Error:", err)
	}
	opfPath = con.Rootfiles.Rootfile.FullPath
	return
}

// 获取html文件里的第一个img标签
func findAttrValue(r io.Reader, imgKey string) (value string) {
	//NewTokenizer 为给定的 Reader 返回一个新的 HTML 分词器（Tokenizer）。假定输入是 UTF-8 编码
	tokenizer := html.NewTokenizer(r)
	for tokenizer.Token().Data != "html" {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return
			}
			logger.Infof("Error: %v", tokenizer.Err())
			return
		}
		tagName, _ := tokenizer.TagName()
		//第一个img标签
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
		logger.Info("getOPFPath Error:", err)
		return
	}
	b, err := getDataFromEpub(epubPath, opfPath)
	if err != nil {
		logger.Info("getDataFromEpub Error:", err)
		return
	}
	err = xml.Unmarshal(b, pack)
	if err != nil {
		logger.Info("XML Unmarshal Error:", err)
		return
	}
	//顺序信息
	itemRef := pack.Spine.Itemref
	//资源列表
	item := pack.Manifest.Item
	//获取Spine里面排好序的HTML文件ID，并生成Html文件列表
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
	//根据有序的html列表，解析html，生成有序图片列表，供其他模块排序用
	for i := 0; i < len(htmlList); i++ {
		data, err := getDataFromEpub(epubPath, htmlList[i])
		if err != nil {
			logger.Info(err)
			continue
		}
		reader := bytes.NewReader(data)
		tempSrc := findAttrValue(reader, "src")
		src := absUrl(tempSrc, htmlList[i])
		imageList = append(imageList, src)
		//logger.Info(src)
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
		logger.Info("getOPFPath Error:", err)
		return
	}
	b, err := getDataFromEpub(epubPath, opfPath)
	if err != nil {
		logger.Info("getDataFromEpub Error:", err)
		return
	}
	err = xml.Unmarshal(b, pack)
	if err != nil {
		logger.Info("XML Unmarshal Error:", err)
		return
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
