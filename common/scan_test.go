package common

import (
	"fmt"
	"testing"
)

func TestScanPath(t *testing.T) {
	bookList, err := ScanPath("../test")
	if err != nil {
		t.Errorf("TestScanPath error")
	}
	fmt.Println(len(bookList))
	t.Log("hello world")
}
