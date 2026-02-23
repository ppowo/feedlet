package httpclient

import (
	"context"
	"crypto/tls"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var client *retryablehttp.Client

// userAgents is a pool of realistic browser User-Agent strings to rotate through.
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.2 Safari/605.1.15",
}

var (
	rng   *rand.Rand
	rngMu sync.Mutex
)

func init() {
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Use a custom transport that disables HTTP/2.
	// Go's default HTTP/2 TLS fingerprint is well-known and blocked
	// by anti-bot services (e.g. Reddit/Cloudflare).
	transport := &http.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	client = retryablehttp.NewClient()
	client.HTTPClient.Timeout = 30 * time.Second
	client.HTTPClient.Transport = transport
	client.RetryMax = 3
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 30 * time.Second
	client.Logger = nil

	client.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			return true, nil
		}
		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode <= 504) {
			return true, nil
		}
		return false, nil
	}
}

func GetClient() *retryablehttp.Client {
	return client
}

// RandomUserAgent returns a random User-Agent string from the pool.
func RandomUserAgent() string {
	rngMu.Lock()
	defer rngMu.Unlock()
	return userAgents[rng.Intn(len(userAgents))]
}
