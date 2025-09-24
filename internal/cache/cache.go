package cache

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

type CachedResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
	Timestamp  time.Time
}

type Cache struct {
	mu    sync.RWMutex
	store map[string]CachedResponse
}

func New() *Cache {
	return &Cache{
		store: make(map[string]CachedResponse),
	}
}

func (c *Cache) Get(key string) (CachedResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	resp, ok := c.store[key]
	return resp, ok
}

func (c *Cache) Set(key string, resp *http.Response, body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Deep copy headers
	headers := make(http.Header)
	for k, v := range resp.Header {
		copied := make([]string, len(v))
		copy(copied, v)
		headers[k] = copied
	}

	c.store[key] = CachedResponse{
		StatusCode: resp.StatusCode,
		Header:     headers,
		Body:       body,
		Timestamp:  time.Now(),
	}
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store = make(map[string]CachedResponse)
}

// Utility: make body reusable
func ReadAndCopyBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(body)) // reset body for re-read
	return body, nil
}
