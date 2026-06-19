package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/pz7-docker/tasks/internal/authclient"
)

func mockAuthServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/validate" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") == "Bearer demo-token" {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status":"ok"}`))
			return
		}
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
}

func TestListTasksSuccess(t *testing.T) {
	authSrv := mockAuthServer(t)
	defer authSrv.Close()

	h := New(authclient.New(authSrv.URL))
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer demo-token")
	req.Header.Set("X-Request-ID", "test-001")
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func TestListTasksUnauthorized(t *testing.T) {
	authSrv := mockAuthServer(t)
	defer authSrv.Close()

	h := New(authclient.New(authSrv.URL))
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rr.Code)
	}
}
