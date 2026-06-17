package player

import (
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/templ/common"
	"github.com/yumenaka/comigo/templ/pages/error_page"
	"github.com/yumenaka/comigo/tools/comigo_remote"
	"github.com/yumenaka/comigo/tools/logger"
)

type PlayerMediaItem struct {
	ID       string                `json:"id"`
	Title    string                `json:"title"`
	Type     model.SupportFileType `json:"type"`
	RawURL   string                `json:"rawUrl"`
	CoverURL string                `json:"coverUrl"`
	MimeType string                `json:"mimeType"`
}

// PlayerData 是播放器页面专用的前端数据载荷。
// 只输出可公开访问的 URL 与 MIME，不把 BookPath/StoreUrl 这类内部路径交给浏览器。
type PlayerData struct {
	Current  PlayerMediaItem   `json:"current"`
	Playlist []PlayerMediaItem `json:"playlist"`
}

// PlayerModeHandler 播放器页面处理器
func PlayerModeHandler(c echo.Context) error {
	// 读取url参数，获取书籍ID
	bookID := c.Param("id")
	var book *model.Book
	if remoteBook, ok, remoteErr := comigo_remote.ResolveBook(config.GetCfg().StoreUrls, c.QueryParam(comigo_remote.RemoteStoreQuery), bookID, "default", config.GetCfg().TimeoutLimitForScan); ok {
		if remoteErr != nil {
			logger.Infof(locale.GetString("log_getbook_error_scroll"), remoteErr)
			return renderPlayerNotFound(c)
		}
		book = remoteBook
	} else {
		var err error
		book, err = model.IStore.GetBook(bookID)
		if err != nil {
			logger.Infof(locale.GetString("log_getbook_error_scroll"), err)
			return renderPlayerNotFound(c)
		}
	}
	if !isMediaBook(book.BookInfo) {
		return renderPlayerNotFound(c)
	}

	// 获取同文件夹下的所有书籍作为播放列表
	playlist := common.QuickJumpBarBooks(book)

	// 过滤出当前目录下的音视频文件。
	mediaList := buildMediaPlaylist(playlist)
	playerData := buildPlayerData(book, mediaList)

	// 定义模板主体内容。
	playerPage := PlayerPage(c, book, playerData)
	// 拼接页面
	indexHtml := common.Html(
		c,
		playerPage, // define body content
		[]string{},
	)
	// 渲染页面
	return common.RenderHTML(c, indexHtml)
}

func renderPlayerNotFound(c echo.Context) error {
	// 没有找到可播放媒体，显示 HTTP 404 错误
	indexHtml := common.Html(
		c,
		error_page.NotFound404(c),
		[]string{},
	)
	c.Response().WriteHeader(http.StatusNotFound)
	return common.RenderHTML(c, indexHtml)
}

// buildMediaPlaylist 从同目录书籍列表中提取播放器列表。
func buildMediaPlaylist(playlist *model.BookInfos) model.BookInfos {
	if playlist == nil {
		return nil
	}

	mediaList := make(model.BookInfos, 0, len(*playlist))
	for _, item := range *playlist {
		if item.Type != model.TypeVideo && item.Type != model.TypeAudio {
			continue
		}
		mediaList = append(mediaList, item)
	}
	return mediaList
}

// buildPlayerData 组装播放器专用数据，确保当前媒体一定存在于播放列表中。
func buildPlayerData(book *model.Book, mediaList model.BookInfos) PlayerData {
	current := buildPlayerMediaItem(book.BookInfo)
	items := make([]PlayerMediaItem, 0, len(mediaList))
	hasCurrent := false
	for _, item := range mediaList {
		if !isMediaBook(item) {
			continue
		}
		mediaItem := buildPlayerMediaItem(item)
		if mediaItem.ID == current.ID {
			hasCurrent = true
		}
		items = append(items, mediaItem)
	}
	if !hasCurrent {
		items = append(items, current)
	}
	return PlayerData{
		Current:  current,
		Playlist: items,
	}
}

// buildPlayerMediaItem 将内部 BookInfo 转为前端可直接使用的公开媒体项。
func buildPlayerMediaItem(book model.BookInfo) PlayerMediaItem {
	return PlayerMediaItem{
		ID:       book.BookID,
		Title:    book.Title,
		Type:     book.Type,
		RawURL:   common.RawBookInfoURL(book),
		CoverURL: playerCoverURL(book),
		MimeType: mediaMimeType(book),
	}
}

// isMediaBook 判断书籍类型是否能进入播放器页面。
func isMediaBook(book model.BookInfo) bool {
	return book.Type == model.TypeVideo || book.Type == model.TypeAudio
}

// playerCoverURL 统一生成播放器封面 URL，避免前端依赖本地路径。
func playerCoverURL(book model.BookInfo) string {
	coverURL := config.PrefixPath("/api/get-cover?id=" + url.QueryEscape(book.BookID) + "&resize_height=352")
	if book.RemoteStoreKey != "" {
		coverURL += "&remote_store=" + url.QueryEscape(book.RemoteStoreKey)
	}
	return coverURL
}

// mediaMimeType 根据标题扩展名生成浏览器 source type，减少音视频探测失败。
func mediaMimeType(book model.BookInfo) string {
	fileName := book.Title
	if fileName == "" {
		fileName = filepath.Base(book.BookPath)
	}
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".m4a":
		return "audio/mp4"
	case ".aac":
		return "audio/aac"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".wma":
		return "audio/x-ms-wma"
	case ".mp4":
		return "video/mp4"
	case ".m4v":
		return "video/x-m4v"
	case ".mov":
		return "video/quicktime"
	case ".webm":
		return "video/webm"
	case ".avi":
		return "video/x-msvideo"
	case ".flv":
		return "video/x-flv"
	}
	return mime.TypeByExtension(ext)
}
