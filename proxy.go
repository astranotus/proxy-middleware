package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxy(target *url.URL) *Proxy {
	return &Proxy{
		target: target,
		proxy:  httputil.NewSingleHostReverseProxy(target),
	}
}

func (p *Proxy) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api/") {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Hello from API"))
				return
			}
			p.proxy.ServeHTTP(w, r)
		})
	}
}
