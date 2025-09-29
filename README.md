# Caching Proxy Server 🚀

A simple **caching proxy server** written in **Go**, which forwards requests to an origin server and caches responses to improve performance.

---

## Features

* Forward HTTP requests to any origin server.
* Cache responses in-memory with TTL (time-to-live).
* Adds `X-Cache` headers to indicate cache hits or misses:

  * `X-Cache: HIT` → response served from cache.
  * `X-Cache: MISS` → response fetched from origin server.
* CLI commands:

  * `serve` → start proxy server.
  * `clear-cache` → clear the cache.
* Thread-safe cache with automatic expiration.
* Structured logging with `logrus`.
* Unit tests for caching behavior.

---

## Installation

Make sure you have **Go >= 1.23** installed.

```bash
git clone https://github.com/rohan44942/caching-proxy.git
cd caching-proxy
go mod download
go build -o caching-proxy ./server
```

---

## Usage

### Start the proxy server

```bash
./caching-proxy serve --port 3000 --origin http://dummyjson.com
#start 
go run main.go serve --port 3000 --origin http://dummyjson.com --ttl 30
```

* **`--port`** → port where proxy server runs (e.g., 3000)
* **`--origin`** → the origin server to forward requests to
* **`--ttl`** → total time to live for cache

### Clear cache

```bash
./caching-proxy clear-cache
```

---

## Example

```bash
curl -i http://localhost:3000/products
```

* **First request** → cache MISS, fetches from origin:

```
X-Cache: MISS
```

* **Second request (same endpoint)** → cache HIT:

```
X-Cache: HIT
```

---

## Logging

* Logs are printed to **stdout** by default.
* Debug, info, and error logs available.
* Optional: write logs to a file by configuring `logrus` output.

Example:

```
INFO[0000] Proxy server running port=3000 origin=http://dummyjson.com
DEBUG[0002] Cache MISS, fetched from origin url=http://dummyjson.com/products
DEBUG[0005] Cache HIT url=/products cachedAge=5s
```

---

## Docker

### Build Docker image

```bash
docker build -t caching-proxy .
```

### Run Docker container

```bash
docker run -p 3000:3000 caching-proxy serve --port 3000 --origin http://dummyjson.com --ttl 30
```

* The proxy is available at **[http://localhost:3000](http://localhost:3000)** inside Docker.


## Project Structure (Not complete)

```
caching-proxy/
│
├── main.go   # entrypoint
├── cmd/internal/cache/             # caching logic
├── cmd/internal/config/            # configuration parsing
├── server/                     # server logic
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

---

## Future Improvements

* Replace in-memory cache with **Redis** for persistence.
* Add **config file support** (`YAML`/`JSON`) for ports, origin, TTL, and log settings.
* Add **unit tests for server** routes (`X-Cache` headers).
* Add **CI/CD pipeline** for Docker image builds and tests.

---

## License

MIT License © 2025 rohan44942
