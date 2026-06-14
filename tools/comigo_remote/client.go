package comigo_remote

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/jxskiss/base62"
	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

const RemoteStoreQuery = "remote_store"

// Client 访问另一个 Comigo 服务，负责登录态和 BasePath 拼接。
type Client struct {
	baseURL  *url.URL
	username string
	password string
	http     *http.Client
}

type successBody struct {
	Data json.RawMessage `json:"data"`
}

// ShelfKey 生成远端顶级书库的本地内部标识。
// 旧版 Comigo top-shelf 不输出真实 StoreUrl，只能用远端主页、顺序和显示名稳定区分。
func ShelfKey(storeURL string, shelf model.StoreBookInfo, index int) string {
	publicBase, err := PublicBaseURL(storeURL)
	if err != nil {
		publicBase = storeURL
	}
	raw := fmt.Sprintf("%s\x00%d\x00%s", publicBase, index, shelf.DisplayName)
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// PublicBaseURL 返回不含账号密码的远端 Comigo 主页 URL。
func PublicBaseURL(storeURL string) (string, error) {
	parsedURL, err := url.Parse(storeURL)
	if err != nil {
		return "", err
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("not a Comigo URL: %s", storeURL)
	}
	parsedURL.User = nil
	parsedURL.RawQuery = ""
	parsedURL.Fragment = ""
	parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	if parsedURL.Path == "/" {
		parsedURL.Path = ""
	}
	return parsedURL.String(), nil
}

// StoreKey 生成公开、稳定、短于原始 URL 的远端书库标识。
func StoreKey(storeURL string) (string, error) {
	publicBase, err := PublicBaseURL(storeURL)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString([]byte(publicBase)), nil
}

// MatchStoreByKey 从当前配置的书库列表中找到 remote_store 对应的远端 Comigo URL。
func MatchStoreByKey(storeURLs []string, key string) (string, error) {
	for _, storeURL := range storeURLs {
		gotKey, err := StoreKey(storeURL)
		if err == nil && gotKey == key {
			return storeURL, nil
		}
	}
	return "", fmt.Errorf("remote store not found: %s", key)
}

// LocalBookID 根据远端书库和远端 BookID 生成本机 BookID，避免跨书库冲突。
func LocalBookID(storeURL string, remoteBookID string) string {
	key, err := StoreKey(storeURL)
	if err != nil {
		key = storeURL
	}
	fullID := base62.EncodeToString([]byte(tools.Md5string(tools.Md5string(key + "\x00" + remoteBookID))))
	if len(fullID) > 12 {
		return fullID[:12]
	}
	return fullID
}

// NewClient 创建远端 Comigo 客户端。URL 中的 userinfo 用于远端登录。
func NewClient(storeURL string, timeoutSeconds int) (*Client, error) {
	parsedURL, err := url.Parse(storeURL)
	if err != nil {
		return nil, err
	}
	publicBase, err := PublicBaseURL(storeURL)
	if err != nil {
		return nil, err
	}
	baseURL, err := url.Parse(publicBase)
	if err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(timeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 20 * time.Second
	}
	client := &Client{
		baseURL: baseURL,
		http: &http.Client{
			Jar:     jar,
			Timeout: timeout,
		},
	}
	if parsedURL.User != nil {
		client.username = parsedURL.User.Username()
		client.password, _ = parsedURL.User.Password()
	}
	return client, nil
}

// GetTopShelf 获取远端顶层书架。
func (c *Client) GetTopShelf(sortBy string) ([]model.StoreBookInfo, error) {
	query := url.Values{}
	if sortBy != "" {
		query.Set("sort_by", sortBy)
	}
	data, _, err := c.doBytes(http.MethodGet, "/api/top-shelf", query, nil, "")
	if err != nil {
		return nil, err
	}
	var shelves []model.StoreBookInfo
	return shelves, json.Unmarshal(data, &shelves)
}

// GetBook 获取远端书籍详情。
func (c *Client) GetBook(remoteBookID string, sortBy string) (*model.Book, error) {
	query := url.Values{}
	query.Set("id", remoteBookID)
	if sortBy != "" {
		query.Set("sort_by", sortBy)
	}
	data, _, err := c.doBytes(http.MethodGet, "/api/get-book", query, nil, "")
	if err != nil {
		return nil, err
	}
	var body successBody
	if err := json.Unmarshal(data, &body); err != nil {
		return nil, err
	}
	if len(body.Data) == 0 {
		return nil, errors.New("remote get-book response has no data")
	}
	var book model.Book
	if err := json.Unmarshal(body.Data, &book); err != nil {
		return nil, err
	}
	return &book, nil
}

// GetAllBookmarks 获取远端所有书签。
func (c *Client) GetAllBookmarks() ([]model.BookinfoWithBookMark, error) {
	data, _, err := c.doBytes(http.MethodGet, "/api/all-bookmarks", nil, nil, "")
	if err != nil {
		return nil, err
	}
	var bookmarks []model.BookinfoWithBookMark
	return bookmarks, json.Unmarshal(data, &bookmarks)
}

// GetBytes 代理远端二进制 API。
func (c *Client) GetBytes(apiPath string, query url.Values) ([]byte, string, error) {
	return c.doBytes(http.MethodGet, apiPath, query, nil, "")
}

// PostJSON 代理远端 JSON POST。
func (c *Client) PostJSON(apiPath string, payload []byte) ([]byte, string, error) {
	return c.doBytes(http.MethodPost, apiPath, nil, payload, "application/json")
}

// Delete 代理远端 DELETE。
func (c *Client) Delete(apiPath string, query url.Values) ([]byte, string, error) {
	return c.doBytes(http.MethodDelete, apiPath, query, nil, "")
}

// GetResponse 流式代理远端 GET 响应，供音视频等需要 Range 的原始资源使用。
func (c *Client) GetResponse(apiPath string, query url.Values, headers http.Header) (*http.Response, error) {
	if c.username != "" {
		if err := c.login(); err != nil {
			return nil, err
		}
	}
	resp, err := c.doRequest(http.MethodGet, apiPath, query, nil, "", headers)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusUnauthorized && c.username != "" {
		resp.Body.Close()
		if err := c.login(); err != nil {
			return nil, err
		}
		resp, err = c.doRequest(http.MethodGet, apiPath, query, nil, "", headers)
		if err != nil {
			return nil, err
		}
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		return nil, fmt.Errorf("remote Comigo GET %s failed: %s", apiPath, resp.Status)
	}
	return resp, nil
}

func (c *Client) doBytes(method string, apiPath string, query url.Values, body []byte, contentType string) ([]byte, string, error) {
	if c.username != "" {
		if err := c.login(); err != nil {
			return nil, "", err
		}
	}
	data, responseContentType, status, err := c.doOnce(method, apiPath, query, body, contentType)
	if status != http.StatusUnauthorized || c.username == "" {
		return data, responseContentType, err
	}
	if err := c.login(); err != nil {
		return nil, "", err
	}
	data, responseContentType, _, err = c.doOnce(method, apiPath, query, body, contentType)
	return data, responseContentType, err
}

func (c *Client) doOnce(method string, apiPath string, query url.Values, body []byte, contentType string) ([]byte, string, int, error) {
	resp, err := c.doRequest(method, apiPath, query, body, contentType, nil)
	if err != nil {
		return nil, "", 0, err
	}
	defer resp.Body.Close()
	stopBodyLog := c.startWaitLog(method, apiPath+" body", time.Now())
	data, readErr := io.ReadAll(resp.Body)
	stopBodyLog()
	if readErr != nil {
		return nil, "", resp.StatusCode, readErr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return data, resp.Header.Get("Content-Type"), resp.StatusCode, fmt.Errorf("remote Comigo %s %s failed: %s", method, apiPath, resp.Status)
	}
	return data, resp.Header.Get("Content-Type"), resp.StatusCode, nil
}

func (c *Client) doRequest(method string, apiPath string, query url.Values, body []byte, contentType string, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest(method, c.apiURL(apiPath, query), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return c.doHTTP(req, method, apiPath)
}

func (c *Client) login() error {
	form := url.Values{}
	form.Set("username", c.username)
	form.Set("password", c.password)
	req, err := http.NewRequest(http.MethodPost, c.apiURL("/api/login", nil), strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.doHTTP(req, http.MethodPost, "/api/login")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("remote Comigo login failed: %s", resp.Status)
	}
	return nil
}

// doHTTP 在等待远端 Comigo 响应时定期打印进度，避免长超时看起来像卡死。
func (c *Client) doHTTP(req *http.Request, method string, apiPath string) (*http.Response, error) {
	stopLog := c.startWaitLog(method, apiPath, time.Now())
	defer stopLog()
	return c.http.Do(req)
}

// startWaitLog 每 2 秒输出一次远端等待状态，stop 函数可重复调用。
func (c *Client) startWaitLog(method string, apiPath string, start time.Time) func() {
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				logger.Infof(locale.GetString("log_remote_comigo_waiting"), method, apiPath, c.baseURL.String(), time.Since(start).Round(time.Second), c.http.Timeout)
			case <-done:
				return
			}
		}
	}()
	var stopped bool
	return func() {
		if stopped {
			return
		}
		stopped = true
		close(done)
	}
}

func (c *Client) apiURL(apiPath string, query url.Values) string {
	u := *c.baseURL
	u.Path = path.Join(c.baseURL.Path, apiPath)
	if strings.HasSuffix(apiPath, "/") && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	u.RawQuery = query.Encode()
	return u.String()
}

// LocalizeBook 把远端书籍元数据转换成本机书库中的远程书籍元数据。
func LocalizeBook(storeURL string, remoteBook *model.Book) *model.Book {
	return LocalizeBookInShelf(storeURL, remoteBook, "", "")
}

// LocalizeBookInShelf 把远端书籍放进指定的远端顶级书库分组。
func LocalizeBookInShelf(storeURL string, remoteBook *model.Book, shelfKey string, shelfName string) *model.Book {
	if remoteBook == nil {
		return nil
	}
	localBook := *remoteBook
	remoteBookID := remoteBook.BookID
	key, err := StoreKey(storeURL)
	if err != nil {
		key = ""
	}
	localBook.BookID = LocalBookID(storeURL, remoteBookID)
	localBook.StoreUrl = storeURL
	localBook.RemoteURL = storeURL
	localBook.RemoteBookID = remoteBookID
	localBook.RemoteStoreKey = key
	localBook.RemoteShelfKey = shelfKey
	localBook.RemoteShelfName = shelfName
	localBook.IsRemote = true
	localBook.BookMarks = nil
	localBook.ChildBooksID = localizeChildIDs(storeURL, remoteBook.ChildBooksID)
	localBook.Cover.Url = LocalCoverURL(localBook.BookID, key, remoteBook.Cover.Url)
	localBook.PageInfos = make(model.PageInfos, len(remoteBook.PageInfos))
	for i := range remoteBook.PageInfos {
		localBook.PageInfos[i] = remoteBook.PageInfos[i]
		localBook.PageInfos[i].Url = LocalFileURL(localBook.BookID, key, remoteBook.PageInfos[i])
	}
	return &localBook
}

// ResolveBook 根据本机 BookID 和 remote_store 参数实时获取远端书籍详情。
func ResolveBook(storeURLs []string, remoteStoreKey string, localBookID string, sortBy string, timeoutSeconds int) (*model.Book, bool, error) {
	if remoteStoreKey == "" {
		return nil, false, nil
	}
	storeURL, err := MatchStoreByKey(storeURLs, remoteStoreKey)
	if err != nil {
		return nil, true, err
	}
	localBook, err := model.IStore.GetBook(localBookID)
	if err != nil {
		return nil, true, err
	}
	if localBook.RemoteBookID == "" {
		return nil, true, fmt.Errorf("book is not a remote Comigo book: %s", localBookID)
	}
	client, err := NewClient(storeURL, timeoutSeconds)
	if err != nil {
		return nil, true, err
	}
	remoteBook, err := client.GetBook(localBook.RemoteBookID, sortBy)
	if err != nil {
		return nil, true, err
	}
	refreshedBook := LocalizeBookInShelf(storeURL, remoteBook, localBook.RemoteShelfKey, localBook.RemoteShelfName)
	return refreshedBook, true, nil
}

func localizeChildIDs(storeURL string, remoteIDs []string) []string {
	localIDs := make([]string, 0, len(remoteIDs))
	for _, id := range remoteIDs {
		localIDs = append(localIDs, LocalBookID(storeURL, id))
	}
	return localIDs
}

// LocalCoverURL 生成本机封面代理 URL。
func LocalCoverURL(localBookID string, remoteStoreKey string, _ string) string {
	query := url.Values{}
	query.Set("id", localBookID)
	if remoteStoreKey != "" {
		query.Set(RemoteStoreQuery, remoteStoreKey)
	}
	return "/api/get-cover?" + query.Encode()
}

// LocalFileURL 生成本机图片代理 URL。
func LocalFileURL(localBookID string, remoteStoreKey string, page model.PageInfo) string {
	query := url.Values{}
	if page.Url != "" {
		if parsedURL, err := url.Parse(page.Url); err == nil {
			query = parsedURL.Query()
		}
	}
	query.Set("id", localBookID)
	if query.Get("filename") == "" && page.Name != "" {
		query.Set("filename", page.Name)
	}
	if remoteStoreKey != "" {
		query.Set(RemoteStoreQuery, remoteStoreKey)
	}
	return "/api/get-file?" + query.Encode()
}
