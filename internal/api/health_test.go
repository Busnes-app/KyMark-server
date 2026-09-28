package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func healthResponse(t *testing.T, h http.Handler, path string) (int, map[string]any, string) {
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	var body map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	return w.Code, body, w.Body.String()
}

func TestHealthAliasesSharePublicSchemaAndCache(t *testing.T) {
	srv, h, cleanup := setupTestServer(t)
	defer cleanup()

	code, body, raw := healthResponse(t, h, "/healthz")
	if code != http.StatusOK || body["schema"] != "ky.health/1" || body["service"] != "kymark" || body["status"] != "ok" {
		t.Fatalf("healthz=%d %s", code, raw)
	}
	checks, ok := body["checks"].([]any)
	if !ok || len(checks) != 2 || checks[0].(map[string]any)["name"] != "database" || checks[1].(map[string]any)["name"] != "audit" {
		t.Fatalf("unexpected check scope: %s", raw)
	}
	srv.auditFailures.Add(1) // The second alias must use the first alias's cached evaluation.
	legacyCode, legacyBody, legacyRaw := healthResponse(t, h, "/api/health")
	if legacyCode != code || !reflect.DeepEqual(legacyBody, body) || legacyRaw != raw {
		t.Fatalf("aliases differ: /healthz=%d %s, /api/health=%d %s", code, raw, legacyCode, legacyRaw)
	}
	for _, private := range []string{"vault", "account", "device", "user", "auditWriteFailures"} {
		if strings.Contains(strings.ToLower(raw), strings.ToLower(private)) {
			t.Fatalf("public health contains %q: %s", private, raw)
		}
	}
}

func TestHealthDatabaseFailure(t *testing.T) {
	srv, h, cleanup := setupTestServer(t)
	defer cleanup()
	if err := srv.store.DB().Close(); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/healthz", "/api/health"} {
		code, body, raw := healthResponse(t, h, path)
		if code != http.StatusServiceUnavailable || body["status"] != "down" || body["service"] != "kymark" {
			t.Fatalf("%s=%d %s", path, code, raw)
		}
		if strings.Contains(raw, "sql: database is closed") {
			t.Fatalf("health leaked database error: %s", raw)
		}
	}
}

func TestHealthAuditFailureIsDegraded(t *testing.T) {
	srv, h, cleanup := setupTestServer(t)
	defer cleanup()
	srv.auditFailures.Add(1)
	for _, path := range []string{"/healthz", "/api/health"} {
		code, body, raw := healthResponse(t, h, path)
		if code != http.StatusOK || body["status"] != "degraded" || body["service"] != "kymark" || !strings.Contains(raw, "append_disabled") {
			t.Fatalf("%s=%d %s", path, code, raw)
		}
		if strings.Contains(raw, "auditWriteFailures") || strings.Contains(raw, "sql:") {
			t.Fatalf("health leaked internal detail: %s", raw)
		}
	}
}
