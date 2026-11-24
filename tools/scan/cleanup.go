package scan

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/yumenaka/comigo/assets/locale"
	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
	"github.com/yumenaka/comigo/store"
	"github.com/yumenaka/comigo/tools"
	"github.com/yumenaka/comigo/tools/logger"
)

// CleanupPreview 只返回统计与少量示例；删除路径始终由服务器重新计算。
type CleanupPreview struct {
	Token     string   `json:"token"`
	Files     int      `json:"files"`
	Bytes     int64    `json:"bytes"`
	Books     int      `json:"books"`
	Pages     int      `json:"pages"`
	Bookmarks int      `json:"bookmarks"`
	Failures  int      `json:"failures"`
	Skipped   int      `json:"skipped"`
	Samples   []string `json:"samples"`
}

type cleanupFile struct {
	Root     string
	Name     string
	Size     int64
	Modified int64
}

type cleanupPlan struct {
	Preview        CleanupPreview
	Files          []cleanupFile
	Books          []*model.Book
	Orphans        map[string][]int64
	Failures       scanFailureCache
	RemoveFailures []string
}

// 数据库额外提供软删除书籍与无主页面、书签；内存模式无需 SQL 依赖。
type cleanupDatabase interface {
	ListCleanupBooks() ([]*model.Book, error)
	ListOrphanRows() (map[string][]int64, error)
	DeleteOrphanRows(map[string][]int64) error
	ReclaimCleanupSpace() error
}

var cleanupID = regexp.MustCompile(`^[0-9A-Za-z]{7,64}$`)
var cleanupCover = regexp.MustCompile(`^([0-9A-Za-z]{7,64})(?:_h[0-9]+|_original)?\.jpg$`)
var ErrCleanupChanged = errors.New("cleanup preview changed")

// CleanupInvalidData 与扫描串行，确认时重算并核对整个计划，拒绝执行过期预览。
func CleanupInvalidData(token string) (*CleanupPreview, error) {
	scanMutex.Lock()
	defer scanMutex.Unlock()
	scanFailureCacheMu.Lock()
	defer scanFailureCacheMu.Unlock()
	plan, err := collectCleanup()
	if err != nil {
		return nil, err
	}
	if token == "" {
		return &plan.Preview, nil
	}
	if token != plan.Preview.Token {
		return nil, ErrCleanupChanged
	}
	// 先清理索引，再处理文件；任何失败都向用户报告，保留剩余项供重新预览。
	for _, book := range plan.Books {
		if err := model.IStore.DeleteBook(book.BookID); err != nil {
			return nil, err
		}
	}
	if db, ok := model.IStore.(cleanupDatabase); ok {
		if err := db.DeleteOrphanRows(plan.Orphans); err != nil {
			return nil, err
		}
	}
	if len(plan.RemoveFailures) > 0 {
		for _, key := range plan.RemoveFailures {
			delete(plan.Failures, key)
		}
		if err := saveScanFailureCache(plan.Failures); err != nil {
			return nil, err
		}
	}
	deletedFiles := map[string]bool{}
	for _, f := range plan.Files {
		root, err := os.OpenRoot(f.Root)
		if err != nil {
			return nil, err
		}
		err = root.Remove(f.Name)
		// 仅尝试移除空目录，不递归删除未展示的文件。
		if err == nil || errors.Is(err, fs.ErrNotExist) {
			for dir := filepath.Dir(f.Name); dir != "."; dir = filepath.Dir(dir) {
				if root.Remove(dir) != nil {
					break
				}
			}
		}
		root.Close()
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			return nil, err
		}
		deletedFiles[filepath.Join(f.Root, f.Name)] = true
	}
	if ram, ok := model.IStore.(*store.StoreInRam); ok {
		configDir, err := config.GetConfigDir()
		if err != nil {
			return nil, err
		}
		metadata, err := filepath.Abs(filepath.Join(configDir, "metadata"))
		if err != nil {
			return nil, err
		}
		metadata, _ = filepath.EvalSymlinks(metadata)
		ram.PendingBooks.Range(func(key, value any) bool {
			book := value.(*model.Book)
			if deletedFiles[filepath.Join(metadata, book.GetStoreID(), book.BookID+".json")] {
				ram.PendingBooks.Delete(key)
			}
			return true
		})
	}
	if len(plan.Books) > 0 {
		if err := model.IStore.GenerateBookGroup(); err != nil {
			return nil, err
		}
	}
	// 所有删除及书组重建提交后再回收；回收失败不改变已完成的清理结果。
	if db, ok := model.IStore.(cleanupDatabase); ok {
		if err := db.ReclaimCleanupSpace(); err != nil {
			logger.Infof(locale.GetString("log_cleanup_sqlite_reclaim_failed"), err)
		}
	}
	return &plan.Preview, nil
}

// collectCleanup 保留可恢复的旧索引；版本过旧本身不代表书籍已经失效。
func collectCleanup() (*cleanupPlan, error) {
	p := &cleanupPlan{Preview: CleanupPreview{Samples: []string{}}, Orphans: map[string][]int64{}}
	var books []*model.Book
	var err error
	if db, ok := model.IStore.(cleanupDatabase); ok {
		books, err = db.ListCleanupBooks()
		if err != nil {
			return nil, err
		}
		p.Orphans, err = db.ListOrphanRows()
		if err != nil {
			return nil, err
		}
	} else {
		books, err = model.IStore.ListBooks()
		if err != nil {
			return nil, err
		}
	}
	stores := map[string]bool{}
	for _, url := range config.GetCfg().StoreUrls {
		if !tools.IsRemoteStoreURL(url) {
			url, err = filepath.Abs(url)
			if err != nil {
				return nil, err
			}
		}
		stores[url] = true
	}
	// 即使书库已从配置移除，也不能把其原目录当作可递归清空的缓存。
	protected := maps.Clone(stores)
	for _, b := range books {
		if b.StoreUrl != "" {
			protected[b.StoreUrl] = true
		}
	}
	// 无法访问整个书库时按离线处理，不把磁盘未挂载误判为全库删除。
	available := map[string]bool{}
	for url := range stores {
		if tools.IsRemoteStoreURL(url) {
			continue
		}
		if _, err := os.Stat(url); err == nil {
			available[url] = true
		} else {
			p.Preview.Skipped++
		}
	}
	missing := func(url, path string, remote bool) bool {
		if url == "" {
			return false
		}
		if !stores[url] {
			return true
		}
		if remote || !available[url] || path == "" {
			return false
		}
		_, err := os.Stat(path)
		return errors.Is(err, fs.ErrNotExist)
	}
	live := map[string]bool{}
	removed := map[string]bool{}
	replacements := map[[3]string]string{}
	for _, b := range books {
		if config.SupportsBookVersion(b.CreatedByVersion) && !b.Deleted && !missing(b.StoreUrl, b.BookPath, b.IsRemote) {
			replacements[b.BookIdentity()] = b.BookID
		}
	}
	invalid := func(b *model.Book) bool {
		replacement := replacements[b.BookIdentity()]
		return b.Deleted || missing(b.StoreUrl, b.BookPath, b.IsRemote) ||
			(!config.SupportsBookVersion(b.CreatedByVersion) && replacement != "" && replacement != b.BookID)
	}
	for _, b := range books {
		if invalid(b) {
			removed[b.BookID] = true
		} else {
			live[b.BookID] = true
		}
	}
	if ram, ok := model.IStore.(*store.StoreInRam); ok {
		ram.PendingBooks.Range(func(_, value any) bool {
			b := value.(*model.Book)
			if !invalid(b) {
				live[b.BookID] = true
			}
			return true
		})
	}
	// 连带移除没有有效子书的书组，避免悬空的父级引用。
	for changed := true; changed; {
		changed = false
		for _, b := range books {
			if b.Type != model.TypeBooksGroup || removed[b.BookID] {
				continue
			}
			hasChild := slices.ContainsFunc(b.ChildBooksID, func(id string) bool { return live[id] })
			if !hasChild {
				removed[b.BookID] = true
				delete(live, b.BookID)
				changed = true
			}
		}
	}
	for _, b := range books {
		if removed[b.BookID] {
			p.Books = append(p.Books, b)
			p.Preview.Pages += len(b.PageInfos)
			p.Preview.Bookmarks += len(b.BookMarks)
		}
	}
	p.Preview.Books = len(p.Books)
	p.Preview.Pages += len(p.Orphans["page_infos"])
	p.Preview.Bookmarks += len(p.Orphans["bookmarks"])
	configDir, err := config.GetConfigDir()
	if err != nil {
		return nil, err
	}
	// 配置目录一定存在；统一为真实绝对路径，保护相对配置路径和目录别名。
	configDir, err = filepath.Abs(configDir)
	if err != nil {
		return nil, err
	}
	configDir, err = filepath.EvalSymlinks(configDir)
	if err != nil {
		return nil, err
	}
	metadata := filepath.Join(configDir, "metadata")
	// 先读取所有 JSON，避免封面或图片缓存先于其待升级索引被判为无主。
	err = p.walk(metadata, protected, func(root *os.Root, name string, info fs.FileInfo) error {
		parts := strings.Split(filepath.ToSlash(name), "/")
		id := strings.TrimSuffix(filepath.Base(name), ".json")
		if len(parts) != 2 || filepath.Ext(name) != ".json" || !cleanupID.MatchString(id) {
			return nil
		}
		data, err := root.ReadFile(name)
		if err != nil {
			return err
		}
		var b struct {
			model.Book
			Path      string `json:"book_path"`
			Store     string `json:"store_url"`
			RemoteID  string `json:"remote_book_id"`
			RemoteKey string `json:"remote_store_key"`
		}
		if json.Unmarshal(data, &b) != nil || b.BookID != id || b.Store == "" || (b.Path == "" && b.RemoteID == "") {
			live[id] = true
			p.Preview.Skipped++
			return nil
		}
		b.BookPath, b.StoreUrl = b.Path, b.Store
		b.RemoteBookID, b.RemoteStoreKey = b.RemoteID, b.RemoteKey
		protected[b.StoreUrl] = true
		if parts[0] != b.GetStoreID() {
			live[id] = true
			p.Preview.Skipped++
			return nil
		}
		if invalid(&b.Book) || removed[id] {
			p.addFile(root.Name(), name, info)
		} else {
			live[id] = true
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = p.walk(metadata, protected, func(root *os.Root, name string, info fs.FileInfo) error {
		if len(strings.Split(filepath.ToSlash(name), "/")) != 2 {
			return nil
		}
		match := cleanupCover.FindStringSubmatch(filepath.Base(name))
		if len(match) > 0 && !live[match[1]] {
			p.addFile(root.Name(), name, info)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	cacheDir := config.GetCfg().CacheDir
	if cacheDir != "" {
		// metadata 独立清理；图片缓存遍历不触碰配置目录中的插件等用户数据。
		err = p.walk(cacheDir, protected, func(root *os.Root, name string, info fs.FileInfo) error {
			if withinPath(configDir, filepath.Join(root.Name(), name)) {
				return nil
			}
			parts := strings.Split(filepath.ToSlash(name), "/")
			if len(parts) < 2 {
				return nil
			}
			id := parts[0]
			if id == "remote_books" {
				if len(parts) < 3 {
					return nil
				}
				id = parts[1]
			}
			if cleanupID.MatchString(id) && !live[id] {
				p.addFile(root.Name(), name, info)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	p.Failures, err = loadScanFailureCache()
	if err != nil {
		return nil, err
	}
	for key, r := range p.Failures {
		if missing(r.StoreURL, r.FilePath, r.IsRemote) || shouldRetryFailedArchiveByVersion(r.CreatedByVersion, config.GetVersion()) {
			p.RemoveFailures = append(p.RemoveFailures, key)
		}
	}
	p.Preview.Failures = len(p.RemoveFailures)
	slices.SortFunc(p.Files, func(a, b cleanupFile) int { return strings.Compare(a.Root+"/"+a.Name, b.Root+"/"+b.Name) })
	slices.SortFunc(p.Books, func(a, b *model.Book) int { return strings.Compare(a.BookID, b.BookID) })
	slices.Sort(p.RemoveFailures)
	for _, b := range p.Books {
		p.sample("index/" + b.BookID + " · " + b.Title)
	}
	// 指纹覆盖所有候选项及配置，而非仅覆盖 UI 中的示例。
	data, err := json.Marshal([]any{p.Files, p.Books, p.Orphans, p.RemoveFailures, p.Failures, config.GetCfg().StoreUrls, config.GetCfg().CacheDir, configDir})
	if err != nil {
		return nil, err
	}
	digest := sha256.Sum256(data)
	p.Preview.Token = hex.EncodeToString(digest[:])
	return p, nil
}

// sample 限制响应体大小，统计仍覆盖所有清理对象。
func (p *cleanupPlan) sample(name string) {
	if len(p.Preview.Samples) < 20 {
		p.Preview.Samples = append(p.Preview.Samples, name)
	}
}

// addFile 记录文件大小、时间和相对路径，用于预览及确认时校验。
func (p *cleanupPlan) addFile(root, name string, info fs.FileInfo) {
	p.Files = append(p.Files, cleanupFile{root, name, info.Size(), info.ModTime().UnixNano()})
	p.Preview.Files++
	p.Preview.Bytes += info.Size()
	p.sample(filepath.Base(root) + "/" + filepath.ToSlash(name))
}

// walk 仅检查生成目录内的普通文件；跳过软链接和与源书库重叠的根目录。
func (p *cleanupPlan) walk(path string, stores map[string]bool, visit func(*os.Root, string, fs.FileInfo) error) error {
	resolved, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	resolved, err = filepath.EvalSymlinks(resolved)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for url := range stores {
		if tools.IsRemoteStoreURL(url) {
			continue
		}
		source, err := filepath.EvalSymlinks(url)
		if err != nil {
			source = url
		}
		if withinPath(resolved, source) || withinPath(source, resolved) {
			p.Preview.Skipped++
			return nil
		}
	}
	root, err := os.OpenRoot(resolved)
	if err != nil {
		return err
	}
	defer root.Close()
	return fs.WalkDir(root.FS(), ".", func(name string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			p.Preview.Skipped++
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		return visit(root, filepath.FromSlash(name), info)
	})
}

// withinPath 按路径层级比较，不能把同名前缀目录当作子目录。
func withinPath(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	return err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}
