package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

func main() {
	target, _ := url.Parse("http://hugo:1313")
	proxy := NewProxy(target)

	r := chi.NewRouter()
	r.Use(proxy.Middleware())

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	log.Println("Proxy server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
