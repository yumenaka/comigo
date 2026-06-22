package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/yumenaka/comigo/config"
	"github.com/yumenaka/comigo/model"
)

// 提供生成书籍组测试所需的最小书库实现，避免 model.NewBook 去重检查访问空书库。
type fakeModelStore struct{}

func (s *fakeModelStore) StoreBook(b *model.Book) error            { return nil }
func (s *fakeModelStore) GetBook(id string) (*model.Book, error)   { return nil, os.ErrNotExist }
func (s *fakeModelStore) DeleteBook(id string) error               { return nil }
func (s *fakeModelStore) ListBooks() ([]*model.Book, error)        { return fakeBooks, nil }
func (s *fakeModelStore) GenerateBookGroup() error                 { return nil }
func (s *fakeModelStore) StoreBookMark(mark *model.BookMark) error { return nil }
func (s *fakeModelStore) GetBookMarks(bookID string) (*model.BookMarks, error) {
	return &model.BookMarks{}, nil
}
func (s *fakeModelStore) DeleteBookMark(bookID string, markType model.MarkType, pageIndex int) error {
	return nil
}

var fakeBooks []*model.Book

func setFakeBooks(t *testing.T, books ...*model.Book) {
	t.Helper()
	oldBooks := fakeBooks
	oldStore := model.IStore
	fakeBooks = books
	model.IStore = &fakeModelStore{}
	t.Cleanup(func() {
		fakeBooks = oldBooks
		model.IStore = oldStore
	})
}

func useTempConfigDir(t *testing.T) string {
	t.Helper()
	configDir := t.TempDir()
	t.Setenv("COMIGO_CONFIG_DIR", configDir)
	oldConfigFile := config.GetCfg().ConfigFile
	config.GetCfg().ConfigFile = ""
	t.Cleanup(func() {
		config.GetCfg().ConfigFile = oldConfigFile
	})
	return configDir
}

func metaJSONPath(configDir string, book *model.Book) string {
	return filepath.Join(configDir, "metadata", book.GetStoreID(), book.BookID+".json")
}

func newStoreGroupTestBook(storeURL, id string) *model.Book {
	return &model.Book{BookInfo: model.BookInfo{
		BookID:   id,
		BookPath: filepath.Join(storeURL, id+".cbz"),
		StoreUrl: storeURL,
		Type:     model.TypeZip,
	}}
}

// 验证内存书库写入书籍时不会额外落盘元数据。
func TestStoreBookInMemoryDoesNotPersistMeta(t *testing.T) {
	configDir := useTempConfigDir(t)
	storeURL := filepath.Join(t.TempDir(), "books")
	book := newStoreGroupTestBook(storeURL, "book-in-memory")
	ramStore := &StoreInRam{}

	if err := ramStore.storeBookInMemory(book); err != nil {
		t.Fatalf("storeBookInMemory error: %v", err)
	}
	if _, err := ramStore.GetBook(book.BookID); err != nil {
		t.Fatalf("storeBookInMemory 没有写入内存 Store: %v", err)
	}
	if _, err := os.Stat(metaJSONPath(configDir, book)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("storeBookInMemory 不应写 metadata 文件，stat err=%v", err)
	}
}

// 验证需要持久化时书籍元数据会写入本地文件。
func TestStoreBookPersistsMetaAfterInMemoryStore(t *testing.T) {
	configDir := useTempConfigDir(t)
	storeURL := filepath.Join(t.TempDir(), "books")
	book := newStoreGroupTestBook(storeURL, "book-persisted")
	ramStore := &StoreInRam{}

	if err := ramStore.StoreBook(book); err != nil {
		t.Fatalf("StoreBook error: %v", err)
	}
	if _, err := ramStore.GetBook(book.BookID); err != nil {
		t.Fatalf("StoreBook 没有写入内存 Store: %v", err)
	}
	if _, err := os.Stat(metaJSONPath(configDir, book)); err != nil {
		t.Fatalf("StoreBook 应写 metadata 文件: %v", err)
	}
}

// 验证公开 JSON 会隐藏本地路径，但内部元数据仍保留真实路径。
func TestPublicJSONHidesLocalPathsButMetadataKeepsThem(t *testing.T) {
	storeURL := filepath.Join(t.TempDir(), "books")
	book := newStoreGroupTestBook(storeURL, "book-private-fields")
	book.RemoteURL = "sftp://user:password@example.com:22/books"
	pagePath := filepath.Join(storeURL, "001.jpg")
	book.PageInfos = model.PageInfos{{
		Name: "001.jpg",
		Path: pagePath,
		Url:  "/api/get-file?id=book-private-fields&filename=001.jpg",
	}}
	book.Cover = book.PageInfos[0]

	publicJSON, err := json.Marshal(book)
	if err != nil {
		t.Fatalf("marshal public book json: %v", err)
	}
	if strings.Contains(string(publicJSON), "book_path") ||
		strings.Contains(string(publicJSON), "store_url") ||
		strings.Contains(string(publicJSON), "remote_url") ||
		strings.Contains(string(publicJSON), book.RemoteURL) ||
		strings.Contains(string(publicJSON), pagePath) {
		t.Fatalf("普通 JSON 不应输出内部路径字段: %s", publicJSON)
	}

	metaJSON, err := marshalBookMetaJSON(book)
	if err != nil {
		t.Fatalf("marshal metadata json: %v", err)
	}
	metaText := string(metaJSON)
	if !strings.Contains(metaText, `"book_path"`) ||
		!strings.Contains(metaText, `"store_url"`) ||
		!strings.Contains(metaText, `"remote_url"`) {
		t.Fatalf("metadata JSON 必须保留内部路径字段: %s", metaText)
	}
	if !strings.Contains(metaText, pagePath) {
		t.Fatalf("metadata JSON 必须保留页面真实路径: %s", metaText)
	}

	var restored model.Book
	if err := json.Unmarshal(metaJSON, &restored); err != nil {
		t.Fatalf("unmarshal metadata json: %v", err)
	}
	if err := restoreBookMetaPrivateFields(metaJSON, &restored); err != nil {
		t.Fatalf("restore metadata private fields: %v", err)
	}
	if restored.BookPath != book.BookPath || restored.StoreUrl != book.StoreUrl || restored.RemoteURL != book.RemoteURL {
		t.Fatalf(
			"metadata 内部字段恢复失败: got path=%q store=%q remote=%q",
			restored.BookPath,
			restored.StoreUrl,
			restored.RemoteURL,
		)
	}
	if len(restored.PageInfos) != 1 || restored.PageInfos[0].Path != pagePath || restored.Cover.Path != pagePath {
		t.Fatalf("metadata 页面路径恢复失败: pages=%#v cover=%#v", restored.PageInfos, restored.Cover)
	}

	shelfJSON, err := json.Marshal(model.StoreBookInfo{
		StoreUrl:     storeURL,
		DisplayName:  filepath.Base(storeURL),
		ChildBookNum: 1,
		BookInfos:    model.BookInfos{book.BookInfo},
	})
	if err != nil {
		t.Fatalf("marshal store book info json: %v", err)
	}
	if strings.Contains(string(shelfJSON), "store_url") || strings.Contains(string(shelfJSON), storeURL) {
		t.Fatalf("书架分组普通 JSON 不应输出书库真实路径: %s", shelfJSON)
	}
}

// 验证多级本地目录会生成中间层书籍组。
// 当中间目录本身没有图片、但更深层有目录书籍时，仍应生成可导航的文件夹书组。
func TestGenerateBookGroup_ShouldCreateFolderGroupsForIntermediateDirs(t *testing.T) {
	setFakeBooks(t)

	tmp := t.TempDir()
	root := filepath.Join(tmp, "test")
	leaf1 := filepath.Join(root, "Steam", "Vol1")   // 叶子目录书籍
	leaf2 := filepath.Join(root, "Steam 2", "Vol1") // 叶子目录书籍

	if err := os.MkdirAll(leaf1, 0o755); err != nil {
		t.Fatalf("mkdir leaf1: %v", err)
	}
	if err := os.MkdirAll(leaf2, 0o755); err != nil {
		t.Fatalf("mkdir leaf2: %v", err)
	}

	s := &Store{
		StoreInfo: StoreInfo{BackendURL: root},
	}

	// 注意：这里不依赖扫描逻辑，直接构造“叶子目录型书籍”，模拟深层图片目录已被成功入库。
	// 深度以书库根目录为起点，root/test 下的 Steam/Vol1 深度为 1。
	b1 := &model.Book{BookInfo: model.BookInfo{
		BookID:   "leaf-steam-vol1",
		BookPath: leaf1,
		StoreUrl: root,
		Depth:    1,
		Type:     model.TypeDir,
		Modified: time.Now(),
	}}
	b2 := &model.Book{BookInfo: model.BookInfo{
		BookID:   "leaf-steam2-vol1",
		BookPath: leaf2,
		StoreUrl: root,
		Depth:    1,
		Type:     model.TypeDir,
		Modified: time.Now(),
	}}
	s.BookMap.Store(b1.BookID, b1)
	s.BookMap.Store(b2.BookID, b2)

	if err := s.GenerateBookGroup(); err != nil {
		t.Fatalf("GenerateBookGroup error: %v", err)
	}

	expectGroupPaths := map[string]bool{
		filepath.Join(root, "Steam"):   false,
		filepath.Join(root, "Steam 2"): false,
	}

	for _, value := range s.BookMap.Range {
		b := value.(*model.Book)
		if b.Type != model.TypeBooksGroup {
			continue
		}
		if _, ok := expectGroupPaths[b.BookPath]; ok {
			expectGroupPaths[b.BookPath] = true
		}
	}

	for p, ok := range expectGroupPaths {
		if !ok {
			t.Fatalf("expected books group for path %s, but it was missing", p)
		}
	}
}

// 验证首页书库名称只显示本地目录名，避免暴露父路径。
func TestDisplayStoreNameHidesLocalParentPath(t *testing.T) {
	localStore := filepath.Join(t.TempDir(), "books")
	if got := displayStoreName(localStore); got != "books" {
		t.Fatalf("本地书库显示名不正确: got %q want %q", got, "books")
	}
}

// 验证远程书库名称显示主机而不是完整带密码地址。
func TestDisplayStoreNameUsesRemoteHost(t *testing.T) {
	const remoteStore = "sftp://user:pass@example.com:22/manga"
	if got := displayStoreName(remoteStore); got != "example.com:22" {
		t.Fatalf("远程书库显示名不正确: got %q want %q", got, "example.com:22")
	}
}

// 验证远程 Comigo 书库会按远端书架拆分展示。
func TestTopOfShelfInfoSplitsRemoteComigoShelves(t *testing.T) {
	oldCfg := config.CopyCfg()
	t.Cleanup(func() {
		*config.GetCfg() = oldCfg
	})
	const remoteStore = "https://example.com"
	config.GetCfg().StoreUrls = []string{remoteStore}

	remoteBookA := bookForFolderTest("remote-a", "/remote/a.cbz", remoteStore, model.TypeZip)
	remoteBookA.Depth = 0
	remoteBookA.IsRemote = true
	remoteBookA.RemoteURL = remoteStore
	remoteBookA.RemoteBookID = "a"
	remoteBookA.RemoteStoreKey = "remote-key"
	remoteBookA.RemoteShelfKey = "shelf-a"
	remoteBookA.RemoteShelfName = "Shelf A"
	remoteBookB := bookForFolderTest("remote-b", "/remote/b.cbz", remoteStore, model.TypeZip)
	remoteBookB.Depth = 0
	remoteBookB.IsRemote = true
	remoteBookB.RemoteURL = remoteStore
	remoteBookB.RemoteBookID = "b"
	remoteBookB.RemoteStoreKey = "remote-key"
	remoteBookB.RemoteShelfKey = "shelf-b"
	remoteBookB.RemoteShelfName = "Shelf B"
	setFakeBooks(t, remoteBookA, remoteBookB)

	got, err := TopOfShelfInfo("filename")
	if err != nil {
		t.Fatalf("TopOfShelfInfo error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("远端顶级书库数量 = %d, want 2: %#v", len(got), got)
	}
	if got[0].DisplayName != "Shelf A" || got[1].DisplayName != "Shelf B" {
		t.Fatalf("远端顶级书库显示名不正确: %#v", got)
	}
	if got[0].ChildBookNum != 1 || got[1].ChildBookNum != 1 {
		t.Fatalf("远端顶级书库计数不正确: %#v", got)
	}
}

// 验证远程书籍组会继承远程访问所需字段。
func TestGeneratedRemoteBookGroupInheritsRemoteFields(t *testing.T) {
	const remoteStore = "sftp://example.com/manga"
	childA := bookForFolderTest("child-a", "/remote/album/a.cbz", remoteStore, model.TypeZip)
	childA.Depth = 1
	childA.IsRemote = true
	childA.RemoteURL = remoteStore
	childA.RemoteStoreKey = "remote-key"
	childA.RemoteShelfKey = "shelf-a"
	childA.RemoteShelfName = "Shelf A"
	childB := bookForFolderTest("child-b", "/remote/album/b.cbz", remoteStore, model.TypeZip)
	childB.Depth = 1
	childB.IsRemote = true
	childB.RemoteURL = remoteStore
	childB.RemoteStoreKey = "remote-key"
	childB.RemoteShelfKey = "shelf-a"
	childB.RemoteShelfName = "Shelf A"
	setFakeBooks(t, childA, childB)

	store := &Store{StoreInfo: StoreInfo{BackendURL: remoteStore}}
	store.BookMap.Store(childA.BookID, childA)
	store.BookMap.Store(childB.BookID, childB)
	if err := store.GenerateBookGroup(); err != nil {
		t.Fatalf("GenerateBookGroup error: %v", err)
	}

	var group *model.Book
	for _, value := range store.BookMap.Range {
		book := value.(*model.Book)
		if book.Type == model.TypeBooksGroup {
			group = book
		}
	}
	if group == nil {
		t.Fatal("未生成远程书组")
	}
	if !group.IsRemote || group.RemoteStoreKey != "remote-key" || group.RemoteShelfKey != "shelf-a" {
		t.Fatalf("远程书组没有继承远端字段: %#v", group.BookInfo)
	}
}

// 验证远程 Comigo 原始分组不会被本地自动分组覆盖。
func TestGenerateBookGroupKeepsRemoteComigoOriginalGroup(t *testing.T) {
	const remoteStore = "https://example.com"
	child := bookForFolderTest("child", "/remote/album/a.cbz", remoteStore, model.TypeZip)
	child.Depth = 1
	child.IsRemote = true
	child.RemoteURL = remoteStore
	child.RemoteBookID = "remote-child"
	child.RemoteStoreKey = "remote-key"
	child.RemoteShelfKey = "shelf-a"
	child.RemoteShelfName = "Shelf A"
	remoteGroup := bookForFolderTest("remote-group", "/remote/album", remoteStore, model.TypeBooksGroup)
	remoteGroup.Depth = 0
	remoteGroup.IsRemote = true
	remoteGroup.RemoteURL = remoteStore
	remoteGroup.RemoteBookID = "remote-group"
	remoteGroup.RemoteStoreKey = "remote-key"
	remoteGroup.RemoteShelfKey = "shelf-a"
	remoteGroup.RemoteShelfName = "Shelf A"
	remoteGroup.ChildBooksID = []string{child.BookID}
	setFakeBooks(t, remoteGroup, child)

	store := &Store{StoreInfo: StoreInfo{BackendURL: remoteStore}}
	store.BookMap.Store(remoteGroup.BookID, remoteGroup)
	store.BookMap.Store(child.BookID, child)
	if err := store.GenerateBookGroup(); err != nil {
		t.Fatalf("GenerateBookGroup error: %v", err)
	}
	if _, ok := store.BookMap.Load(remoteGroup.BookID); !ok {
		t.Fatal("远端 Comigo 原本返回的书组被本地书组生成流程删除")
	}
	groupCount := 0
	for _, value := range store.BookMap.Range {
		book := value.(*model.Book)
		if book.Type == model.TypeBooksGroup {
			groupCount++
		}
	}
	if groupCount != 1 {
		t.Fatalf("Comigo 远端书组不应被本地路径规则重复生成: got %d want 1", groupCount)
	}
}

// 验证按文件夹取书时本地书库使用真实父路径匹配。
func TestGetBookInfoListByBookFolderUsesRealLocalParentPath(t *testing.T) {
	tmp := t.TempDir()
	rootA := filepath.Join(tmp, "root-a")
	rootB := filepath.Join(tmp, "root-b")
	currentPath := filepath.Join(rootA, "Album", "01.mp4")
	siblingPath := filepath.Join(rootA, "Album", "02.mp4")
	collidingPath := filepath.Join(rootB, "Album", "03.mp4")

	current := bookForFolderTest("current", currentPath, rootA, model.TypeVideo)
	sibling := bookForFolderTest("sibling", siblingPath, rootA, model.TypeAudio)
	colliding := bookForFolderTest("colliding", collidingPath, rootB, model.TypeVideo)
	setFakeBooks(t, current, sibling, colliding)

	got, err := GetBookInfoListByBookFolder(current)
	if err != nil {
		t.Fatalf("GetBookInfoListByBookFolder error: %v", err)
	}
	assertBookIDs(t, *got, "current", "sibling")
}

// 验证按文件夹取书时不同远程书库不会互相串数据。
func TestGetBookInfoListByBookFolderSeparatesRemoteStores(t *testing.T) {
	current := bookForFolderTest("current", "/manga/Album/01.cbz", "webdav://a.example/manga", model.TypeZip)
	current.IsRemote = true
	sibling := bookForFolderTest("sibling", "/manga/Album/02.cbz", "webdav://a.example/manga", model.TypeZip)
	sibling.IsRemote = true
	colliding := bookForFolderTest("colliding", "/manga/Album/03.cbz", "webdav://b.example/manga", model.TypeZip)
	colliding.IsRemote = true
	setFakeBooks(t, current, sibling, colliding)

	got, err := GetBookInfoListByBookFolder(current)
	if err != nil {
		t.Fatalf("GetBookInfoListByBookFolder error: %v", err)
	}
	assertBookIDs(t, *got, "current", "sibling")
}

func bookForFolderTest(id, bookPath, storeURL string, bookType model.SupportFileType) *model.Book {
	parentFolder := filepath.Base(filepath.Dir(bookPath))
	return &model.Book{BookInfo: model.BookInfo{
		BookID:       id,
		Title:        id,
		BookPath:     bookPath,
		ParentFolder: parentFolder,
		StoreUrl:     storeURL,
		Type:         bookType,
	}}
}

func assertBookIDs(t *testing.T, got model.BookInfos, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("同目录书籍数量不正确: got %d want %d, list=%v", len(got), len(want), got)
	}
	gotIDs := make(map[string]bool, len(got))
	for _, book := range got {
		gotIDs[book.BookID] = true
	}
	for _, id := range want {
		if !gotIDs[id] {
			t.Fatalf("缺少同目录书籍 %q, got=%v", id, gotIDs)
		}
	}
}
