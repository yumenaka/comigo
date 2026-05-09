package store

import (
	"os"
	"path/filepath"
	"testing"
	"time"

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
