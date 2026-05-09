package websocket

import "testing"

func TestHandSyncPageMessageToFlipMode(t *testing.T) {
	debug := false
	WsDebug = &debug

	msg := Message{
		Type:       "flip_mode_sync_page",
		DataString: `{"book_id":"book-1","now_page_num":3}`,
	}
	got, ok := handSyncPageMessageToFlipMode(msg, "client-1")
	if !ok {
		t.Fatal("valid flip sync message should pass")
	}
	if got.Detail != "同步页数。" {
		t.Fatalf("Detail got %q", got.Detail)
	}

	invalid := Message{
		Type:       "flip_mode_sync_page",
		DataString: `{"book_id":"","now_page_num":3}`,
	}
	if _, ok := handSyncPageMessageToFlipMode(invalid, "client-1"); ok {
		t.Fatal("empty book_id should fail")
	}
}

func TestHandSyncPageMessageToScrollMode(t *testing.T) {
	debug := false
	WsDebug = &debug

	msg := Message{
		Type:       "scroll_mode_sync_page",
		DataString: `{"book_id":"book-1","now_page_num":3,"now_page_num_percent":0.5}`,
	}
	got, ok := handSyncPageMessageToScrollMode(msg, "client-1")
	if !ok {
		t.Fatal("valid scroll sync message should pass")
	}
	if got.Detail != "同步页数。" {
		t.Fatalf("Detail got %q", got.Detail)
	}

	invalid := Message{
		Type:       "scroll_mode_sync_page",
		DataString: `{"book_id":"book-1","now_page_num":3,"now_page_num_percent":1.5}`,
	}
	if _, ok := handSyncPageMessageToScrollMode(invalid, "client-1"); ok {
		t.Fatal("percent greater than 1 should fail")
	}
}
