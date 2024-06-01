package file

import (
	"errors"
	"github.com/yumenaka/comi/util/locale"
	"github.com/yumenaka/comi/util/logger"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/elliotchance/orderedmap"
	"github.com/elliotchance/pie/v2"
)

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

// SaveFileToCache 读取过一次的文件，就保存到硬盘上加快读取
func SaveFileToCache(id string, filename string, data []byte, queryString string, contentType string, isCover bool, cachePath string, debug bool) error {
	err := os.MkdirAll(filepath.Join(cachePath, id), os.ModePerm)
	if err != nil {
		println(locale.GetString("saveFileToCache_error"))
	}
	//特殊字符转义，避免保存不了
	filename = url.PathEscape(filename)
	//如果是封面，另存为comigo_cover.png、comigo_cover.jpeg
	if isCover {
		filename = "comigo_cover" + path.Ext(filename)
	}
	err = os.WriteFile(filepath.Join(cachePath, id, filename), data, 0644)
	if err != nil {
		logger.Infof("%s", err)
	}
	key := cacheKey{bookID: id, queryString: queryString}
	//将ContentType存入Map
	sMap.writeMap(key, contentType)
	return err
}

// 根据query生成一个key string，用到两个第三方库
func GetQueryString(query url.Values) string {
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
func GetFileFromCache(id string, filename string, queryString string, isCover bool, cachePath string, debug bool) ([]byte, string, error) {
	//queryStringKey
	key := cacheKey{bookID: id, queryString: queryString}
	contentType, ok := sMap.readMap(key)
	if !ok {
		if debug {
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
	loadedImage, err := os.ReadFile(filepath.Join(cachePath, id, filename))
	if err != nil && debug {
		logger.Infof("getFileFromCache,file not found:%s", filename)
	}
	return loadedImage, contentType, err
}
