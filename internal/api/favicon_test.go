package api

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestFaviconUsesActualContentType(t *testing.T) {
	s, h, cleanup := setupTestServer(t)
	defer cleanup()
	s.cfg.WebDir = t.TempDir()
	path := filepath.Join(s.cfg.WebDir, "favicon.ico")
	if err := os.WriteFile(path, []byte{0, 0, 1, 0}, 0600); err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/favicon.ico", nil))
	if rec.Code != 200 || rec.Header().Get("Content-Type") != "image/x-icon" {
		t.Fatalf("ICO: %d %s", rec.Code, rec.Header().Get("Content-Type"))
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/favicon.ico", nil))
	if rec.Header().Get("Content-Type") != "image/svg+xml" {
		t.Fatal("SVG fallback must keep its content type")
	}
}
