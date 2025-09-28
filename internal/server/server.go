package server

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/rohan44942/caching-proxy/internal/cache"
	"github.com/rohan44942/caching-proxy/internal/config"
	"github.com/sirupsen/logrus"
)

func Start(cfg config.Config) error {
	originURL, err := url.Parse(cfg.Origin)
	if err != nil {
		return fmt.Errorf("invalid origin URL: %w", err)
	}

	// Configure logger
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	file, err := os.OpenFile("proxy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Warn("Failed to log to file, using default stderr")
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		key := r.Method + ":" + r.URL.String()

		// Try cache
		if cachedResp, ok := cache.GlobalCache.Get(key); ok {
			logrus.WithFields(logrus.Fields{
				"url":        r.URL.String(),
				"method":     r.Method,
				"cache":      "HIT",
				"cachedAge":  time.Since(cachedResp.Timestamp).String(),
				"statusCode": cachedResp.StatusCode,
			}).Debug("serving from cache")

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
			logrus.WithError(err).Error("failed to create request to origin")
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
			logrus.WithError(err).Error("failed to reach origin server")
			http.Error(w, "failed to reach origin server", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logrus.WithError(err).Error("failed to read origin response body")
			http.Error(w, "failed to read origin response", http.StatusInternalServerError)
			return
		}
		cache.GlobalCache.Set(key, resp, body)

		logrus.WithFields(logrus.Fields{
			"url":        target.String(),
			"method":     r.Method,
			"cache":      "MISS",
			"statusCode": resp.StatusCode,
		}).Info("fetched from origin and cached")

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

	logrus.WithFields(logrus.Fields{
		"port":   cfg.Port,
		"origin": cfg.Origin,
	}).Info("Proxy server running")

	return http.ListenAndServe(addr, nil)
}
