package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rohan44942/caching-proxy/internal/cache"
	"github.com/rohan44942/caching-proxy/internal/config"
)

func Start(cfg config.Config) error {
	originURL, err := url.Parse(cfg.Origin)
	if err != nil {
		return fmt.Errorf("invalid origin URL: %w", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + ":" + r.URL.String()
		cachedResp, ok := cache.GlobalCache.Get(key)
		// fmt.Printf("inside the function %v value of ok %v\n", len(cachedResp.Body), ok)
		if ok {
			for name, values := range cachedResp.Header {
				for _, value := range values {
					w.Header().Add(name, value)
				}
			}
			w.Header().Set("X-Cache", "HIT")
			w.WriteHeader(cachedResp.StatusCode)
			w.Write(cachedResp.Body)
			return
		}

		target := originURL.ResolveReference(r.URL)
		req, err := http.NewRequest(r.Method, target.String(), r.Body)
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "failed to reach origin server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "failed to read origin response", http.StatusInternalServerError)
			return
		}
		cache.GlobalCache.Set(key, resp, body)

		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	}

	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Proxy server running on %s, forwarding to %s\n", addr, cfg.Origin)
	return http.ListenAndServe(addr, nil)
}
