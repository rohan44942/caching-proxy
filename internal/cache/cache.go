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
	TTL   time.Duration
}

func New(ttl time.Duration) *Cache {
	return &Cache{
		store: make(map[string]CachedResponse),
		TTL:   ttl,
	}
}

func (c *Cache) Get(key string) (CachedResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	resp, ok := c.store[key]
	// fmt.Printf("value of resp%v ", (resp))
	if !ok {
		return CachedResponse{}, false
	}
	if c.TTL > 0 && time.Since(resp.Timestamp) > c.TTL {
		// expired, remove and return miss
		c.mu.RUnlock()
		c.mu.Lock()
		delete(c.store, key)
		c.mu.Unlock()
		c.mu.RLock()
		return CachedResponse{}, false
	}
	return resp, true
}

func (c *Cache) Set(key string, resp *http.Response, body []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

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
	// fmt.Print("value of key, and body length is ", key, len(body))
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
	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
