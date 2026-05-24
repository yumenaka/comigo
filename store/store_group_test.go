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

// fakeModelStore 仅用于 model.NewBook 过程中避免 nil IStore
// 说明：GenerateBookGroup 内部会调用 model.NewBook，而 NewBook 会用到 model.IStore.ListBooks() 做去重检查。
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

// TestGenerateBookGroup_ShouldCreateFolderGroupsForIntermediateDirs
// 目标：当一个目录自身不满足“目录型书籍入库条件”（例如无图片），但其更深层有图片目录/书籍时，
// 应该能生成对应的“文件夹/书组(TypeBooksGroup)”节点用于导航（例如 Steam / Steam 2）。
func TestGenerateBookGroup_ShouldCreateFolderGroupsForIntermediateDirs(t *testing.T) {
	setFakeBooks(t)

	tmp := t.TempDir()
	root := filepath.Join(tmp, "test")
	leaf1 := filepath.Join(root, "Steam", "Vol1")   // leaf dir book
	leaf2 := filepath.Join(root, "Steam 2", "Vol1") // leaf dir book

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
	// depth: root/test 下的 Steam/Vol1 => 相对路径 "Steam/Vol1" 深度为 1
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

func TestDisplayStoreNameHidesLocalParentPath(t *testing.T) {
	localStore := filepath.Join(t.TempDir(), "books")
	if got := displayStoreName(localStore); got != "books" {
		t.Fatalf("本地书库显示名不正确: got %q want %q", got, "books")
	}
}

func TestDisplayStoreNameUsesRemoteHost(t *testing.T) {
	const remoteStore = "sftp://user:pass@example.com:22/manga"
	if got := displayStoreName(remoteStore); got != "example.com:22" {
		t.Fatalf("远程书库显示名不正确: got %q want %q", got, "example.com:22")
	}
}

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
