package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestAPIHandler(t *testing.T) {
	target, _ := url.Parse("http://example.com") 
	proxy := NewProxy(target)

	r := chi.NewRouter()
	r.Use(proxy.Middleware())
	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", resp.StatusCode)
	}
	if string(body) != "Hello from API" {
		t.Errorf("Ожидалось тело 'Hello from API', получено %q", string(body))
	}
}

func TestProxyHandler(t *testing.T) {
	proxied := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxied = true
		w.WriteHeader(http.StatusTeapot)
		w.Write([]byte("proxied"))
	}))
	defer ts.Close()

	target, _ := url.Parse(ts.URL)
	proxy := NewProxy(target)

	r := chi.NewRouter()
	r.Use(proxy.Middleware())
	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/non-api", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if !proxied {
		t.Errorf("Ожидался вызов прокси, но он не был выполнен")
	}
	if resp.StatusCode != http.StatusTeapot {
		t.Errorf("Ожидался статус 418, получено %d", resp.StatusCode)
	}
	if string(body) != "proxied" {
		t.Errorf("Ожидалось тело 'proxied', получено %q", string(body))
	}
}
