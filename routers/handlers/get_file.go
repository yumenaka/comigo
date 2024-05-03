package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/elliotchance/orderedmap"
	"github.com/elliotchance/pie/v2"
	"github.com/gin-gonic/gin"
	"github.com/yumenaka/comi/arch"
	"github.com/yumenaka/comi/config"
	"github.com/yumenaka/comi/locale"
	"github.com/yumenaka/comi/logger"
	"github.com/yumenaka/comi/types"
)

// GetFile 示例 URL： 127.0.0.1:1234/get_file?id=2b17a13&filename=1.jpg
// 缩放文件，会转化为jpeg：http://127.0.0.1:1234/api/get_file?resize_width=300&resize_height=400&id=597e06&filename=01.jpeg
// 相关参数：
// id：书籍的ID，必须项目       							&id=2B17a
// filename:获取的文件名，必须项目   							&filename=01.jpg
// //可选参数：
// resize_width:指定宽度，缩放图片  							&resize_width=300
// resize_height:指定高度，缩放图片 							&resize_height=300
// thumbnail_mode:缩略图模式，同时指定宽高的时候要不要剪切图片		&thumbnail_mode=true
// resize_max_width:指定宽度上限，图片宽度大于这个上限时缩小图片  	&resize_max_width=740
// resize_max_height:指定高度上限，图片高度大于这个上限时缩小图片  	&resize_max_height=300
// auto_crop:自动切白边，数字是切白边的阈值，范围是0~100 			&auto_crop=10
// gray:黑白化												&gray=true
// blurhash:获取对应图片的blurhash，不是原始图片 				&blurhash=3
// blurhash_image:获取对应图片的blurhash图片，不是原始图片  	    &blurhash_image=3
func GetFile(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	needFile := c.DefaultQuery("filename", "")
	//没有指定这两项，直接返回
	if id == "" && needFile == "" {
		return
	}
	noCache := c.DefaultQuery("no-cache", "false")
	//如果启用了本地缓存
	if config.Config.UseCache && noCache == "false" {
		//获取所有的参数键值对
		query := c.Request.URL.Query()
		//如果有缓存，直接读取本地获取缓存文件并返回
		cacheData, ct, errGet := getFileFromCache(id, needFile, query, c.DefaultQuery("thumbnail_mode", "false") == "true")
		if errGet == nil && cacheData != nil {
			c.Data(http.StatusOK, ct, cacheData)
			return
		}
	}

	bookByID, err := types.GetBookByID(id, "")
	if err != nil {
		logger.Infof("%s", err)
	}

	//读取图片Resize用的resizeWidth
	resizeWidth, errX := strconv.Atoi(c.DefaultQuery("resize_width", "0"))
	if errX != nil {
		resizeWidth = 0
	}
	//读取图片Resize用的resizeHeight
	resizeHeight, errY := strconv.Atoi(c.DefaultQuery("resize_height", "0"))
	if errY != nil {
		resizeHeight = 0
	}
	//自动切白边参数
	autoCrop, errCrop := strconv.Atoi(c.DefaultQuery("auto_crop", "-1"))
	if errCrop != nil {
		autoCrop = -1
	}
	//图片Resize, 按照 maxWidth 限制大小
	resizeMaxWidth, errMX := strconv.Atoi(c.DefaultQuery("resize_max_width", "0"))
	if errMX != nil {
		resizeMaxWidth = 0
	}
	//图片Resize, 按照 MaxHeight 限制大小
	resizeMaxHeight, errMY := strconv.Atoi(c.DefaultQuery("resize_max_height", "0"))
	if errMY != nil {
		resizeMaxHeight = 0
	}

	blurhash, blurErr := strconv.Atoi(c.DefaultQuery("blurhash", "0"))
	if blurErr != nil {
		blurhash = 0
	}

	blurhashImage, blurImageErr := strconv.Atoi(c.DefaultQuery("blurhash_image", "0"))
	if blurImageErr != nil {
		blurhashImage = 0
	}

	option := arch.GetPictureDataOption{
		PictureName:      needFile,
		BookIsPDF:        bookByID.Type == types.TypePDF,
		BookIsDir:        bookByID.Type == types.TypeDir,
		BookIsNonUTF8Zip: bookByID.NonUTF8Zip,
		BookFilePath:     bookByID.FilePath,
		Debug:            config.Config.Debug,
		UseCache:         config.Config.UseCache,
		ResizeWidth:      resizeWidth,
		ResizeHeight:     resizeHeight,
		ResizeMaxWidth:   resizeMaxWidth,
		ResizeMaxHeight:  resizeMaxHeight,
		ThumbnailMode:    c.DefaultQuery("thumbnail_mode", "false") == "true",
		AutoCrop:         autoCrop,
		Gray:             c.DefaultQuery("gray", "false") == "true",
		BlurHash:         blurhash,
		BlurHashImage:    blurhashImage,
	}
	imgData, contentType, err := arch.GetPictureData(option)
	if err != nil {
		c.String(http.StatusBadRequest, "GetPictureData error:%s", err)
	}

	//如果启用了本地缓存
	if config.Config.UseCache && noCache == "false" && bookByID.Type != types.TypeDir {
		//获取所有的参数键值对
		query := c.Request.URL.Query()
		//缓存文件到本地，避免重复解压
		errSave := saveFileToCache(id, needFile, imgData, query, contentType, c.DefaultQuery("thumbnail_mode", "false") == "true")
		if errSave != nil {
			logger.Info(errSave)
		}
	}
	c.Data(http.StatusOK, contentType, imgData)
}

// 储存文件信息的key
type cacheKey struct {
	bookID      string
	queryString string
}

// SyncMap 有读写锁的map
type SyncMap struct {
	sync.RWMutex                       // map不是并发安全的 , 当有多个并发的goroutine读写同一个map时会出现panic错误(fatal error: concurrent map writes)
	mapContentType map[cacheKey]string //需要一个Map保存ContentType
}

// 读写锁
func (l *SyncMap) readMap(key cacheKey) (string, bool) {
	l.RLock()
	value, ok := l.mapContentType[key]
	l.RUnlock()
	return value, ok
}

// 读写锁
func (l *SyncMap) writeMap(key cacheKey, value string) {
	l.Lock()
	l.mapContentType[key] = value
	l.Unlock()
}

// SyncMap 有读写锁的map.除此之外，还可以使用channel，或sync.map保证map的并发安全
var sMap SyncMap

func init() {
	sMap.mapContentType = make(map[cacheKey]string)
}

// 读取过一次的文件，就保存到硬盘上加快读取
func saveFileToCache(id string, filename string, data []byte, query url.Values, contentType string, isCover bool) error {
	err := os.MkdirAll(filepath.Join(config.Config.CachePath, id), os.ModePerm)
	if err != nil {
		println(locale.GetString("saveFileToCache_error"))
	}
	//特殊字符转义，避免保存不了
	filename = url.PathEscape(filename)
	//如果是封面，另存为comigo_cover.png、comigo_cover.jpeg
	if isCover {
		filename = "comigo_cover" + path.Ext(filename)
	}
	err = os.WriteFile(filepath.Join(config.Config.CachePath, id, filename), data, 0644)
	if err != nil {
		logger.Infof("%s", err)
	}
	qS := getQueryStringKey(query)
	key := cacheKey{bookID: id, queryString: qS}
	//将ContentType存入Map
	sMap.writeMap(key, contentType)
	return err
}

// 根据query生成一个key string，用到两个第三方库
func getQueryStringKey(query url.Values) string {
	//因为map没有排序，相同参数每次形成的string都不一样,所以需要第三方库，建立一个有序map。
	//OrderedMap按照插入顺序排序迭代，所以插入的时候也要保证顺序
	m := orderedmap.NewOrderedMap()
	//构建一个key列表，并用pie排序
	var keyList []string
	for k := range query {
		keyList = append(keyList, k)
	}
	//pie.Sort()返回一个排好序的slice
	sortKeyList := pie.Sort(keyList)
	//按照顺序插入
	for _, sortKey := range sortKeyList {
		m.Set(sortKey, query[sortKey])
	}
	queryString := ""
	//取出values与key，然后通过类型断言转换成原类型 Keys()按照插入顺序迭代
	for _, key := range m.Keys() {
		values, _ := m.Get(key)
		//go 类型断言
		if V, ok := values.([]string); ok {
			temp := ""
			for _, v := range V {
				temp = temp + v
			}
			//go 类型断言
			if K, ok := key.(string); ok {
				queryString = queryString + K + temp
			}
		}
	}
	//logger.Infof("queryString:" + queryString)
	return queryString
}

// 读取缓存，加快第二次访问的速度
func getFileFromCache(id string, filename string, query url.Values, isCover bool) ([]byte, string, error) {
	//将query转换为字符串，后面用来当key
	qS := getQueryStringKey(query)
	key := cacheKey{bookID: id, queryString: qS}
	contentType, ok := sMap.readMap(key)
	if !ok {
		if config.Config.Debug {
			return nil, contentType, errors.New("getFileFromCache key not found")
		}
		return nil, contentType, nil
	}
	//文件名经过转义，避免保存不了，所以这里也必须转义才能取到本地文件
	filename = url.PathEscape(filename)
	//如果是封面，另存为comigo_cover.png、comigo_cover.jpeg
	if isCover {
		filename = "comigo_cover" + path.Ext(filename)
	}
	loadedImage, err := os.ReadFile(filepath.Join(config.Config.CachePath, id, filename))
	if err != nil && config.Config.Debug {
		logger.Infof("getFileFromCache,file not found:%s", filename)
	}
	return loadedImage, contentType, err
}
