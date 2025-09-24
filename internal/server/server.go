package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/rohan44942/caching-proxy/internal/config"
)

func Start(cfg config.Config) error {
	originURL, err := url.Parse(cfg.Origin)
	if err != nil {
		return fmt.Errorf("invalid origin URL: %w", err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		target := originURL.ResolveReference(r.URL) // combine origin + request path

		req, err := http.NewRequest(r.Method, target.String(), r.Body)
		if err != nil {
			http.Error(w, "failed to create request", http.StatusInternalServerError)
			return
		}

		// copy headers
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// forward request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "failed to reach origin server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// copy headers from origin
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// mark as MISS (Phase 2 always MISS)
		w.Header().Set("X-Cache", "MISS")
		w.WriteHeader(resp.StatusCode)

		// copy body
		io.Copy(w, resp.Body)
	}

	http.HandleFunc("/", handler)
	addr := fmt.Sprintf(":%d", cfg.Port)
	fmt.Printf("Proxy server running on %s, forwarding to %s\n", addr, cfg.Origin)
	return http.ListenAndServe(addr, nil)
}
