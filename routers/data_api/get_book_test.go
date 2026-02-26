package data_api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestGetBook_MissingID_ReturnsContract(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/get-book", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if err := GetBook(c); err != nil {
		t.Fatalf("GetBook returned error: %v", err)
	}
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if body["code"] != "missing_param" {
		t.Fatalf("expected code missing_param, got %v", body["code"])
	}
	if body["message"] == "" {
		t.Fatalf("expected non-empty message")
	}
}
