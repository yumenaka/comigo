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
func (s *fakeModelStore) ListBooks() ([]*model.Book, error)        { return []*model.Book{}, nil }
func (s *fakeModelStore) GenerateBookGroup() error                 { return nil }
func (s *fakeModelStore) StoreBookMark(mark *model.BookMark) error { return nil }
func (s *fakeModelStore) GetBookMarks(bookID string) (*model.BookMarks, error) {
	return &model.BookMarks{}, nil
}

// TestGenerateBookGroup_ShouldCreateFolderGroupsForIntermediateDirs
// 目标：当一个目录自身不满足“目录型书籍入库条件”（例如无图片），但其更深层有图片目录/书籍时，
// 应该能生成对应的“文件夹/书组(TypeBooksGroup)”节点用于导航（例如 Steam / Steam 2）。
func TestGenerateBookGroup_ShouldCreateFolderGroupsForIntermediateDirs(t *testing.T) {
	model.IStore = &fakeModelStore{}

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
