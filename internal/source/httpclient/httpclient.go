package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var client *retryablehttp.Client

func init() {
	client = retryablehttp.NewClient()
	client.HTTPClient.Timeout = 30 * time.Second
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